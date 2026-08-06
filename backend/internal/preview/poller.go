package preview

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sriram687/GIT_WIZ/backend/internal/domain"
)

// DefaultPollInterval is the preview poller's scan interval when none is configured.
const DefaultPollInterval = 250 * time.Millisecond

type sessionPreviewSource interface {
	ListAllSessions(ctx context.Context) ([]domain.SessionRecord, error)
	GetProject(ctx context.Context, id string) (domain.ProjectRecord, bool, error)
}

type previewSetter interface {
	SetPreview(ctx context.Context, id domain.SessionID, previewURL string) (domain.Session, error)
}

// PollerConfig configures preview poller timing and logging.
type PollerConfig struct {
	Interval time.Duration
	Logger   *slog.Logger
}

// Poller watches active worker workspaces for static frontend entrypoints and
// persists preview URL refreshes through the normal session service path.
type Poller struct {
	source     sessionPreviewSource
	setter     previewSetter
	baseURL    string
	interval   time.Duration
	logger     *slog.Logger
	seen       map[domain.SessionID]entryState
	devServers map[domain.SessionID]*devServer
	mu         sync.Mutex
}

type devServer struct {
	cmd  *exec.Cmd
	port int
}

type entryState struct {
	path    string
	modUnix int64
	size    int64
	// cleared is set when the poller itself cleared the preview URL because the
	// workspace entry was missing. When the file reappears, shouldRefresh uses
	// this to re-discover even though the revision was bumped by the clear.
	cleared bool
}

// NewPoller constructs a preview poller over the supplied session source and setter.
func NewPoller(source sessionPreviewSource, setter previewSetter, baseURL string, cfg PollerConfig) *Poller {
	p := &Poller{
		source:     source,
		setter:     setter,
		baseURL:    baseURL,
		interval:   cfg.Interval,
		logger:     cfg.Logger,
		seen:       map[domain.SessionID]entryState{},
		devServers: map[domain.SessionID]*devServer{},
	}
	if p.interval <= 0 {
		p.interval = DefaultPollInterval
	}
	if p.logger == nil {
		p.logger = slog.Default()
	}
	return p
}

// Start runs an immediate poll followed by interval polling until ctx is
// cancelled. The returned channel closes after the goroutine exits.
func (p *Poller) Start(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer p.cleanupAll()
		p.pollAndLog(ctx)
		ticker := time.NewTicker(p.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				p.pollAndLog(ctx)
			}
		}
	}()
	return done
}

func (p *Poller) pollAndLog(ctx context.Context) {
	if err := p.Poll(ctx); err != nil {
		p.logger.Error("preview poller: poll failed", "err", err)
	}
}

// Poll performs one deterministic scan of active worker sessions.
func (p *Poller) Poll(ctx context.Context) error {
	if p.source == nil || p.setter == nil {
		return nil
	}
	sessions, err := p.source.ListAllSessions(ctx)
	if err != nil {
		return fmt.Errorf("preview poller list sessions: %w", err)
	}
	activeIDs := make(map[domain.SessionID]struct{}, len(sessions))
	for _, sess := range sessions {
		if sess.IsTerminated {
			if sess.Metadata.WorkspacePath == "" || !pathExists(sess.Metadata.WorkspacePath) {
				continue
			}
		}
		activeIDs[sess.ID] = struct{}{}

		p.ensureDevServer(ctx, sess)

		if sess.IsTerminated {
			continue
		}
		if sess.Kind != domain.KindWorker {
			continue
		}
		storedEntry, workspaceOwned := StoredWorkspaceEntry(sess.Metadata.PreviewURL, sess.ID)
		entry, ok := Entry{}, false
		if workspaceOwned {
			entry, ok = EntryAtPath(sess.Metadata.WorkspacePath, storedEntry)
		}
		if !ok {
			entry, ok = DiscoverEntry(sess.Metadata.WorkspacePath)
		}
		if !ok {
			if workspaceOwned {
				if _, err := p.setter.SetPreview(ctx, sess.ID, ""); err != nil {
					p.logger.Error("preview poller: failed to clear stale preview",
						"session", sess.ID, "err", err)
				}
				p.seen[sess.ID] = entryState{cleared: true}
			}
			continue
		}
		state := stateFor(entry)
		previous, seenBefore := p.seen[sess.ID]
		if seenBefore && previous == state {
			continue
		}
		target, err := FileURL(p.baseURL, sess.ID, entry.Path)
		if err != nil {
			p.logger.Error("preview poller: cannot build isolated preview URL", "session", sess.ID, "err", err)
			p.seen[sess.ID] = state
			continue
		}
		if !p.shouldRefresh(sess, target, seenBefore, workspaceOwned) {
			p.seen[sess.ID] = state
			continue
		}
		if _, err := p.setter.SetPreview(ctx, sess.ID, target); err != nil {
			return fmt.Errorf("preview poller set preview %s: %w", sess.ID, err)
		}
		p.seen[sess.ID] = state
	}
	for id := range p.seen {
		if _, ok := activeIDs[id]; !ok {
			delete(p.seen, id)
		}
	}
	p.mu.Lock()
	for id, server := range p.devServers {
		if _, ok := activeIDs[id]; !ok {
			p.logger.Info("preview: killing dev server for inactive session", "session", id)
			if server.cmd.Process != nil {
				_ = server.cmd.Process.Kill()
			}
			delete(p.devServers, id)
		}
	}
	p.mu.Unlock()
	return nil
}

func (p *Poller) shouldRefresh(sess domain.SessionRecord, target string, seenBefore, workspaceOwned bool) bool {
	current := strings.TrimSpace(sess.Metadata.PreviewURL)
	if current == "" {
		if !seenBefore {
			return sess.Metadata.PreviewRevision == 0
		}
		previous := p.seen[sess.ID]
		return previous.cleared
	}
	if current == target {
		return seenBefore
	}
	return workspaceOwned
}

func stateFor(entry Entry) entryState {
	return entryState{path: entry.Path, modUnix: entry.ModTime.UnixNano(), size: entry.Size}
}

func (p *Poller) ensureDevServer(ctx context.Context, sess domain.SessionRecord) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.devServers[sess.ID]; exists {
		return
	}

	pkgPath, found := findPackageJSON(sess.Metadata.WorkspacePath)
	if !found {
		return
	}

	// Try to symlink node_modules if missing
	targetNodeModules := filepath.Join(filepath.Dir(pkgPath), "node_modules")
	if !pathExists(targetNodeModules) {
		proj, ok, err := p.source.GetProject(ctx, string(sess.ProjectID))
		if err == nil && ok && proj.Path != "" {
			relDir, err := filepath.Rel(sess.Metadata.WorkspacePath, filepath.Dir(pkgPath))
			if err == nil {
				sourceNodeModules := filepath.Join(proj.Path, relDir, "node_modules")
				if pathExists(sourceNodeModules) {
					p.logger.Info("preview: target node_modules missing, linking from main repo", "session", sess.ID, "src", sourceNodeModules, "dest", targetNodeModules)
					if err := linkDirectory(sourceNodeModules, targetNodeModules); err != nil {
						p.logger.Error("preview: failed to link node_modules", "session", sess.ID, "err", err)
					}
				}
			}
		}
	}

	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return
	}
	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return
	}
	devScript, hasDev := pkg.Scripts["dev"]
	if !hasDev {
		devScript, hasDev = pkg.Scripts["start"]
	}
	if !hasDev {
		return
	}

	dir := filepath.Dir(pkgPath)
	p.logger.Info("preview: found dev server package.json", "session", sess.ID, "dir", dir, "script", devScript)

	port := detectDefaultPort(pkgPath, data)
	freePort(port, p.logger)

	cmdName := "npm"
	args := []string{"run", "dev"}
	if !hasDev {
		args = []string{"run", "start"}
	}
	if fileExists(filepath.Join(dir, "pnpm-lock.yaml")) {
		cmdName = "pnpm"
	} else if fileExists(filepath.Join(dir, "yarn.lock")) {
		cmdName = "yarn"
		if !hasDev {
			args = []string{"start"}
		} else {
			args = []string{"dev"}
		}
	}

	if runtime.GOOS == "windows" {
		cmdName += ".cmd"
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PORT="+strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		p.logger.Error("preview: failed to start dev server", "session", sess.ID, "err", err)
		return
	}

	p.logger.Info("preview: started dev server process", "session", sess.ID, "pid", cmd.Process.Pid, "port", port)
	p.devServers[sess.ID] = &devServer{cmd: cmd, port: port}

	go func(id domain.SessionID, c *exec.Cmd) {
		_ = c.Wait()
		p.mu.Lock()
		if current, exists := p.devServers[id]; exists && current.cmd == c {
			p.logger.Info("preview: dev server process exited, cleaning up map", "session", id)
			delete(p.devServers, id)
		}
		p.mu.Unlock()
	}(sess.ID, cmd)

	go func(id domain.SessionID, targetPort int) {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		timeout := time.After(15 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-timeout:
				p.logger.Warn("preview: timeout waiting for dev server port to listen", "session", id, "port", targetPort)
				return
			case <-ticker.C:
				if isPortListening(targetPort) {
					p.logger.Info("preview: dev server port is listening, updating session preview URL", "session", id, "port", targetPort)
					targetURL := fmt.Sprintf("http://localhost:%d", targetPort)
					_, err := p.setter.SetPreview(context.Background(), id, targetURL)
					if err != nil {
						p.logger.Error("preview: failed to set preview URL", "session", id, "err", err)
					}
					return
				}
			}
		}
	}(sess.ID, port)
}

func (p *Poller) cleanupAll() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for id, server := range p.devServers {
		p.logger.Info("preview: cleaning up dev server", "session", id)
		if server.cmd.Process != nil {
			_ = server.cmd.Process.Kill()
		}
	}
	p.devServers = make(map[domain.SessionID]*devServer)
}

type packageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

func findPackageJSON(root string) (string, bool) {
	path := filepath.Join(root, "package.json")
	if fileExists(path) {
		return path, true
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return "", false
	}
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "node_modules" && !strings.HasPrefix(entry.Name(), ".") {
			path = filepath.Join(root, entry.Name(), "package.json")
			if fileExists(path) {
				return path, true
			}
			subEntries, err := os.ReadDir(filepath.Join(root, entry.Name()))
			if err == nil {
				for _, subEntry := range subEntries {
					if subEntry.IsDir() && subEntry.Name() != "node_modules" && !strings.HasPrefix(subEntry.Name(), ".") {
						path = filepath.Join(root, entry.Name(), subEntry.Name(), "package.json")
						if fileExists(path) {
							return path, true
						}
					}
				}
			}
		}
	}
	return "", false
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func linkDirectory(source, target string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "mklink", "/J", target, source)
		return cmd.Run()
	}
	return os.Symlink(source, target)
}

func isPortListening(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 200*time.Millisecond)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func detectDefaultPort(packagePath string, content []byte) int {
	str := string(content)
	if strings.Contains(str, "\"next\"") || strings.Contains(str, "next dev") {
		return 3000
	}
	if strings.Contains(str, "\"vite\"") || strings.Contains(str, "vite") {
		return 5173
	}
	return 3000
}

func freePort(port int, logger *slog.Logger) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("netstat", "-aon").Output()
		if err != nil {
			return
		}
		lines := strings.Split(string(out), "\n")
		target := fmt.Sprintf(":%d", port)
		for _, line := range lines {
			if strings.Contains(line, target) && strings.Contains(line, "LISTENING") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					pidStr := fields[len(fields)-1]
					pid, err := strconv.Atoi(pidStr)
					if err == nil && pid > 0 {
						logger.Info("preview: killing process locking port", "port", port, "pid", pid)
						_ = exec.Command("taskkill", "/F", "/PID", pidStr).Run()
					}
				}
			}
		}
	} else {
		out, err := exec.Command("lsof", "-t", "-i", fmt.Sprintf(":%d", port)).Output()
		if err == nil {
			pids := strings.Split(strings.TrimSpace(string(out)), "\n")
			for _, pidStr := range pids {
				pidStr = strings.TrimSpace(pidStr)
				if pidStr != "" {
					logger.Info("preview: killing process locking port", "port", port, "pid", pidStr)
					_ = exec.Command("kill", "-9", pidStr).Run()
				}
			}
		}
	}
}

package iris

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sriram687/GIT_WIZ/backend/internal/domain"
)

type fakeStore struct {
	sessions []domain.SessionRecord
	project  domain.ProjectRecord
	updates  []domain.SessionRecord
}

func (f *fakeStore) ListAllSessions(ctx context.Context) ([]domain.SessionRecord, error) {
	return f.sessions, nil
}

func (f *fakeStore) GetSession(ctx context.Context, id domain.SessionID) (domain.SessionRecord, bool, error) {
	for _, s := range f.sessions {
		if s.ID == id {
			return s, true, nil
		}
	}
	return domain.SessionRecord{}, false, nil
}

func (f *fakeStore) UpdateSession(ctx context.Context, rec domain.SessionRecord) error {
	f.updates = append(f.updates, rec)
	for i, s := range f.sessions {
		if s.ID == rec.ID {
			f.sessions[i] = rec
		}
	}
	return nil
}

func (f *fakeStore) GetProject(ctx context.Context, id string) (domain.ProjectRecord, bool, error) {
	if f.project.ID == id {
		return f.project, true, nil
	}
	return domain.ProjectRecord{}, false, nil
}

func setupGitRepo(t *testing.T) (string, string) {
	// Create a temp origin repo and worker checkout repo
	tmpDir := t.TempDir()
	originDir := filepath.Join(tmpDir, "origin")
	workerDir := filepath.Join(tmpDir, "worker")

	err := os.MkdirAll(originDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Init origin
	runGit(t, originDir, "init", "--initial-branch=main")
	runGit(t, originDir, "config", "user.name", "Test User")
	runGit(t, originDir, "config", "user.email", "test@example.com")
	runGit(t, originDir, "config", "commit.gpgsign", "false")

	// Create initial commit
	dummyFile := filepath.Join(originDir, "file.txt")
	err = os.WriteFile(dummyFile, []byte("hello\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, originDir, "add", "file.txt")
	runGit(t, originDir, "commit", "-m", "initial commit")

	// Clone to worker
	runGit(t, tmpDir, "clone", originDir, "worker")
	runGit(t, workerDir, "config", "user.name", "Test User")
	runGit(t, workerDir, "config", "user.email", "test@example.com")
	runGit(t, workerDir, "config", "commit.gpgsign", "false")

	// Create a worker branch
	runGit(t, workerDir, "checkout", "-b", "worker-branch")

	return originDir, workerDir
}

func runGit(t *testing.T, dir string, args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed in %s: %v\nOutput: %s", args, dir, err, string(out))
	}
}

func quietLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestIrisObserver_Tick_UpToDate(t *testing.T) {
	_, workerDir := setupGitRepo(t)

	store := &fakeStore{
		sessions: []domain.SessionRecord{
			{
				ID:        "sess-1",
				ProjectID: "proj-1",
				Activity:  domain.Activity{State: domain.ActivityIdle},
				Metadata: domain.SessionMetadata{
					WorkspacePath: workerDir,
					Branch:        "worker-branch",
				},
			},
		},
		project: domain.ProjectRecord{
			ID: "proj-1",
			Config: domain.ProjectConfig{
				DefaultBranch: "main",
			},
		},
	}

	o := New(store, Config{Logger: quietLogger()})
	err := o.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Should not have set behind_upstream
	if store.sessions[0].BehindUpstream {
		t.Fatal("expected behind_upstream to be false")
	}
}

func TestIrisObserver_Tick_BehindCleanMerge(t *testing.T) {
	originDir, workerDir := setupGitRepo(t)

	// Add commit to origin/main
	dummyFile := filepath.Join(originDir, "other.txt")
	err := os.WriteFile(dummyFile, []byte("world\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, originDir, "add", "other.txt")
	runGit(t, originDir, "commit", "-m", "second commit")

	store := &fakeStore{
		sessions: []domain.SessionRecord{
			{
				ID:        "sess-1",
				ProjectID: "proj-1",
				Activity:  domain.Activity{State: domain.ActivityIdle},
				Metadata: domain.SessionMetadata{
					WorkspacePath: workerDir,
					Branch:        "worker-branch",
				},
			},
		},
		project: domain.ProjectRecord{
			ID: "proj-1",
			Config: domain.ProjectConfig{
				DefaultBranch: "main",
			},
		},
	}

	o := New(store, Config{Logger: quietLogger()})
	err = o.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Since session was idle and merge is clean, it should have merged it and kept BehindUpstream false.
	if store.sessions[0].BehindUpstream {
		t.Fatal("expected behind_upstream to remain false after successful clean merge")
	}

	// Verify other.txt exists in worker dir
	if _, err := os.Stat(filepath.Join(workerDir, "other.txt")); os.IsNotExist(err) {
		t.Fatal("expected clean merge to pull other.txt")
	}
}

func TestIrisObserver_Tick_BehindCleanMergeBusy(t *testing.T) {
	originDir, workerDir := setupGitRepo(t)

	// Add commit to origin/main
	dummyFile := filepath.Join(originDir, "other.txt")
	err := os.WriteFile(dummyFile, []byte("world\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, originDir, "add", "other.txt")
	runGit(t, originDir, "commit", "-m", "second commit")

	store := &fakeStore{
		sessions: []domain.SessionRecord{
			{
				ID:        "sess-1",
				ProjectID: "proj-1",
				Activity:  domain.Activity{State: domain.ActivityActive}, // Working/busy!
				Metadata: domain.SessionMetadata{
					WorkspacePath: workerDir,
					Branch:        "worker-branch",
				},
			},
		},
		project: domain.ProjectRecord{
			ID: "proj-1",
			Config: domain.ProjectConfig{
				DefaultBranch: "main",
			},
		},
	}

	o := New(store, Config{Logger: quietLogger()})
	err = o.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Since session is active, it should flag behind_upstream as true instead of merging immediately
	if !store.sessions[0].BehindUpstream {
		t.Fatal("expected behind_upstream to be flagged true since session is busy")
	}

	// Verify other.txt does NOT exist in worker dir yet
	if _, err := os.Stat(filepath.Join(workerDir, "other.txt")); err == nil {
		t.Fatal("expected merge to be postponed")
	}
}

func TestIrisObserver_Tick_BehindConflict(t *testing.T) {
	originDir, workerDir := setupGitRepo(t)

	// Add commit to origin/main editing file.txt
	err := os.WriteFile(filepath.Join(originDir, "file.txt"), []byte("hello from origin\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, originDir, "add", "file.txt")
	runGit(t, originDir, "commit", "-m", "origin edit")

	// Add conflicting commit to worker editing file.txt
	err = os.WriteFile(filepath.Join(workerDir, "file.txt"), []byte("hello from worker\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	runGit(t, workerDir, "add", "file.txt")
	runGit(t, workerDir, "commit", "-m", "worker edit")

	store := &fakeStore{
		sessions: []domain.SessionRecord{
			{
				ID:        "sess-1",
				ProjectID: "proj-1",
				Activity:  domain.Activity{State: domain.ActivityIdle},
				Metadata: domain.SessionMetadata{
					WorkspacePath: workerDir,
					Branch:        "worker-branch",
				},
			},
		},
		project: domain.ProjectRecord{
			ID: "proj-1",
			Config: domain.ProjectConfig{
				DefaultBranch: "main",
			},
		},
	}

	o := New(store, Config{Logger: quietLogger()})
	err = o.Tick(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	// Since it conflicts, it must flag behind_upstream as true
	if !store.sessions[0].BehindUpstream {
		t.Fatal("expected behind_upstream to be flagged true due to merge conflict")
	}
}

package iris

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/sriram687/GIT_WIZ/backend/internal/domain"
	"github.com/sriram687/GIT_WIZ/backend/internal/observe"
	aoprocess "github.com/sriram687/GIT_WIZ/backend/internal/process"
)

const (
	// DefaultTickInterval is the polling frequency for upstream checks.
	DefaultTickInterval = 1 * time.Minute
)

// Store defines the database contract required by the Iris syncer.
type Store interface {
	ListAllSessions(ctx context.Context) ([]domain.SessionRecord, error)
	GetSession(ctx context.Context, id domain.SessionID) (domain.SessionRecord, bool, error)
	UpdateSession(ctx context.Context, rec domain.SessionRecord) error
	GetProject(ctx context.Context, id string) (domain.ProjectRecord, bool, error)
}

// Config holds observer knobs. Zero values use defaults.
type Config struct {
	Tick   time.Duration
	Logger *slog.Logger
}

// Observer is the background agent that fetches and attempts to sync local worktrees.
type Observer struct {
	store  Store
	tick   time.Duration
	logger *slog.Logger
	mu     sync.Mutex
}

// New constructs an Iris syncer observer.
func New(store Store, cfg Config) *Observer {
	o := &Observer{
		store:  store,
		tick:   cfg.Tick,
		logger: cfg.Logger,
	}
	if o.tick <= 0 {
		o.tick = DefaultTickInterval
	}
	if o.logger == nil {
		o.logger = slog.Default()
	}
	return o
}

// Start launches the background synchronization polling loop.
func (o *Observer) Start(ctx context.Context) <-chan struct{} {
	return observe.StartPollLoop(ctx, o.tick, o.Tick, o.logger, "iris observer")
}

// Tick executes a single round of checks and syncs.
func (o *Observer) Tick(ctx context.Context) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	sessions, err := o.store.ListAllSessions(ctx)
	if err != nil {
		return fmt.Errorf("iris: list sessions: %w", err)
	}

	for _, s := range sessions {
		if s.IsTerminated || s.Metadata.WorkspacePath == "" {
			continue
		}

		if err := o.syncSession(ctx, s); err != nil {
			o.logger.Error("iris: sync failed for session", "session_id", s.ID, "err", err)
		}
	}

	return nil
}

func (o *Observer) syncSession(ctx context.Context, s domain.SessionRecord) error {
	defaultBranch := "main"
	proj, ok, err := o.store.GetProject(ctx, string(s.ProjectID))
	if err == nil && ok && proj.Config.DefaultBranch != "" {
		defaultBranch = proj.Config.DefaultBranch
	}

	dir := s.Metadata.WorkspacePath

	// 1. Run git fetch to update remote branches
	fetchCmd := aoprocess.CommandContext(ctx, "git", "-C", dir, "fetch", "origin")
	if err := fetchCmd.Run(); err != nil {
		// Log but proceed (offline mode fallback)
		o.logger.Warn("iris: git fetch failed", "session_id", s.ID, "dir", dir, "err", err)
	}

	// 2. Identify the target ref
	targetRef := "origin/" + defaultBranch
	revCmd := aoprocess.CommandContext(ctx, "git", "-C", dir, "rev-list", "--count", "HEAD.."+targetRef)
	countBytes, err := revCmd.Output()
	if err != nil {
		// Fallback to local default branch
		targetRef = defaultBranch
		revCmd = aoprocess.CommandContext(ctx, "git", "-C", dir, "rev-list", "--count", "HEAD.."+targetRef)
		countBytes, err = revCmd.Output()
		if err != nil {
			return fmt.Errorf("git rev-list check failed: %w", err)
		}
	}

	countStr := strings.TrimSpace(string(countBytes))
	var count int
	if _, err := fmt.Sscanf(countStr, "%d", &count); err != nil {
		return fmt.Errorf("failed to parse revision count: %w", err)
	}

	if count == 0 {
		// Up-to-date: ensure behind_upstream is cleared
		if s.BehindUpstream {
			s.BehindUpstream = false
			if err := o.store.UpdateSession(ctx, s); err != nil {
				return fmt.Errorf("clear behind_upstream state: %w", err)
			}
			o.logger.Info("iris: session is up-to-date now", "session_id", s.ID)
		}
		return nil
	}

	// 3. We are behind. Run trial in-memory merge using git merge-tree
	mergeTreeCmd := aoprocess.CommandContext(ctx, "git", "-C", dir, "merge-tree", "HEAD", targetRef)
	mergeClean := false
	if err := mergeTreeCmd.Run(); err == nil {
		mergeClean = true
	} else if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		mergeClean = false
	} else {
		return fmt.Errorf("git merge-tree check failed: %w", err)
	}

	// 4. Handle merge applicability
	if mergeClean {
		if s.Activity.State == domain.ActivityIdle {
			// Clean & idle: apply merge directly
			mergeCmd := aoprocess.CommandContext(ctx, "git", "-C", dir, "merge", targetRef, "--no-edit", "-m", "Merge upstream changes")
			if err := mergeCmd.Run(); err != nil {
				// Failed to merge (unexpected workspace lock or state)
				s.BehindUpstream = true
				_ = o.store.UpdateSession(ctx, s)
				return fmt.Errorf("git merge failed: %w", err)
			}
			s.BehindUpstream = false
			if err := o.store.UpdateSession(ctx, s); err != nil {
				return fmt.Errorf("update session after clean merge: %w", err)
			}
			o.logger.Info("iris: successfully merged upstream changes cleanly", "session_id", s.ID, "target", targetRef)
		} else {
			// Clean but agent is busy: flag behind_upstream until it goes idle
			if !s.BehindUpstream {
				s.BehindUpstream = true
				if err := o.store.UpdateSession(ctx, s); err != nil {
					return fmt.Errorf("set behind_upstream state (clean/busy): %w", err)
				}
				o.logger.Info("iris: session is behind upstream cleanly but currently busy", "session_id", s.ID)
			}
		}
	} else {
		// Merge conflict: flag behind_upstream
		if !s.BehindUpstream {
			s.BehindUpstream = true
			if err := o.store.UpdateSession(ctx, s); err != nil {
				return fmt.Errorf("set behind_upstream state (conflicted): %w", err)
			}
			o.logger.Warn("iris: session is behind upstream with conflicts", "session_id", s.ID)
		}
	}

	return nil
}

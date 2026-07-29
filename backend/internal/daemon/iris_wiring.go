package daemon

import (
	"context"
	"log/slog"

	irisobserve "github.com/sriram687/GIT_WIZ/backend/internal/observe/iris"
	"github.com/sriram687/GIT_WIZ/backend/internal/storage/sqlite"
)

func startIrisObserver(ctx context.Context, store *sqlite.Store, logger *slog.Logger) <-chan struct{} {
	observer := irisobserve.New(store, irisobserve.Config{Logger: logger})
	return observer.Start(ctx)
}

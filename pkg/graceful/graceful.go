package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WaitForSignalContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		select {
		case <-ctx.Done():
		case <-sigChan:
			cancel()
		}
	}()
	return ctx, cancel
}

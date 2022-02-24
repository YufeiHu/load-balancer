package context_helper

import (
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

const (
    WorkRequestKey = "workRequest"
)

// GetInterruptableCtx creates a child context that gets cancelled
// if SIGTERM/SIGINT signal is caught.
func GetInterruptableCtx(pctx context.Context) (cctx context.Context, cancel context.CancelFunc) {
    cctx, cancel = context.WithCancel(pctx)

    go func() {
        defer cancel()

        c := make(chan os.Signal, 1)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

        select {
        case <-c:
        case <-pctx.Done():
        }
    }()

    return
}

// Delay sleeps for the specific amount or until context cancel/timeout
// and returns context error if context was canceled/timed out.
func Delay(ctx context.Context, duration time.Duration) error {
    select {
    case <-time.After(duration):
    case <-ctx.Done():
    }

    return ctx.Err()
}


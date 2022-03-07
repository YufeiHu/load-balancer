package log_helper

import (
	"context"
	"fmt"
	"load-balancer/internal/helper/context_helper"
	"time"
)

type level string

const (
    infoLevel  = level("INFO")
    errorLevel = level("ERRO")
    warnLevel  = level("WARN")
)

func LogInfo(ctx context.Context, format string, args ...interface{}) {
    log(ctx, infoLevel, format, args...)
}

func LogError(ctx context.Context, format string, args ...interface{}) {
    log(ctx, errorLevel, format, args...)
}

func LogWarn(ctx context.Context, format string, args ...interface{}) {
    log(ctx, warnLevel, format, args...)
}

func log(ctx context.Context, l level, format string, args ...interface{}) {
    workRequestId := context_helper.WorkRequestId(ctx)
    fmt.Printf("[%s] %s [%s]: %s\n", l, time.Now().Format("2006-01-02 15:04:05 MST"), workRequestId, fmt.Sprintf(format, args...))
}

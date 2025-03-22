package log

import (
	"context"
)

type ContextValue int

// Logger type registry.
const ContextLogger ContextValue = 0

func FromContext(ctx context.Context) Logger {
	return ctx.Value(ContextLogger).(Logger)
}

func ToContext(ctx context.Context, logger Logger) context.Context {
	if ctx == nil || logger == nil {
		return ctx
	}
	return context.WithValue(ctx, ContextLogger, logger)
}

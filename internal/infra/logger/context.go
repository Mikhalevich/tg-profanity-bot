package logger

import (
	"context"
)

type contextKey string

const (
	contextLogger = contextKey("contextLogger")
)

func FromContext(ctx context.Context) Logger {
	v := ctx.Value(contextLogger)
	if v == nil {
		return StdLogger()
	}

	l, ok := v.(Logger)
	if !ok {
		return StdLogger()
	}

	return l
}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, contextLogger, l)
}

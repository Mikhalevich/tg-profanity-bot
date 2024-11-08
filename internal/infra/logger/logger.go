package logger

import (
	"context"
)

type Fields map[string]interface{}

//nolint:interfacebloat
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})

	WithContext(ctx context.Context) Logger
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
}

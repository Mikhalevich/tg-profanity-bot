package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logrus struct {
	l *logrus.Entry
}

func NewLogrus() *Logrus {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &Logrus{
		l: logrus.NewEntry(logger),
	}
}

func NewLogrusWithLevel(lvl string) (*Logrus, error) {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	return &Logrus{
		l: logrus.NewEntry(logger),
	}, nil
}

func (lw *Logrus) Debugf(format string, args ...interface{}) {
	lw.l.Debugf(format, args...)
}

func (lw *Logrus) Infof(format string, args ...interface{}) {
	lw.l.Infof(format, args...)
}

func (lw *Logrus) Warnf(format string, args ...interface{}) {
	lw.l.Warnf(format, args...)
}

func (lw *Logrus) Errorf(format string, args ...interface{}) {
	lw.l.Errorf(format, args...)
}

func (lw *Logrus) Debug(args ...interface{}) {
	lw.l.Debug(args...)
}

func (lw *Logrus) Info(args ...interface{}) {
	lw.l.Info(args...)
}

func (lw *Logrus) Warn(args ...interface{}) {
	lw.l.Warn(args...)
}

func (lw *Logrus) Error(args ...interface{}) {
	lw.l.Error(args...)
}

func (lw *Logrus) WithContext(ctx context.Context) Logger {
	return &Logrus{
		l: lw.l.WithContext(ctx),
	}
}

func (lw *Logrus) WithError(err error) Logger {
	return &Logrus{
		l: lw.l.WithError(err),
	}
}

func (lw *Logrus) WithField(key string, value interface{}) Logger {
	return &Logrus{
		l: lw.l.WithField(key, value),
	}
}

func (lw *Logrus) WithFields(fields Fields) Logger {
	return &Logrus{
		l: lw.l.WithFields(logrus.Fields(fields)),
	}
}

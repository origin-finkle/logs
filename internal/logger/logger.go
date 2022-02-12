package logger

import (
	"context"

	"github.com/origin-finkle/logs/internal/version"
	"github.com/sirupsen/logrus"
)

type loggerType string

const (
	Logger loggerType = "internal/logger"
)

func ContextWithLogger(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, Logger, log)
}

func FromContext(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(Logger)
	if logger == nil {
		return logrus.WithContext(ctx).WithField("app_version", version.Version)
	}
	return logger.(*logrus.Entry)
}

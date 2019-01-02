package service

import (
	"context"
	"errors"
)

type tRequestContextKey string

const (
	// ContextKeyLoggerProvider ...
	ContextKeyLoggerProvider tRequestContextKey = "rck-logger-provider"
)

// GetLoggerProviderFromContext ...
func GetLoggerProviderFromContext(ctx context.Context) (LoggerInterface, error) {
	lp, ok := ctx.Value(ContextKeyLoggerProvider).(LoggerInterface)
	if !ok {
		return lp, errors.New("LoggerInterface not found in Context")
	}
	return lp, nil
}

// ContextWithLoggerProvider ...
func ContextWithLoggerProvider(ctx context.Context, lp LoggerInterface) context.Context {
	return context.WithValue(ctx, ContextKeyLoggerProvider, lp)
}

package service

import (
	"context"
	"errors"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
)

type tRequestContextKey string

const (
	// ContextKeyLoggerProvider ...
	ContextKeyLoggerProvider tRequestContextKey = "rck-logger-provider"
	// ContextKeyClient ...
	ContextKeyClient tRequestContextKey = "rck-dogstatsd-metrics"
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

// GetClientFromContext ...
func GetClientFromContext(ctx context.Context) (metrics.Interface, error) {
	dsdi, ok := ctx.Value(ContextKeyClient).(metrics.Interface)
	if !ok {
		return dsdi, errors.New("DogStatsD not found in Context")
	}
	return dsdi, nil
}

// ContextWithClient ...
func ContextWithClient(ctx context.Context, dsdi metrics.Interface) context.Context {
	return context.WithValue(ctx, ContextKeyClient, dsdi)
}

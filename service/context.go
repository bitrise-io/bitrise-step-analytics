package service

import (
	"context"
	"errors"

	"github.com/bitrise-team/bitrise-step-analytics/metrics"
)

type tRequestContextKey string

const (
	// ContextKeyLoggerProvider ...
	ContextKeyLoggerProvider tRequestContextKey = "rck-logger-provider"
	// ContextKeyDogStatsDMetrics ...
	ContextKeyDogStatsDMetrics tRequestContextKey = "rck-dogstatsd-metrics"
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

// GetDogStatsDMetricsFromContext ...
func GetDogStatsDMetricsFromContext(ctx context.Context) (metrics.DogStatsDInterface, error) {
	dsdi, ok := ctx.Value(ContextKeyDogStatsDMetrics).(metrics.DogStatsDInterface)
	if !ok {
		return dsdi, errors.New("DogStatsD not found in Context")
	}
	return dsdi, nil
}

// ContextWithDogStatsDMetrics ...
func ContextWithDogStatsDMetrics(ctx context.Context, dsdi metrics.DogStatsDInterface) context.Context {
	return context.WithValue(ctx, ContextKeyDogStatsDMetrics, dsdi)
}

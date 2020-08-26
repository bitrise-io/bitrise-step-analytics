package service

import (
	"context"
	"errors"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
)

type tRequestContextKey string

const (
	// ContextKeyClient ...
	ContextKeyClient tRequestContextKey = "rck-dogstatsd-metrics"
)

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

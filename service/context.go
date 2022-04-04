package service

import (
	"context"
	"errors"

	"github.com/bitrise-io/bitrise-step-analytics/event"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
)

type tRequestContextKey string

const (
	ContextKeyClient  tRequestContextKey = "rck-dogstatsd-metrics"
	ContextKeyTracker tRequestContextKey = "rck-event-tracker"
)

func GetClientFromContext(ctx context.Context) (metrics.Interface, error) {
	dsdi, ok := ctx.Value(ContextKeyClient).(metrics.Interface)
	if !ok {
		return dsdi, errors.New("DogStatsD not found in Context")
	}
	return dsdi, nil
}

func GetTrackerFromContext(ctx context.Context) (event.Tracker, error) {
	tracker, ok := ctx.Value(ContextKeyTracker).(event.Tracker)
	if !ok {
		return nil, errors.New("event tracker not found in Context")
	}
	return tracker, nil
}

func ContextWithClient(ctx context.Context, dsdi metrics.Interface) context.Context {
	return context.WithValue(ctx, ContextKeyClient, dsdi)
}

func ContextWithTracker(ctx context.Context, tracker event.Tracker) context.Context {
	return context.WithValue(ctx, ContextKeyTracker, tracker)
}

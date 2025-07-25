package service

import (
	"context"
	"errors"

	"github.com/bitrise-io/bitrise-step-analytics/event"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
)

type tRequestContextKey string

const (
	ContextKeyClient  tRequestContextKey = "segment-client"
	ContextKeyTracker tRequestContextKey = "rck-event-tracker"
)

func GetClientFromContext(ctx context.Context) (metrics.Interface, error) {
	segmentClient, ok := ctx.Value(ContextKeyClient).(metrics.Interface)
	if !ok {
		return segmentClient, errors.New("segment client not found in Context")
	}
	return segmentClient, nil
}

func GetTrackerFromContext(ctx context.Context) (event.Tracker, error) {
	tracker, ok := ctx.Value(ContextKeyTracker).(event.Tracker)
	if !ok {
		return nil, errors.New("event tracker not found in Context")
	}
	return tracker, nil
}

func ContextWithClient(ctx context.Context, segmentClient metrics.Interface) context.Context {
	return context.WithValue(ctx, ContextKeyClient, segmentClient)
}

func ContextWithTracker(ctx context.Context, tracker event.Tracker) context.Context {
	return context.WithValue(ctx, ContextKeyTracker, tracker)
}

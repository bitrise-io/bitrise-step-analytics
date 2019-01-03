package service_test

import "github.com/bitrise-team/bitrise-step-analytics/metrics"

type testDogStatsdMetrics struct{}

func (m *testDogStatsdMetrics) Track(t metrics.Trackable, metricName string) {}

func (m *testDogStatsdMetrics) Close() {}

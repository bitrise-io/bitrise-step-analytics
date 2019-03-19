package service_test

import "github.com/bitrise-io/bitrise-step-analytics/metrics"

type testDogStatsdMetrics struct{}

func (m *testDogStatsdMetrics) Track(t metrics.Trackable, metricName string) {}

func (m *testDogStatsdMetrics) Close() {}

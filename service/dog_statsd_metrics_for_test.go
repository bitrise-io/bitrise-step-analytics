package service_test

import "github.com/bitrise-io/bitrise-step-analytics/metrics"

type testClient struct{}

func (m *testClient) Track(t metrics.Trackable, metricName string) {}

func (m *testClient) Close() {}

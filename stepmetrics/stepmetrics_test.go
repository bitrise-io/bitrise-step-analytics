package stepmetrics

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_normalizeEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		want     string
	}{
		{
			endpoint: "/v1/certificates",
			want:     "/v1/certificates",
		},
		{
			endpoint: "/v1/devices",
			want:     "/v1/devices",
		},
		{
			endpoint: "/v1/profiles",
			want:     "/v1/profiles",
		},
		{
			endpoint: "/v1/bundleids",
			want:     "/v1/bundleids",
		},
		{
			endpoint: "/v1/profiles/qx2rdp4p9r",
			want:     "/v1/profiles/{id}",
		},
		{
			endpoint: "/v1/profiles/22x3m5y2l3/bundleid",
			want:     "/v1/profiles/{id}/bundleid",
		},
		{
			endpoint: "/v1/profiles/22x3m5y2l3/certificates",
			want:     "/v1/profiles/{id}/certificates",
		},
		{
			endpoint: "/v1/bundleids/22x3m5y2l3/bundleidcapabilities",
			want:     "/v1/bundleids/{id}/bundleidcapabilities",
		},
		{
			endpoint: "/v1/certificates/abc123def456",
			want:     "/v1/certificates/{id}",
		},
		{
			endpoint: "/v1/devices/xyz789abc123",
			want:     "/v1/devices/{id}",
		},
		{
			endpoint: "/v1/certificates/1234567890",
			want:     "/v1/certificates/{id}",
		},
		{
			endpoint: "/v2/profiles/someid",
			want:     "/v2/profiles/someid",
		},
		{
			endpoint: "/v1/bundleids/verylongid123456/bundleidcapabilities",
			want:     "/v1/bundleids/{id}/bundleidcapabilities",
		},
		{
			endpoint: "",
			want:     "",
		},
		{
			endpoint: "/v1/profiles/abc-123_def456",
			want:     "/v1/profiles/{id}",
		},
	}

	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt.want, "/", "_"), func(t *testing.T) {
			got := normalizeEndpoint(tt.endpoint)
			require.Equal(t, tt.want, got)
		})
	}
}

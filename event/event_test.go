package event

import (
	"testing"
)

func Test_convertEpochToBigQueryTimestampFormat(t *testing.T) {

	tests := []struct {
		name      string
		timestamp int64
		want      string
	}{
		{name: "Test proper conversion", timestamp: 1644253819659123, want: "2022-02-07 17:10:19.659123"},
		{name: "Zero timestamp", timestamp: 0, want: "1970-01-01 00:00:00.000000"},
		{name: "Negative timestamp", timestamp: -1, want: "1969-12-31 23:59:59.999999"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertEpochInMicrosecondsToBigQueryTimestampFormat(tt.timestamp)
			if got != tt.want {
				t.Errorf("convertEpochInMicrosecondsToBigQueryTimestampFormat() got = %v, want %v", got, tt.want)
			}
		})
	}
}

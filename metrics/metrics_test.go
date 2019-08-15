package metrics

import (
	"reflect"
	"testing"
)

func Test_parseTrackableFields(t *testing.T) {
	tests := []struct {
		name  string
		model interface{}
		want  map[string]interface{}
	}{
		{
			name: "track tag check",
			model: struct {
				Field1 string `track:"Field1Test"`
				Field2 bool   `track:"Field 2 bool test"`
			}{
				Field1: "test value",
				Field2: true,
			},
			want: map[string]interface{}{
				"Field1Test":        "test value",
				"Field 2 bool test": true,
			},
		},
		{
			name: "track tag check - one field has no track tag",
			model: struct {
				Field1          string `track:"Field1Test"`
				FieldNoTrackTag int    `json:"Field 2 bool test"`
				Field2          bool   `track:"Field 2 bool test"`
			}{
				Field1:          "test value",
				FieldNoTrackTag: 99,
				Field2:          true,
			},
			want: map[string]interface{}{
				"Field1Test":        "test value",
				"Field 2 bool test": true,
			},
		},
		{
			name: "track tag check - multiple tags",
			model: struct {
				Field1          string `json:"json tag" track:"Field1Test"`
				FieldNoTrackTag int    `json:"Field 2 bool test"`
				Field2          bool   `track:"Field 2 bool test"`
			}{
				Field1:          "test value",
				FieldNoTrackTag: 99,
				Field2:          true,
			},
			want: map[string]interface{}{
				"Field1Test":        "test value",
				"Field 2 bool test": true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTrackableFields(tt.model); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTrackableFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

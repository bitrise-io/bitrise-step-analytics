package models

import "fmt"

type TrackEvent struct {
	ID         string         `json:"id"`
	EventName  string         `json:"event_name"`
	Timestamp  int64          `json:"timestamp"`
	Properties map[string]any `json:"properties"`
}

func (e TrackEvent) StringProp(key string) string {
	if val, ok := e.Properties[key].(string); ok {
		return val
	}
	return ""
}

func (e TrackEvent) IntProp(key string) int {
	// event.Properties is untyped and it's unmarshalled from a JSON payload.
	// JSON numbers are unmarshaled as float64 because the JSON spec only has a `number` type
	if val, ok := e.Properties[key].(float64); ok {
		return int(val)
	}
	// Fallback for cases where it might actually be an int
	if val, ok := e.Properties[key].(int); ok {
		return val
	}
	return 0
}

func (e TrackEvent) BoolProp(key string) (bool, error) {
	if val, ok := e.Properties[key].(bool); ok {
		return val, nil
	}
	return false, fmt.Errorf("property %q is not a bool", key)
}

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

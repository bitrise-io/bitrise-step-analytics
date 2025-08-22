package models

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

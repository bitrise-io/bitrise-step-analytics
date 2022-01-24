package models

// TrackEvent ...
type TrackEvent struct {
	ID         string                 `json:"id"`
	EventName  string                 `json:"event_name"`
	Timestamp  int64                  `json:"timestamp"`
	Properties map[string]interface{} `json:"properties"`
}

package models

// RemoteLog ...
type RemoteLog struct {
	LogLevel string                 `json:"log_level"`
	Message  string                 `json:"message"`
	Data     map[string]interface{} `json:"data"`
}

// Event ...
func (a RemoteLog) Event() string {
	return "remote_log"
}

// Model ...
func (a RemoteLog) Model() interface{} {
	return a
}

// UserID ...
func (a RemoteLog) UserID() string {
	stepID, _ := a.Data["step_id"].(string)
	return stepID
}

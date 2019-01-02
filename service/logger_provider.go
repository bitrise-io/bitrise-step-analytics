package service

import "go.uber.org/zap"

// LoggerInterface ...
type LoggerInterface interface {
	GetLogger() *zap.Logger
}

// LoggerProvider ...
type LoggerProvider struct {
	logger *zap.Logger
}

// NewLoggerProvider ...
func NewLoggerProvider(l *zap.Logger) *LoggerProvider {
	return &LoggerProvider{
		logger: l,
	}
}

// GetLogger ...
func (l *LoggerProvider) GetLogger() *zap.Logger {
	return l.logger
}

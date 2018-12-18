package utils

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetLogger ...
func GetLogger() (*zap.Logger, func() error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"./log/steps.log",
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, func() error {
			return errors.WithStack(err)
		}
	}
	return logger, logger.Sync
}

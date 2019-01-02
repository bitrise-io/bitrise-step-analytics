package utils

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// GetLogger ...
func GetLogger() (*zap.Logger, func() error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, func() error {
			return errors.WithStack(err)
		}
	}
	return logger, logger.Sync
}

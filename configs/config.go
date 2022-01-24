package configs

import (
	"errors"
	"os"
)

// ConfigModel ...
type ConfigModel struct {
	Port, EnvMode, SegmentWriteKey, TrackerWriteKey string
}

// Validate ...
func (c ConfigModel) Validate() error {
	if len(c.Port) < 1 {
		return errors.New("Port must be specified")
	}
	if len(c.EnvMode) < 1 {
		return errors.New("Env mode must be specified")
	}
	if len(c.SegmentWriteKey) < 1 {
		return errors.New("Segment write key must be specified")
	}
	return nil
}

func createFromEnvs() (ConfigModel, error) {
	c := ConfigModel{
		Port:            os.Getenv("PORT"),
		EnvMode:         os.Getenv("GO_ENV"),
		SegmentWriteKey: os.Getenv("SEGMENT_WRITE_KEY"),
		TrackerWriteKey: os.Getenv("TRACKER_WRITE_KEY"),
	}
	return c, nil
}

// CreateAndValidate - reads the config from envs etc.
func CreateAndValidate() (ConfigModel, error) {
	conf, err := createFromEnvs()
	if err != nil {
		return conf, err
	}

	return conf, conf.Validate()
}

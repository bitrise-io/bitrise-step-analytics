package configs

import (
	"errors"
	"os"
)

// ConfigModel ...
type ConfigModel struct {
	Port string
}

// Validate ...
func (c ConfigModel) Validate() error {
	if len(c.Port) < 1 {
		return errors.New("Port must be specified")
	}
	return nil
}

func createFromEnvs() (ConfigModel, error) {
	c := ConfigModel{
		Port: os.Getenv("PORT"),
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

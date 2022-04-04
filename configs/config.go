package configs

import (
	"github.com/bitrise-io/go-steputils/v2/stepconf"
)

type Config struct {
	Port              string `env:"PORT,required"`
	EnvMode           string `env:"GO_ENV,required"`
	SegmentWriteKey   string `env:"SEGMENT_WRITE_KEY,required"`
	PubSubCredentials string `env:"PUBSUB_CREDENTIALS,required"`
	PubSubProject     string `env:"PUBSUB_PROJECT,required"`
	PubSubTopic       string `env:"PUBSUB_TOPIC,required"`
}

// Parse - reads the config from envs etc.
func Parse(parser stepconf.InputParser) (Config, error) {
	var cfg Config
	if err := parser.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

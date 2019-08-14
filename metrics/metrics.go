package metrics

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	segment "gopkg.in/segmentio/analytics-go.v3"
)

// Interface ...
type Interface interface {
	Track(t Trackable)
	Close()
}

// Client ...
type Client struct {
	client segment.Client
}

// Trackable defines a configuration of a
// trackable piece of the execution stack
// It's used to track supervisor proccess stacks
type Trackable interface {
	Model() interface{} // returns a struct which has `track` field tags defined
	Event() string
	UserID() string
}

// NewClient ...
func NewClient(segmentWriteKey string) *Client {
	return &Client{
		client: segment.New(segmentWriteKey),
	}
}

// Track ...
func (b *Client) Track(t Trackable) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Failed to sync logger: %s", err)
		}
	}()

	if err := b.client.Enqueue(segment.Track{
		Event:      t.Event(),
		UserId:     t.UserID(),
		Properties: parseTrackableFields(t.Model()),
	}); err != nil {
		logger.Error("DogStatsD Diagnostic backend has failed to track",
			zap.String("profile_name", t.UserID()),
			zap.Any("error_details", errors.WithStack(err)),
		)
	}
}

// Close ...
func (b *Client) Close() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Failed to sync logger: %s", err)
		}
	}()

	if err := b.client.Close(); err != nil {
		logger.Error("DogStatsD Diagnostic backend has failed to close its client",
			zap.Any("error_details", errors.WithStack(err)),
		)
	}
}

func parseTrackableFields(model interface{}) map[string]interface{} {
	properties := map[string]interface{}{}
	for i, v := 0, reflect.ValueOf(model); i < v.NumField(); i++ {
		if key, ok := v.Type().Field(i).Tag.Lookup("track"); ok {
			properties[key] = v.Field(i).Interface()
		}
	}
	return properties
}

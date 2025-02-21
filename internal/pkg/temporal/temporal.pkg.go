package temporal

import (
	"context"

	"go.temporal.io/sdk/client"
)

var (
	TemporalClient client.Client
)

// InitTemporalConnection initializes connection to temporal server
func InitTemporalConnection(ctx context.Context) error {
	var err error

	// Create temporal client
	TemporalClient, err = client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})

	if err != nil {
		return err
	}

	return nil
}

// GetTemporalClient returns temporal client instance
func GetTemporalClient() client.Client {
	return TemporalClient
}

// CloseTemporalConnection closes the temporal client connection
func CloseTemporalConnection() {
	if TemporalClient != nil {
		TemporalClient.Close()
	}
}

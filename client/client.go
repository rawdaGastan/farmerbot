// Package client is responsible for client communication
package client

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/threefoldtech/zbus"
)

// FarmerbotClient for interacting with the farmerbot
type FarmerbotClient struct {
	zBusClient zbus.Client
	version    zbus.Version
}

// NewFarmerClient creates a new client
func NewFarmerClient(zBusClient zbus.Client) *FarmerbotClient {
	return &FarmerbotClient{
		zBusClient: zBusClient,
		version:    "1.0.0",
	}
}

// Call calls a specific cmd with its date and returns the result
func (f *FarmerbotClient) Call(ctx context.Context, fn string, data []interface{}, result interface{}) error {

	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize request body with error: %v", err)
	}

	cmd := strings.Split(fn, ".")
	if len(cmd) != 3 {
		return fmt.Errorf("invalid command length %s", cmd)
	}

	if cmd[0] != "farmerbot" {
		return fmt.Errorf("invalid command parent %s", cmd)
	}

	if cmd[1] != "farmmanager" && cmd[1] != "powermanager" && cmd[1] != "nodemanager" {
		return fmt.Errorf("invalid command sub parent %s", cmd)
	}

	output, err := f.zBusClient.RequestContext(ctx, cmd[0], zbus.ObjectID{
		Name:    cmd[1],
		Version: f.version,
	}, cmd[2], data...)

	if err != nil {
		return fmt.Errorf("invalid request %v, with data: %v", cmd, payload)
	}

	if output.Output.Error != nil {
		return output.Output.Error
	}

	if len(output.Output.Data) == 0 {
		return nil
	}

	loader := zbus.Loader{
		&result,
	}
	return output.Unmarshal(&loader)
}

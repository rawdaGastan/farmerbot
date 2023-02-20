// Package constants for farmerbot constants
package constants

import "time"

const (
	//TimeoutPowerStateChange a timeout for changing nodes power
	TimeoutPowerStateChange = time.Minute * 30

	//TimeoutClaimedResources a number of rounds to delay claims updates
	TimeoutClaimedResources = 6

	//DefaultWakeUpThreshold default threshold to wake up a new node
	DefaultWakeUpThreshold = uint64(80)
	//MinWakeUpThreshold min threshold to wake up a new node
	MinWakeUpThreshold = uint64(50)
	//MaxWakeUpThreshold max threshold to wake up a new node
	MaxWakeUpThreshold = uint64(80)
)

const (
	mainNetwork string = "main"
	testNetwork string = "test"
	devNetwork  string = "dev"
	qaNetwork   string = "qa"
)

// SubstrateURLs for substrate urls
var SubstrateURLs = map[string][]string{
	testNetwork: {"wss://tfchain.test.grid.tf/ws", "wss://tfchain.test.grid.tf:443"},
	mainNetwork: {"wss://tfchain.grid.tf/ws", "wss://tfchain.grid.tf:443"},
	devNetwork:  {"wss://tfchain.dev.grid.tf/ws", "wss://tfchain.dev.grid.tf:443"},
	qaNetwork:   {"wss://tfchain.qa.grid.tf/ws", "wss://tfchain.qa.grid.tf:443"},
}

// RelayURLS relay urls
var RelayURLS = map[string]string{
	devNetwork:  "wss://relay.dev.grid.tf",
	testNetwork: "wss://relay.test.grid.tf",
	qaNetwork:   "wss://relay.qa.grid.tf",
	mainNetwork: "wss://relay.grid.tf",
}

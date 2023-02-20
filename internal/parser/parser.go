// Package parser for parsing cmd configs
package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog/log"
)

// ReadFile reads a file and returns its contents
func ReadFile(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

// ParseJSONIntoConfig parses the JSON configuration
func ParseJSONIntoConfig(content []byte) (models.Config, error) {
	c := models.Config{}

	err := json.Unmarshal(content, &c)
	if err != nil {
		return models.Config{}, err
	}

	// default values
	for i := range c.Nodes {
		if c.Nodes[i].Resources.OverProvisionCPU == 0 {
			c.Nodes[i].Resources.OverProvisionCPU = 1
		}
		if c.Nodes[i].Resources.OverProvisionCPU < 1 || c.Nodes[i].Resources.OverProvisionCPU > 4 {
			return models.Config{}, fmt.Errorf("overProvision cpu should be a value between 1 and 4 not %v", c.Nodes[i].Resources.OverProvisionCPU)
		}

		c.Nodes[i].PowerState.ON = true
	}

	if c.Power.WakeUpThreshold == 0 {
		c.Power.WakeUpThreshold = constants.DefaultWakeUpThreshold
	}

	if c.Power.WakeUpThreshold < constants.MinWakeUpThreshold {
		log.Warn().Msgf("setting wakeUpThreshold should be in the range [%d, %d] not %d", constants.MinWakeUpThreshold, constants.MaxWakeUpThreshold, c.Power.WakeUpThreshold)
		c.Power.WakeUpThreshold = constants.MinWakeUpThreshold
	}

	if c.Power.WakeUpThreshold > constants.MaxWakeUpThreshold {
		log.Warn().Msgf("setting wakeUpThreshold should be in the range [%d, %d] not %d", constants.MinWakeUpThreshold, constants.MaxWakeUpThreshold, c.Power.WakeUpThreshold)
		c.Power.WakeUpThreshold = constants.MaxWakeUpThreshold
	}

	c.Power.PeriodicWakeup = models.WakeupDate(c.Power.PeriodicWakeup.PeriodicWakeupStart())

	// required values for farm
	if c.Farm.ID == 0 {
		return c, errors.New("farm ID is required")
	}

	// required values for node
	for i, n := range c.Nodes {
		if n.ID == 0 {
			return c, fmt.Errorf("node ID with index %d is required", i)
		}
		if n.TwinID == 0 {
			return c, fmt.Errorf("node twin ID with index %d is required", i)
		}
		if n.Resources.Total.SRU == 0 {
			return c, fmt.Errorf("node total SRU with index %d is required", i)
		}
		if n.Resources.Total.CRU == 0 {
			return c, fmt.Errorf("node total CRU with index %d is required", i)
		}
		if n.Resources.Total.MRU == 0 {
			return c, fmt.Errorf("node total MRU with index %d is required", i)
		}
		if n.Resources.Total.HRU == 0 {
			return c, fmt.Errorf("node total HRU with index %d is required", i)
		}
	}

	return c, nil
}

// ParseJSONIntoFarm parses JSON into farm
func ParseJSONIntoFarm(content []byte) (models.Farm, error) {
	farm := models.Farm{}

	err := json.Unmarshal(content, &farm)
	if err != nil {
		return models.Farm{}, err
	}

	// required values for farm
	if farm.ID == 0 {
		return models.Farm{}, errors.New("farm ID is required")
	}

	return farm, nil
}

// ParseJSONIntoPower parses JSON into power
func ParseJSONIntoPower(content []byte) (models.Power, error) {
	power := models.Power{}

	err := json.Unmarshal(content, &power)
	if err != nil {
		return models.Power{}, err
	}

	if power.WakeUpThreshold == 0 {
		power.WakeUpThreshold = constants.DefaultWakeUpThreshold
	}

	if power.WakeUpThreshold < constants.MinWakeUpThreshold {
		log.Warn().Msgf("setting wakeUpThreshold should be in the range [%d, %d] not %d", constants.MinWakeUpThreshold, constants.MaxWakeUpThreshold, power.WakeUpThreshold)
		power.WakeUpThreshold = constants.MinWakeUpThreshold
	}

	if power.WakeUpThreshold > constants.MaxWakeUpThreshold {
		log.Warn().Msgf("setting wakeUpThreshold should be in the range [%d, %d] not %d", constants.MinWakeUpThreshold, constants.MaxWakeUpThreshold, power.WakeUpThreshold)
		power.WakeUpThreshold = constants.MaxWakeUpThreshold
	}

	power.PeriodicWakeup = models.WakeupDate(power.PeriodicWakeup.PeriodicWakeupStart())

	return power, nil
}

// ParseJSONIntoNode parses JSON into node
func ParseJSONIntoNode(content []byte) (models.Node, error) {
	node := models.Node{}

	err := json.Unmarshal(content, &node)
	if err != nil {
		return models.Node{}, err
	}

	// default values
	if node.Resources.OverProvisionCPU == 0 {
		node.Resources.OverProvisionCPU = 1
	}
	if node.Resources.OverProvisionCPU < 1 || node.Resources.OverProvisionCPU > 4 {
		return models.Node{}, fmt.Errorf("overProvision cpu should be a value between 1 and 4 not %v", node.Resources.OverProvisionCPU)
	}

	node.PowerState.ON = true

	// required values for node
	if node.ID == 0 {
		return models.Node{}, fmt.Errorf("node %d ID  is required", node.ID)
	}
	if node.TwinID == 0 {
		return models.Node{}, fmt.Errorf("node %d twin ID is required", node.ID)
	}
	if node.Resources.Total.SRU == 0 {
		return models.Node{}, fmt.Errorf("node %d total SRU is required", node.ID)
	}
	if node.Resources.Total.CRU == 0 {
		return models.Node{}, fmt.Errorf("node %d total CRU is required", node.ID)
	}
	if node.Resources.Total.MRU == 0 {
		return models.Node{}, fmt.Errorf("node %d total MRU is required", node.ID)
	}
	if node.Resources.Total.HRU == 0 {
		return models.Node{}, fmt.Errorf("node %d total HRU is required", node.ID)
	}

	return node, nil
}

// ParseJSONIntoNodeOptions parses JSON into node options
func ParseJSONIntoNodeOptions(content []byte) (models.NodeOptions, error) {
	options := models.NodeOptions{}

	err := json.Unmarshal(content, &options)
	if err != nil {
		return models.NodeOptions{}, err
	}

	return options, nil
}

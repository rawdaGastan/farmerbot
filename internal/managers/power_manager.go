// Package manager provides how to manage nodes, farms and power
package manager

import (
	"fmt"
	"time"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/substrate-client"
)

// PowerHandler interface for mocks
type PowerHandler interface {
	Configure(power models.Power) error
	PowerOn(node models.Node) error
	PowerOff(node models.Node) error

	PeriodicWakeup()
	PowerManagement()
}

// PowerManager manages the power of nodes
type PowerManager struct {
	logger   zerolog.Logger
	db       models.RedisDB
	identity substrate.Identity
	subConn  models.Sub
}

// NewPowerManager creates a new PowerManager
func NewPowerManager(network string, mnemonics string, address string, logger zerolog.Logger) (PowerManager, error) {
	substrateManager := substrate.NewManager(constants.SubstrateURLs[network]...)
	subConn, err := substrateManager.Substrate()
	if err != nil {
		return PowerManager{}, err
	}

	identity, err := substrate.NewIdentityFromSr25519Phrase(mnemonics)
	if err != nil {
		return PowerManager{}, err
	}

	return PowerManager{logger, models.NewRedisDB(address), identity, subConn}, nil
}

// Configure configure a power
func (p *PowerManager) Configure(jsonContent []byte) error {
	power, err := parser.ParseJSONIntoPower(jsonContent)
	if err != nil {
		return fmt.Errorf("failed to get power from json content %v", err)
	}

	p.logger.Debug().Msgf("power configuration threshold is %v, wake up time is %v", power.WakeUpThreshold, time.Time(power.PeriodicWakeup))
	return p.db.SetPower(power)
}

// PowerOn sets the node power state ON
func (p *PowerManager) PowerOn(nodeID uint32) error {
	p.logger.Info().Msgf("POWER ON: %d", nodeID)

	node, err := p.db.GetNode(nodeID)
	if err != nil {
		return err
	}

	if err := node.SetNodePower(p.identity, p.subConn, true); err != nil {
		return err
	}

	return p.db.UpdatesNodes(node)
}

// PowerOff sets the node power state OFF
func (p *PowerManager) PowerOff(nodeID uint32) error {
	p.logger.Info().Msgf("POWER OFF: %d", nodeID)

	onNodes, err := p.db.FilterOnNodes()
	if err != nil {
		return err
	}

	if len(onNodes) < 2 {
		return fmt.Errorf("cannot power off node %d, at least one node should be on in the farm", nodeID)
	}

	node, err := p.db.GetNode(nodeID)
	if err != nil {
		return err
	}

	if err := node.SetNodePower(p.identity, p.subConn, false); err != nil {
		return err
	}

	return p.db.UpdatesNodes(node)
}

func (p *PowerManager) PeriodicWakeup() error {
	nodes, err := p.db.GetNodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes from db with error: %v", err)
	}

	power, err := p.db.GetPower()
	if err != nil {
		return fmt.Errorf("failed to get power from db with error: %v", err)
	}

	now := time.Now()
	periodicWakeupStart := power.PeriodicWakeup.PeriodicWakeupStart()
	p.logger.Debug().Msgf("periodic wakeup time is %v", periodicWakeupStart)

	if periodicWakeupStart.Before(now) {
		for _, node := range nodes {
			if node.PowerState.OFF && node.LastTimeAwake.Before(periodicWakeupStart) {
				if err := p.PowerOn(node.ID); err != nil {
					return fmt.Errorf("power on node %d failed with error: %v", node.ID, err)
				}
				// reboot one at a time others will be rebooted 5 min later
				break

			}
		}
	}

	return nil
}

func (p *PowerManager) PowerManagement() error {
	nodes, err := p.db.GetNodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes from db with error: %v", err)
	}

	power, err := p.db.GetPower()
	if err != nil {
		return fmt.Errorf("failed to get power from db with error: %v", err)
	}

	if len(models.FilterWakingOrShuttingNodes(nodes)) > 0 {
		// in case one of the nodes is waking up or shutting down do nothing until the timeouts occur or the nodes are up or down.
		return nil
	}

	usedResources, totalResources, err := p.calculateResourceUsage()
	if err != nil {
		return err
	}
	if totalResources == 0 {
		return nil
	}

	// usage > threshold
	resourceUsage := 100 * usedResources / totalResources
	if resourceUsage >= power.WakeUpThreshold {
		sleepingNodes := models.FilterOffNodes(nodes)
		if len(sleepingNodes) > 0 {
			node := sleepingNodes[0]
			p.logger.Debug().Msgf("too much resource usage: %d. Turning on node %d", resourceUsage, node.ID)
			if err := p.PowerOn(node.ID); err != nil {
				return fmt.Errorf("power on node %d failed with error: %v", node.ID, err)
			}
		}
	} else {
		unusedNodes := models.FilterUnusedOnNodes(nodes)
		if len(unusedNodes) > 1 {
			// shutdown a node if there is more then 1 unused node (aka keep at least one node online)
			newUsedResources := usedResources
			newTotalResources := totalResources
			nodesLeftOnline := len(unusedNodes)
			for _, node := range unusedNodes {
				// check that we have at least one unused node left online
				if nodesLeftOnline == 1 {
					break
				}
				// nodes with public config can't be shutdown
				if node.PublicConfig {
					continue
				}

				nodesLeftOnline -= 1
				newUsedResources -= node.Resources.Used.HRU + node.Resources.Used.SRU + node.Resources.Used.MRU + node.Resources.Used.CRU
				newTotalResources -= node.Resources.Total.HRU + node.Resources.Total.SRU + node.Resources.Total.MRU + node.Resources.Total.CRU
				if 100*newUsedResources/newTotalResources < power.WakeUpThreshold {
					// we need to keep the resource percentage lower than the threshold
					p.logger.Debug().Msgf("too low resource usage: %d. Turning off unused node %d", resourceUsage, node.ID)
					if err := p.PowerOff(node.ID); err != nil {
						return fmt.Errorf("power off node %d failed with error: %v", node.ID, err)
					}
				}
			}
		} else {
			p.logger.Debug().Msg("nothing to shutdown")
		}
	}
	return nil
}

func (p *PowerManager) calculateResourceUsage() (uint64, uint64, error) {
	nodes, err := p.db.GetNodes()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get nodes from db with error: %v", err)
	}

	usedResources := models.Capacity{}
	totalResources := models.Capacity{}

	for _, node := range nodes {
		if node.PowerState.ON {
			usedResources.Add(node.Resources.Used)
			totalResources.Add(node.Resources.Total)
		}
	}

	used := usedResources.CRU + usedResources.HRU + usedResources.MRU + usedResources.SRU
	total := totalResources.CRU + totalResources.HRU + totalResources.MRU + totalResources.SRU

	return used, total, nil
}

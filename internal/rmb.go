// Package internal for farmerbot internals
package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/rmb-sdk-go"
	"github.com/threefoldtech/rmb-sdk-go/direct"
	"github.com/threefoldtech/substrate-client"
	"github.com/threefoldtech/zos/pkg"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

// RMBClient is an rmb abstract client interface.
type RMBClient interface {
	Call(ctx context.Context, twin uint32, fn string, data interface{}, result interface{}) error
}

type rmbNodeClient struct {
	logger zerolog.Logger
	rmb    rmb.Client //RMBClient
	sub    models.Sub
}

func newRmbNodeClient(sub *substrate.Substrate, mnemonics string, network string, logger zerolog.Logger) (rmbNodeClient, error) {
	sessionID := fmt.Sprintf("tf-%d", os.Getpid())
	rmbClient, err := direct.NewClient("sr25519", mnemonics, constants.RelayURLS[network], sessionID, sub)
	if err != nil {
		return rmbNodeClient{}, fmt.Errorf("failed with error: %w, couldn't create rmb client", err)
	}

	return rmbNodeClient{
		rmb: rmbClient,
	}, nil
}

// PingNode checks state of the node
func (n *rmbNodeClient) pingNode(ctx context.Context, node models.Node) (bool, error) {
	var err error
	if err = n.systemVersion(ctx, node.TwinID); err != nil {
		if node.PowerState.WakingUp {
			if time.Since(node.LastTimePowerStateChanged) < constants.TimeoutPowerStateChange {
				n.logger.Debug().Msgf("Node %d is waking up.", node.ID)
				return false, nil
			}
			err = fmt.Errorf("node %d wakeup was unsuccessful. putting its state back to off", node.ID)
		}

		if node.PowerState.ShuttingDown {
			n.logger.Debug().Msgf("Node %d shutting down was successful", node.ID)
		}

		if node.PowerState.ON {
			err = fmt.Errorf("node %d is not responding while we expect it to", node.ID)
		}

		if node.PowerState.OFF {
			n.logger.Debug().Msgf("Node %d is offline.", node.ID)
		}

		node.PowerState.OFF = true
		node.PowerState.ON = !node.PowerState.OFF
		node.PowerState.WakingUp = !node.PowerState.OFF
		node.PowerState.ShuttingDown = !node.PowerState.OFF
		node.LastTimePowerStateChanged = time.Now()
		return false, err
	}

	if node.PowerState.ShuttingDown {
		if time.Since(node.LastTimePowerStateChanged) < constants.TimeoutPowerStateChange {
			n.logger.Debug().Msgf("Node %d is shutting down.", node.ID)
			return false, nil
		}
		err = fmt.Errorf("node %d shutting down was unsuccessful. putting its state back to on", node.ID)
	} else {
		n.logger.Debug().Msgf("Node %d is online.", node.ID)
	}

	node.PowerState.ON = true
	node.PowerState.OFF = !node.PowerState.OFF
	node.PowerState.ShuttingDown = !node.PowerState.OFF
	node.PowerState.WakingUp = !node.PowerState.OFF
	node.LastTimePowerStateChanged = time.Now()
	node.LastTimeAwake = time.Now()
	return true, err
}

// UpdateNode updates the node statistics
func (n *rmbNodeClient) updateNode(ctx context.Context, node models.Node) error {
	if node.TimeoutClaimedResources.Before(time.Now()) {
		stats, err := n.statistics(ctx, node.TwinID)
		if err != nil {
			return fmt.Errorf("failed to get statistics of node %d with error: %w", node.ID, err)
		}
		node.UpdateResources(stats)

		pools, err := n.getStoragePools(ctx, node.TwinID)
		if err != nil {
			return fmt.Errorf("failed to update storage pools of node %d with error: %w", node.ID, err)
		}
		node.Pools = pools

		rentContract, err := n.sub.GetNodeRentContract(node.ID)
		if err != nil {
			return fmt.Errorf("failed to update contracts of node %d with error: %w", node.ID, err)
		}

		if rentContract != 0 {
			node.HasActiveRentContract = true
		}
	}

	node.PublicConfig = n.networkHasPublicConfig(ctx, node.TwinID)
	if !node.PublicConfig {
		return fmt.Errorf("failed to get public config of node %d", node.ID)
	}

	wgPorts, err := n.networkListWGPorts(ctx, node.TwinID)
	if err != nil {
		return fmt.Errorf("failed to update the wireguard ports used by node %d with error: %w", node.ID, err)
	}
	node.WgPorts = wgPorts

	n.logger.Debug().Msgf("capacity updated for node %d:\n%v\nhas active rent contract: %v", node.ID, node.Resources, node.HasActiveRentContract)
	return nil
}

// GetStoragePools executes zos system version cmd
func (n *rmbNodeClient) getStoragePools(ctx context.Context, nodeTwin uint32) (pools []pkg.PoolMetrics, err error) {
	const cmd = "zos.storage.pools"
	err = n.rmb.Call(ctx, nodeTwin, cmd, nil, &pools)
	return pools, err
}

// SystemVersion executes zos system version cmd
func (n *rmbNodeClient) systemVersion(ctx context.Context, nodeTwin uint32) error {
	const cmd = "zos.system.version"
	return n.rmb.Call(ctx, nodeTwin, cmd, nil, nil)
}

// NetworkHasPublicConfig returns the current public node network configuration. A node with a
// public config can be used as an access node for wireguard.
func (n *rmbNodeClient) networkHasPublicConfig(ctx context.Context, nodeTwin uint32) bool {
	const cmd = "zos.network.public_config_get"

	if err := n.rmb.Call(ctx, nodeTwin, cmd, nil, nil); err != nil {
		return false
	}

	return true
}

// statistics returns some node statistics. Including total and available cpu, memory, storage, etc...
func (n *rmbNodeClient) statistics(ctx context.Context, nodeTwin uint32) (result models.ConsumableResources, err error) {
	const cmd = "zos.statistics.get"
	var res struct {
		Total gridtypes.Capacity `json:"total"`
		Used  gridtypes.Capacity `json:"used"`
	}
	if err = n.rmb.Call(ctx, nodeTwin, cmd, nil, &res); err != nil {
		return
	}
	result.Total = models.Capacity{HRU: uint64(res.Total.HRU), SRU: uint64(res.Total.SRU), CRU: res.Total.CRU, MRU: uint64(res.Total.MRU), Ipv4: res.Total.IPV4U}
	result.Used = models.Capacity{HRU: uint64(res.Used.HRU), SRU: uint64(res.Used.SRU), CRU: res.Used.CRU, MRU: uint64(res.Used.MRU), Ipv4: res.Used.IPV4U}
	return result, nil
}

// networkListWGPorts return a list of all "taken" ports on the node. A new deployment
// should be careful to use a free port for its network setup.
func (n *rmbNodeClient) networkListWGPorts(ctx context.Context, nodeTwin uint32) ([]uint16, error) {
	const cmd = "zos.network.list_wg_ports"
	var result []uint16

	if err := n.rmb.Call(ctx, nodeTwin, cmd, nil, &result); err != nil {
		return nil, err
	}

	return result, nil
}

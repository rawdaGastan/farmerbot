// Package internal for farmerbot internals
package internal

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/rmb-sdk-go"
	"github.com/threefoldtech/rmb-sdk-go/direct"
	"github.com/threefoldtech/substrate-client"
)

// RMBClient is an rmb abstract client interface.
/*type RMBClient interface {
	Call(ctx context.Context, twin uint32, fn string, data interface{}, result interface{}) error
}*/

type rmbNodeClient struct {
	logger zerolog.Logger
	rmb    rmb.Client //RMBClient
}

func newRmbNodeClient(sub substrate.Substrate, identity substrate.Identity, network string, logger zerolog.Logger) (rmbNodeClient, error) {
	sessionID := fmt.Sprintf("tf-%d", os.Getpid())
	db := newTwinDB(sub)
	rmbClient, err := direct.NewClient(context.Background(), identity, constants.RelayURLS[network], sessionID, db)
	if err != nil {
		return rmbNodeClient{}, fmt.Errorf("failed with error: %w, couldn't create rmb client", err)
	}

	return rmbNodeClient{
		rmb: rmbClient,
	}, nil
}

type twinDB struct {
	cache *cache.Cache
	sub   substrate.Substrate
}

// newTwinDB creates a new twinDBImpl instance, with a non expiring cache.
func newTwinDB(sub substrate.Substrate) direct.TwinDB {
	return &twinDB{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
		sub:   sub,
	}
}

// Get gets Twin from cache if present. if not, gets it from substrate client and caches it.
func (t *twinDB) Get(id uint32) (direct.Twin, error) {
	cachedValue, ok := t.cache.Get(fmt.Sprint(id))
	if ok {
		return cachedValue.(direct.Twin), nil
	}

	twin, err := t.sub.GetTwin(id)
	if err != nil {
		return direct.Twin{}, errors.Wrapf(err, "could not get twin of twin with id %d", id)
	}

	directTwin := direct.Twin{
		ID:        id,
		PublicKey: twin.Account.PublicKey(),
	}

	err = t.cache.Add(fmt.Sprint(id), directTwin, cache.DefaultExpiration)
	if err != nil {
		return direct.Twin{}, errors.Wrapf(err, "could not set cache for twin with id %d", id)
	}

	return directTwin, nil
}

// GetByPk returns a twin's id using its public key
func (t *twinDB) GetByPk(pk []byte) (uint32, error) {
	return t.sub.GetTwinByPubKey(pk)
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
	if node.TimeoutClaimedResources == 0 {
		stats, err := n.statistics(ctx, node.TwinID)
		if err != nil {
			return fmt.Errorf("failed to get statistics of node %d with error: %w", node.ID, err)
		}
		node.UpdateResources(stats)
	} else {
		node.TimeoutClaimedResources--
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

	n.logger.Debug().Msgf("capacity updated for node %d:\n%v", node.ID, node.Resources)
	return nil
}

// SystemVersion executes zos system version cmd
func (n *rmbNodeClient) systemVersion(ctx context.Context, nodeTwin uint32) error {
	const cmd = "zos.system.version"

	if err := n.rmb.Call(ctx, nodeTwin, cmd, nil, nil); err != nil {
		return err
	}

	return nil
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

	if err = n.rmb.Call(ctx, nodeTwin, cmd, nil, &result); err != nil {
		return
	}

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

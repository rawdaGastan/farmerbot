// Package manager provides how to manage nodes, farms and power
package manager

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/substrate-client"
)

// NodeManager manages nodes
type NodeManager struct {
	logger   zerolog.Logger
	db       models.RedisManager
	identity substrate.Identity
	subConn  models.Sub
}

// NewNodeManager creates a new NodeManager
func NewNodeManager(mnemonics string, subConn models.Sub, db models.RedisManager, logger zerolog.Logger) (NodeManager, error) {
	identity, err := substrate.NewIdentityFromSr25519Phrase(mnemonics)
	if err != nil {
		return NodeManager{}, err
	}

	return NodeManager{logger, db, identity, subConn}, nil
}

// Define defines a node
func (n *NodeManager) Define(node models.Node) error {
	n.logger.Debug().Msgf("node is %+v", node)
	return n.db.UpdatesNodes(node)
}

// FindNode finds an available node in the farm
func (n *NodeManager) FindNode(nodeOptions models.NodeOptions, nodesToExclude []uint) (uint32, error) {
	nodes, err := n.db.GetNodes()
	if err != nil {
		return 0, errors.New("failed to get nodes from db")
	}

	farm, err := n.db.GetFarm()
	if err != nil {
		return 0, errors.New("failed to get farm from db")
	}

	if nodeOptions.PublicIPs > 0 {
		var publicIPsUsedByNodes uint64

		for _, node := range nodes {
			publicIPsUsedByNodes += node.PublicIPsUsed
		}
		if publicIPsUsedByNodes+nodeOptions.PublicIPs > farm.PublicIPs {
			return 0, fmt.Errorf("no more public ips available for farm %d", farm.ID)
		}
	}

	var possibleNodes []models.Node
	for _, node := range nodes {
		if nodeOptions.Certified && !node.Certified {
			continue
		}

		if nodeOptions.PublicConfig && !node.PublicConfig {
			continue
		}

		if node.HasActiveRentContract {
			continue
		}

		if nodeOptions.Dedicated && (!node.Dedicated || !node.IsUnused()) {
			continue
		}

		// TODO: what if the node resources are used
		if !nodeOptions.Dedicated && node.Dedicated && nodeOptions.Capacity != node.Resources.Total {
			continue
		}

		if contains(nodesToExclude, uint(node.ID)) {
			continue
		}
		if !node.CanClaimResources(nodeOptions.Capacity) {
			continue
		}
		possibleNodes = append(possibleNodes, node)
	}

	if len(possibleNodes) == 0 {
		return 0, fmt.Errorf("could not find a suitable node with the given options: %v", possibleNodes)
	}

	// Sort the nodes on power state (the ones that are ON first)
	sort.Slice(possibleNodes, func(i, j int) bool {
		return possibleNodes[i].PowerState.ON
	})

	nodeFounded := possibleNodes[0]
	n.logger.Debug().Msgf("Found a node: %d", nodeFounded.ID)

	// claim the resources until next update of the data
	// add a timeout (after 30 minutes we update the resources)
	nodeFounded.TimeoutClaimedResources = time.Now().Add(constants.TimeoutPowerStateChange)
	if nodeOptions.Dedicated {
		// claim all capacity
		nodeFounded.ClaimResources(nodeFounded.Resources.Total)
	} else {
		nodeFounded.ClaimResources(nodeOptions.Capacity)
	}

	// claim public ips until next update of the data
	if nodeOptions.PublicIPs > 0 {
		nodeFounded.PublicIPsUsed += nodeOptions.PublicIPs
	}

	if err := n.powerOn(nodeFounded); err != nil {
		return 0, err
	}

	return nodeFounded.ID, nil
}

// PowerOn power on a node
func (n *NodeManager) powerOn(node models.Node) error {
	n.logger.Info().Msgf("POWER ON: %d", node.ID)
	return node.SetNodePower(n.identity, n.subConn, true)
}

// Contains check if a slice contains an element
func contains[T comparable](elements []T, element T) bool {
	for _, e := range elements {
		if element == e {
			return true
		}
	}
	return false
}

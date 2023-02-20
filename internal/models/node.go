// Package models for farmerbot models.
package models

import (
	"fmt"
	"math"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/threefoldtech/substrate-client"
)

// Node represents a node in a farm
type Node struct {
	ID                        uint32              `json:"id"`
	TwinID                    uint32              `json:"twinID"`
	FarmID                    uint32              `json:"farmID,omitempty"`
	Description               string              `json:"description,omitempty"`
	Certified                 bool                `json:"certified,omitempty"`
	Dedicated                 bool                `json:"dedicated,omitempty"`
	PublicConfig              bool                `json:"publicConfig,omitempty"`
	PublicIPsUsed             uint64              `json:"publicIPsUsed,omitempty"`
	WgPorts                   []uint16            `json:"wgPorts,omitempty"`
	Resources                 ConsumableResources `json:"resources"`
	PowerState                PowerState          `json:"powerState,omitempty"`
	TimeoutClaimedResources   uint8               `json:"timeoutClaimedResources,omitempty"`
	LastTimePowerStateChanged time.Time           `json:"lastTimePowerStateChanged,omitempty"`
	LastTimeAwake             time.Time           `json:"lastTimeAwake,omitempty"`
}

// PowerState is the state of node's power
type PowerState struct {
	ON           bool `json:"on,omitempty"`
	WakingUp     bool `json:"wakingUp,omitempty"`
	OFF          bool `json:"off,omitempty"`
	ShuttingDown bool `json:"shuttingDown,omitempty"`
}

// NodeOptions represents the options to find a node
type NodeOptions struct {
	Certified    bool     `json:"certified,omitempty"`
	Dedicated    bool     `json:"dedicated,omitempty"`
	PublicConfig bool     `json:"publicConfig,omitempty"`
	PublicIPs    uint64   `json:"publicIPs,omitempty"`
	Capacity     Capacity `json:"capacity,omitempty"`
}

// Sub is substrate client interface
type Sub interface {
	SetNodePowerState(identity substrate.Identity, up bool) (hash types.Hash, err error)
}

// SetNodePower sets the node power
func (n *Node) SetNodePower(identity substrate.Identity, subConn Sub, on bool) error {
	if on && (n.PowerState.ON || n.PowerState.WakingUp) {
		return nil
	}

	if !on && (n.PowerState.OFF || n.PowerState.ShuttingDown) {
		return nil
	}

	// make sure the node isn't waking up or shutting down
	if err := n.ensureNodeIsOnOrOff(); err != nil {
		return err
	}

	_, err := subConn.SetNodePowerState(identity, on)
	if err != nil {
		return err
	}

	// update nodes
	n.PowerState.ShuttingDown = !on
	n.PowerState.WakingUp = on
	n.LastTimePowerStateChanged = time.Now()

	return nil
}

// UpdateResources updates the node resources
func (n *Node) UpdateResources(cap ConsumableResources) {
	n.Resources.Total = cap.Total
	n.Resources.Used = cap.Used
	n.PublicIPsUsed = cap.Used.Ipv4
}

// IsUnused node is an empty node
func (n *Node) IsUnused() bool {
	return n.Resources.Used.isEmpty()
}

// CanClaimResources checks if a node can claim some resources
func (n *Node) CanClaimResources(cap Capacity) bool {
	total := n.Resources.Total
	total.CRU = uint64(math.Ceil(float64(total.CRU) * n.Resources.OverProvisionCPU))

	free := total.subtract(n.Resources.Used)
	return total.CRU >= cap.CRU && free.CRU >= cap.CRU && free.MRU >= cap.MRU && free.HRU >= cap.HRU && free.SRU >= cap.SRU
}

// ClaimResources claims the resources from a node
func (n *Node) ClaimResources(cap Capacity) {
	n.Resources.Used.Add(cap)
}

// FilterOffNodes filters off nodes
func FilterOffNodes(nodes []Node) []Node {
	out := make([]Node, 0)
	for _, node := range nodes {
		if node.PowerState.OFF {
			out = append(out, node)
		}
	}
	return out
}

// FilterUnusedOnNodes filters nodes that are ON and unused
func FilterUnusedOnNodes(nodes []Node) []Node {
	out := make([]Node, 0)
	for _, node := range nodes {
		if node.PowerState.ON && node.IsUnused() {
			out = append(out, node)
		}
	}
	return out
}

// FilterWakingOrShuttingNodes filters nodes that are waking up or shutting down
func FilterWakingOrShuttingNodes(nodes []Node) []Node {
	out := make([]Node, 0)
	for _, node := range nodes {
		if node.PowerState.WakingUp || node.PowerState.ShuttingDown {
			out = append(out, node)
		}
	}
	return out
}

// EnsureNodeIsOnOrOff make sure node is ON or OFF not waking up or shutting down
func (n *Node) ensureNodeIsOnOrOff() error {
	if n.PowerState.WakingUp {
		return fmt.Errorf("node %d is waking up", n.ID)
	}
	if n.PowerState.ShuttingDown {
		return fmt.Errorf("node %d is shutting down", n.ID)
	}
	return nil
}

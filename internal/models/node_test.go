// Package models for farmerbot models.
package models

import (
	"fmt"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNodeModel(t *testing.T) {
	node := Node{
		ID:     1,
		TwinID: 1,
		Resources: ConsumableResources{
			OverProvisionCPU: 1,
			Total:            cap,
		},
		PowerState: PowerState{
			ON: true,
		},
	}

	t.Run("test ensure node is on/off", func(t *testing.T) {
		err := node.ensureNodeIsOnOrOff()
		assert.NoError(t, err)

		node.PowerState.WakingUp = true
		err = node.ensureNodeIsOnOrOff()
		assert.Error(t, err)

		node.PowerState.WakingUp = false
		node.PowerState.ShuttingDown = true
		err = node.ensureNodeIsOnOrOff()
		assert.Error(t, err)

		node.PowerState.ShuttingDown = false
	})

	// set power from node models tests
	t.Run("test set node power", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		sub := NewMockSub(ctrl)

		// set power on for already on node
		err := node.SetNodePower(nil, sub, true)
		assert.NoError(t, err)

		node.PowerState.OFF = true
		node.PowerState.ON = false

		// set power off for already node off
		err = node.SetNodePower(nil, sub, false)
		assert.NoError(t, err)

		// set power on for node off
		sub.EXPECT().SetNodePowerState(nil, true)
		err = node.SetNodePower(nil, sub, true)
		assert.NoError(t, err)

		// set power off for node waking up -> error
		err = node.SetNodePower(nil, sub, false)
		assert.Error(t, err)

		node.PowerState.WakingUp = false
		node.PowerState.ON = true
		node.PowerState.OFF = false

		// set power on for already on node
		err = node.SetNodePower(nil, sub, true)
		assert.NoError(t, err)

		// set power off for node on but substrate failed -> error
		sub.EXPECT().SetNodePowerState(nil, false).Return(types.Hash{}, fmt.Errorf("error"))
		err = node.SetNodePower(nil, sub, false)
		assert.Error(t, err)
	})

	t.Run("test update node resources", func(t *testing.T) {
		node.UpdateResources(node.Resources)
		assert.True(t, node.Resources.Used.isEmpty())
		assert.True(t, node.IsUnused())
		assert.Equal(t, node.Resources.OverProvisionCPU, float64(1))
		assert.True(t, node.CanClaimResources(node.Resources.Total))

		node.ClaimResources(node.Resources.Total)
		assert.False(t, node.Resources.Used.isEmpty())
		assert.False(t, node.IsUnused())
		assert.False(t, node.CanClaimResources(node.Resources.Total))

		node.Resources.Used = Capacity{}
	})

	t.Run("test node filters", func(t *testing.T) {
		nodes := FilterOffNodes([]Node{node})
		assert.Empty(t, nodes)

		node.PowerState.OFF = true
		node.PowerState.ON = false

		nodes = FilterOffNodes([]Node{node})
		assert.NotEmpty(t, nodes)

		nodes = FilterUnusedOnNodes([]Node{node})
		assert.Empty(t, nodes)

		node.PowerState.OFF = false
		node.PowerState.ON = true

		nodes = FilterUnusedOnNodes([]Node{node})
		assert.NotEmpty(t, nodes)

		nodes = FilterWakingOrShuttingNodes([]Node{node})
		assert.Empty(t, nodes)

		node.PowerState.ShuttingDown = true

		nodes = FilterWakingOrShuttingNodes([]Node{node})
		assert.NotEmpty(t, nodes)
	})
}

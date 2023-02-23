// Package manager provides how to manage nodes, nodes and power
package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/golang/mock/gomock"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var nodeCapacity = models.Capacity{
	CRU: 1,
	SRU: 1,
	MRU: 1,
	HRU: 1,
}

var node = models.Node{
	ID:     1,
	TwinID: 1,
	Resources: models.ConsumableResources{
		OverProvisionCPU: 1,
		Total:            nodeCapacity,
	},
	PowerState: models.PowerState{
		ON: true,
	},
}

func TestNodeManager(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mocks.NewMockRedisManager(ctrl)
	sub := models.NewMockSub(ctrl)

	_, err := NewNodeManager("bad", sub, db, log.Logger)
	assert.Error(t, err)

	nodeManager, err := NewNodeManager("", sub, db, log.Logger)
	assert.NoError(t, err)
	nodeManager.subConn = sub

	nodeOptions := models.NodeOptions{
		PublicIPs: 1,
		Capacity:  nodeCapacity,
	}

	t.Run("test valid define node", func(t *testing.T) {
		db.EXPECT().UpdatesNodes(node).Return(nil)

		nodeBytes, err := json.Marshal(node)
		assert.NoError(t, err)

		err = nodeManager.Define(nodeBytes)
		assert.NoError(t, err)
	})

	t.Run("test invalid define node: db failed", func(t *testing.T) {
		db.EXPECT().UpdatesNodes(node).Return(fmt.Errorf("error"))

		nodeBytes, err := json.Marshal(node)
		assert.NoError(t, err)

		err = nodeManager.Define(nodeBytes)
		assert.Error(t, err)
	})

	t.Run("test invalid define node: wrong input", func(t *testing.T) {
		nodeBytes, err := json.Marshal("node")
		assert.NoError(t, err)

		err = nodeManager.Define(nodeBytes)
		assert.Error(t, err)
	})

	t.Run("test valid find node: found an ON node", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		node, err = nodeManager.FindNode(nodeOptions, []uint{})
		assert.NoError(t, err)
	})

	t.Run("test valid find node: found an OFF node", func(t *testing.T) {
		node.PowerState.OFF = true
		node.PowerState.ON = false

		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		sub.EXPECT().SetNodePowerState(nodeManager.identity, true)

		node, err = nodeManager.FindNode(models.NodeOptions{}, []uint{})
		assert.NoError(t, err)
	})

	t.Run("test invalid find node: found an OFF node but change power failed", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		sub.EXPECT().SetNodePowerState(nodeManager.identity, true).Return(types.Hash{}, fmt.Errorf("error"))

		_, err = nodeManager.FindNode(models.NodeOptions{}, []uint{})
		assert.Error(t, err)

		node.PowerState.ON = true
		node.PowerState.OFF = false
	})

	t.Run("test invalid find node: no more public ips", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(nodeOptions, []uint{})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: certified so no nodes found", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{Certified: true}, []uint{})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: publicConfig so no nodes found", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{PublicConfig: true}, []uint{})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: dedicated so no nodes found", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{Dedicated: true}, []uint{})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: node is dedicated so no nodes found", func(t *testing.T) {
		node.Dedicated = true
		node.Resources.Total = models.Capacity{}

		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{Capacity: nodeCapacity}, []uint{})
		assert.Error(t, err)
		node.Dedicated = false
		node.Resources.Total = nodeCapacity
	})

	t.Run("test invalid find node: node is excluded", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{}, []uint{uint(node.ID)})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: node cannot claim resources", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{Capacity: nodeCapacity}, []uint{})
		assert.Error(t, err)
	})

	t.Run("test valid find node: both are dedicated and node is unused", func(t *testing.T) {
		node.Resources.Used = models.Capacity{}
		node.Dedicated = true
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, nil)

		_, err = nodeManager.FindNode(models.NodeOptions{Dedicated: true}, []uint{})
		assert.NoError(t, err)
	})

	t.Run("test invalid find node: failed DB to get nodes", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, fmt.Errorf("error"))

		_, err = nodeManager.FindNode(models.NodeOptions{}, []uint{})
		assert.Error(t, err)
	})

	t.Run("test invalid find node: failed DB to get farm", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetFarm().Return(testFarm, fmt.Errorf("error"))

		_, err = nodeManager.FindNode(models.NodeOptions{}, []uint{})
		assert.Error(t, err)
	})
}

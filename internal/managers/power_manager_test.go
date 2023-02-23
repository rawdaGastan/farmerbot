// Package manager provides how to manage powers, powers and power
package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	types "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/golang/mock/gomock"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPowerManager(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := mocks.NewMockRedisManager(ctrl)
	sub := models.NewMockSub(ctrl)

	_, err := NewPowerManager("bad", sub, db, log.Logger)
	assert.Error(t, err)

	powerManager, err := NewPowerManager("", sub, db, log.Logger)
	assert.NoError(t, err)
	powerManager.subConn = sub

	power := models.Power{
		WakeUpThreshold: 80,
		PeriodicWakeup:  models.WakeupDate(time.Now()),
	}
	power.PeriodicWakeup = models.WakeupDate(power.PeriodicWakeup.PeriodicWakeupStart())

	t.Run("test valid configure power", func(t *testing.T) {
		db.EXPECT().SetPower(power).Return(nil)

		powerBytes, err := json.Marshal(power)
		assert.NoError(t, err)

		err = powerManager.Configure(powerBytes)
		assert.NoError(t, err)
	})

	t.Run("test invalid configure power: db failed", func(t *testing.T) {
		db.EXPECT().SetPower(power).Return(fmt.Errorf("error"))

		powerBytes, err := json.Marshal(power)
		assert.NoError(t, err)

		err = powerManager.Configure(powerBytes)
		assert.Error(t, err)
	})

	t.Run("test invalid configure power: wrong input", func(t *testing.T) {
		powerBytes, err := json.Marshal("power")
		assert.NoError(t, err)

		err = powerManager.Configure(powerBytes)
		assert.Error(t, err)
	})

	t.Run("test valid power on", func(t *testing.T) {
		node.PowerState.OFF = true
		node.PowerState.ON = false
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, true).Return(types.Hash{}, nil)
		//mockAny because I can't match node state change time
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(nil)

		err = powerManager.PowerOn(node.ID)
		assert.NoError(t, err)
	})

	t.Run("test invalid power on: node not found", func(t *testing.T) {
		db.EXPECT().GetNode(node.ID).Return(node, fmt.Errorf("error"))

		err = powerManager.PowerOn(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power on: set node failed", func(t *testing.T) {
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, true).Return(types.Hash{}, fmt.Errorf("error"))

		err = powerManager.PowerOn(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power on: update nodes failed", func(t *testing.T) {
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, true).Return(types.Hash{}, nil)
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(fmt.Errorf("error"))

		err = powerManager.PowerOn(node.ID)
		assert.Error(t, err)
	})

	t.Run("test valid power off", func(t *testing.T) {
		node.PowerState.OFF = false
		node.PowerState.ON = true
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, false).Return(types.Hash{}, nil)
		//mockAny because I can't match node state change time
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(nil)

		err = powerManager.PowerOff(node.ID)
		assert.NoError(t, err)
	})

	t.Run("test invalid power off: one node is on and cannot be off", func(t *testing.T) {
		db.EXPECT().FilterOnNodes().Return([]models.Node{node}, nil)

		err = powerManager.PowerOff(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power off: filter on nodes error", func(t *testing.T) {
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, fmt.Errorf("error"))

		err = powerManager.PowerOff(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power off: node not found", func(t *testing.T) {
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetNode(node.ID).Return(node, fmt.Errorf("error"))

		err = powerManager.PowerOff(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power off: set node failed", func(t *testing.T) {
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, false).Return(types.Hash{}, fmt.Errorf("error"))

		err = powerManager.PowerOff(node.ID)
		assert.Error(t, err)
	})

	t.Run("test invalid power off: update nodes failed", func(t *testing.T) {
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, false).Return(types.Hash{}, nil)
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(fmt.Errorf("error"))

		err = powerManager.PowerOff(node.ID)
		assert.Error(t, err)
	})

	t.Run("test valid periodic wakeup: already on", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		err = powerManager.PeriodicWakeup()
		assert.NoError(t, err)
	})

	t.Run("test valid periodic wakeup", func(t *testing.T) {
		node.PowerState.OFF = true
		node.PowerState.ON = false
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set node power state on mocks
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, true).Return(types.Hash{}, nil)
		//mockAny because I can't match node state change time
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(nil)

		err = powerManager.PeriodicWakeup()
		assert.NoError(t, err)
	})

	t.Run("test invalid periodic wakeup: failed to get nodes from db", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, fmt.Errorf("error"))

		err = powerManager.PeriodicWakeup()
		assert.Error(t, err)
	})

	t.Run("test invalid periodic wakeup: failed to get power from db", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, fmt.Errorf("error"))

		err = powerManager.PeriodicWakeup()
		assert.Error(t, err)
	})

	t.Run("test invalid periodic wakeup: failed to set power state on", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set node power state on mocks
		db.EXPECT().GetNode(node.ID).Return(node, fmt.Errorf("error"))

		err = powerManager.PeriodicWakeup()
		assert.Error(t, err)
	})

	t.Run("test valid power management: a node to shutdown", func(t *testing.T) {
		node.PowerState.OFF = false
		node.PowerState.ON = true

		db.EXPECT().GetNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set power off to the second node
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, false).Return(types.Hash{}, nil)
		//mockAny because I can't match node state change time
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)
	})

	t.Run("test valid power management: nothing to shut down", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)
	})

	t.Run("test valid power management: cannot shutdown public config", func(t *testing.T) {
		node.PublicConfig = true
		db.EXPECT().GetNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)
		node.PublicConfig = false
	})

	t.Run("test valid power management: node is waking up", func(t *testing.T) {
		node.PowerState.WakingUp = true
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)
		node.PowerState.WakingUp = false
	})

	t.Run("test valid power management: no total resources", func(t *testing.T) {
		node.Resources.Total = models.Capacity{}
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)
		node.Resources.Total = nodeCapacity
	})

	t.Run("test valid/invalid power management: a node to wake up", func(t *testing.T) {
		// add an on node used
		node.Resources.Used = nodeCapacity
		nodes := []models.Node{node}

		node.PowerState.OFF = true
		node.PowerState.ON = false
		nodes = append(nodes, node)

		db.EXPECT().GetNodes().Return(nodes, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set power on to the node
		db.EXPECT().GetNode(node.ID).Return(node, nil)
		sub.EXPECT().SetNodePowerState(powerManager.identity, true).Return(types.Hash{}, nil)
		//mockAny because I can't match node state change time
		db.EXPECT().UpdatesNodes(gomock.Any()).Return(nil)

		err = powerManager.PowerManagement()
		assert.NoError(t, err)

		// invalid
		db.EXPECT().GetNodes().Return(nodes, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set power on to the node
		db.EXPECT().GetNode(node.ID).Return(node, fmt.Errorf("error"))

		err = powerManager.PowerManagement()
		assert.Error(t, err)
		node.Resources.Used = models.Capacity{}
	})

	t.Run("test invalid power management: failed to shutdown node", func(t *testing.T) {
		node.PowerState.OFF = false
		node.PowerState.ON = true

		db.EXPECT().GetNodes().Return([]models.Node{node, node}, nil)
		db.EXPECT().GetPower().Return(power, nil)

		// set power off to the second node
		db.EXPECT().FilterOnNodes().Return([]models.Node{node, node}, fmt.Errorf("error"))

		err = powerManager.PowerManagement()
		assert.Error(t, err)
	})

	t.Run("test invalid power management: failed to get nodes from db", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, fmt.Errorf("error"))

		err = powerManager.PowerManagement()
		assert.Error(t, err)
	})

	t.Run("test invalid power management: failed to get power from db", func(t *testing.T) {
		db.EXPECT().GetNodes().Return([]models.Node{node}, nil)
		db.EXPECT().GetPower().Return(power, fmt.Errorf("error"))

		err = powerManager.PowerManagement()
		assert.Error(t, err)
	})

}

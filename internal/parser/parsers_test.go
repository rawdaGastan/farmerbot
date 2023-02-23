// Package parser for parsing cmd configs
package parser

import (
	"fmt"
	"testing"
	"time"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestParsers(t *testing.T) {
	t.Run("test invalid file", func(t *testing.T) {
		_, err := ReadFile("json.json")
		assert.Error(t, err)
	})

	t.Run("test valid file", func(t *testing.T) {
		_, err := ReadFile("parser.go")
		assert.NoError(t, err)
	})

	t.Run("test invalid json", func(t *testing.T) {
		content := `
		{ 
			"nodes": [ ],
			"farm": , 
			"power": ,
		}
		`

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoFarm([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoPower([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNodeOptions([]byte(content))
		assert.Error(t, err)
	})

	t.Run("test valid node options", func(t *testing.T) {
		content := `
		{ 
			"certified": true,
			"dedicated": true,
			"publicConfig": false,
			"publicIPs": 1, 
			"capacity": { "SRU": 1, "CRU": 2, "HRU": 3, "MRU": 4 }
		}
		`

		options, err := ParseJSONIntoNodeOptions([]byte(content))
		assert.NoError(t, err)
		assert.True(t, options.Certified)
		assert.True(t, options.Dedicated)
		assert.False(t, options.PublicConfig)
		assert.Equal(t, options.PublicIPs, uint64(1))
		assert.Equal(t, options.Capacity.CRU, uint64(2))
		assert.Equal(t, options.Capacity.MRU, uint64(4))
		assert.Equal(t, options.Capacity.SRU, uint64(1))
		assert.Equal(t, options.Capacity.HRU, uint64(3))
	})

	t.Run("test valid json", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "SRU": 1, "CRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM", "WakeUpThreshold": 30 }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		c, err := ParseJSONIntoConfig([]byte(content))
		assert.NoError(t, err)
		assert.Equal(t, c.Power.WakeUpThreshold, constants.MinWakeUpThreshold)
		assert.Equal(t, c.Nodes[0].Resources.OverProvisionCPU, float64(1))

		f, err := ParseJSONIntoFarm([]byte(farmContent))
		assert.NoError(t, err)
		assert.Equal(t, f.ID, uint32(1))

		n, err := ParseJSONIntoNode([]byte(nodeContent))
		assert.NoError(t, err)
		assert.Equal(t, n.ID, uint32(1))
		assert.Equal(t, n.Resources.OverProvisionCPU, float64(1))

		p, err := ParseJSONIntoPower([]byte(powerContent))
		assert.NoError(t, err)
		now := time.Now()
		assert.Equal(t, p.PeriodicWakeup.PeriodicWakeupStart(), time.Time(time.Date(now.Year(), now.Month(), now.Day(), 8, 30, 0, 0, time.Local)))
	})

	t.Run("test invalid json no node ID", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "twinID" : 1, "resources": { "total": { "SRU": 1, "CRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM", "WakeUpThreshold": 90 }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		c, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)
		assert.Equal(t, c.Power.WakeUpThreshold, constants.MaxWakeUpThreshold)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)

		p, err := ParseJSONIntoPower([]byte(powerContent))
		assert.NoError(t, err)
		assert.Equal(t, p.WakeUpThreshold, constants.MaxWakeUpThreshold)
	})

	t.Run("test invalid json no node twin ID", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "resources": { "total": { "SRU": 1, "CRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		c, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)
		assert.Equal(t, c.Power.WakeUpThreshold, constants.DefaultWakeUpThreshold)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)

		p, err := ParseJSONIntoPower([]byte(powerContent))
		assert.NoError(t, err)
		assert.Equal(t, p.WakeUpThreshold, constants.DefaultWakeUpThreshold)
	})

	t.Run("test invalid json no node sru", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "CRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)
	})

	t.Run("test invalid json no cru", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "SRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)
	})

	t.Run("test invalid json no hru", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "SRU": 1, "CRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)
	})

	t.Run("test invalid json no mru", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "SRU": 1, "CRU": 1, "HRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)
	})

	t.Run("test invalid json node over provision CPU", func(t *testing.T) {
		farmContent := `{ "ID": 1 }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "overProvisionCPU": 5, "total": { "SRU": 1, "CRU": 1, "HRU": 1, "CRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoNode([]byte(nodeContent))
		assert.Error(t, err)
	})

	t.Run("test invalid json no farm ID", func(t *testing.T) {
		farmContent := `{ }`
		nodeContent := `{ "ID": 1, "twinID" : 1, "resources": { "total": { "SRU": 1, "CRU": 1, "HRU": 1, "MRU": 1 } } }`
		powerContent := `{ "periodicWakeup": "08:30AM" }`
		content := fmt.Sprintf(`
		{ 
			"nodes": [ %v ],
			"farm": %v, 
			"power": %v
		}
		`, nodeContent, farmContent, powerContent)

		_, err := ParseJSONIntoConfig([]byte(content))
		assert.Error(t, err)

		_, err = ParseJSONIntoFarm([]byte(farmContent))
		assert.Error(t, err)
	})
}

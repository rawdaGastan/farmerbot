// Package models for farmerbot models.
package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPowerModel(t *testing.T) {
	power := Power{
		WakeUpThreshold: 80,
		PeriodicWakeup:  WakeupDate(time.Now()),
	}
	oldPower := time.Time(power.PeriodicWakeup)

	// invalid
	err := power.PeriodicWakeup.UnmarshalJSON([]byte("7:3"))
	assert.Error(t, err)

	// valid
	wakeUpBytes, err := power.PeriodicWakeup.MarshalJSON()
	assert.NoError(t, err)

	err = power.PeriodicWakeup.UnmarshalJSON(wakeUpBytes)
	assert.NoError(t, err)

	assert.Equal(t, time.Time(power.PeriodicWakeup).Hour(), oldPower.Hour())
	assert.Equal(t, time.Time(power.PeriodicWakeup).Minute(), oldPower.Minute())
	assert.NotEqual(t, time.Time(power.PeriodicWakeup).Day(), oldPower.Day())

	power.PeriodicWakeup = WakeupDate(power.PeriodicWakeup.PeriodicWakeupStart())
	assert.Equal(t, time.Time(power.PeriodicWakeup).Day(), oldPower.Day())
}

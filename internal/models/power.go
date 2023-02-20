// Package models for farmerbot models.
package models

import (
	"encoding/json"
	"strings"
	"time"
)

// WakeupDate is the date to wakeup all nodes
type WakeupDate time.Time

// Power represents power configuration
type Power struct {
	WakeUpThreshold uint64     `json:"wakeUpThreshold"`
	PeriodicWakeup  WakeupDate `json:"periodicWakeUp"`
}

// UnmarshalJSON unmarshals the given JSON object into wakeUp date
func (d *WakeupDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("03:04PM", s)
	if err != nil {
		return err
	}
	*d = WakeupDate(t)
	return nil
}

// MarshalJSON marshals the wakeup date
func (d WakeupDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

// PeriodicWakeupStart returns periodic wakeup start date
func (d WakeupDate) PeriodicWakeupStart() time.Time {
	date := time.Time(d)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return today.Local().Add(time.Hour*time.Duration(date.Hour()) +
		time.Minute*time.Duration(date.Minute()) +
		time.Second*time.Duration(date.Second()))
}

// Package manager provides how to manage nodes, farms and power
package manager

import (
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
)

// FarmManager manages farms
type FarmManager struct {
	logger zerolog.Logger
	db     models.RedisManager
}

// NewFarmManager creates a new FarmManager
func NewFarmManager(db models.RedisManager, logger zerolog.Logger) FarmManager {
	return FarmManager{logger, db}
}

// Define defines a farm
func (f *FarmManager) Define(farm models.Farm) error {
	f.logger.Debug().Msgf("farm is %+v", farm)
	return f.db.SetFarm(farm)
}

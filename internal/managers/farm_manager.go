// Package manager provides how to manage nodes, farms and power
package manager

import (
	"fmt"

	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/rs/zerolog"
)

// FarmHandler interface for mocks
type FarmHandler interface {
	Define(farm models.Farm) error
}

// FarmManager manages farms
type FarmManager struct {
	logger zerolog.Logger
	db     models.RedisDB
}

// NewFarmManager creates a new FarmManager
func NewFarmManager(address string, logger zerolog.Logger) FarmManager {
	return FarmManager{logger, models.NewRedisDB(address)}
}

// Define defines a farm
func (f *FarmManager) Define(jsonContent []byte) error {
	farm, err := parser.ParseJSONIntoFarm(jsonContent)
	if err != nil {
		return fmt.Errorf("failed to get farm from json content: %v", err)
	}

	f.logger.Debug().Msgf("farm is %+v", farm)
	return f.db.SetFarm(farm)
}

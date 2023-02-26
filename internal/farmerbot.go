// Package internal for farmerbot internals
package internal

import (
	"context"
	"time"

	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/substrate-client"
)

// TODO: add a msg bus for all commands

// FarmerBot for managing farms
type FarmerBot struct {
	logger        zerolog.Logger
	db            models.RedisDB
	rmbNodeClient rmbNodeClient
	powerManager  manager.PowerManager
}

// NewFarmerBot generates a new farmer bot
func NewFarmerBot(configPath string, network string, mnemonics string, sub *substrate.Substrate, db models.RedisDB, logger zerolog.Logger) (FarmerBot, error) {
	farmerBot := FarmerBot{}
	jsonContent, err := parser.ReadFile(configPath)
	if err != nil {
		return farmerBot, err
	}

	config, err := parser.ParseJSONIntoConfig(jsonContent)
	if err != nil {
		return farmerBot, err
	}

	rmbNodeClient, err := newRmbNodeClient(sub, mnemonics, network, logger)
	if err != nil {
		return farmerBot, err
	}

	err = db.SaveConfig(config)
	if err != nil {
		return farmerBot, err
	}

	powerManager, err := manager.NewPowerManager(mnemonics, sub, &db, logger)
	if err != nil {
		return farmerBot, err
	}

	farmerBot.db = db
	farmerBot.rmbNodeClient = rmbNodeClient
	farmerBot.powerManager = powerManager
	farmerBot.logger = logger
	return farmerBot, nil
}

// Run runs farmerbot to update nodes and power management
func (f *FarmerBot) Run(ctx context.Context) {
	f.logger.Info().Msg("Starting farmer bot...")
	// TODO: change to 5 * time.Minute
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		startTime := time.Now()

		// update nodes
		f.logger.Debug().Msgf("get DB nodes")
		nodes, err := f.db.GetNodes()
		if err != nil {
			f.logger.Error().Err(err).Msg("failed to get nodes from db")
		}

		for _, node := range nodes {
			f.logger.Debug().Msgf("ping node with ID %v", node.ID)
			pong, err := f.rmbNodeClient.pingNode(ctx, node)

			if err != nil {
				f.logger.Error().Err(err).Msgf("failed to ping node with ID %d", node.ID)
				continue
			}

			if !pong {
				continue
			}

			f.logger.Debug().Msgf("update node with ID %v", node.ID)
			if err := f.rmbNodeClient.updateNode(ctx, node); err != nil {
				f.logger.Error().Err(err).Msgf("failed to update node with ID %d", node.ID)
				continue
			}

			if err := f.db.UpdatesNodes(node); err != nil {
				f.logger.Error().Err(err).Msgf("failed to update node %d in DB", node.ID)
				continue
			}
		}

		// wake up a new node in the wakeup time
		f.logger.Debug().Msg("check periodic wakeup")
		err = f.powerManager.PeriodicWakeup()
		if err != nil {
			f.logger.Error().Err(err).Msgf("failed to perform periodic wake up")
		}

		// TODO: add commands for power management and PeriodicWakeup
		// power management
		f.logger.Debug().Msg("check power management")
		err = f.powerManager.PowerManagement()
		if err != nil {
			f.logger.Error().Err(err).Msgf("failed to power management nodes")
		}

		delta := time.Since(startTime)
		f.logger.Debug().Msgf("Elapsed time for update: %v minutes", delta.Minutes())
	}
}

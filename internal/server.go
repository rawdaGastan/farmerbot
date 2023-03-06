// Package internal for farmerbot internals
package internal

import (
	"context"
	"fmt"

	"github.com/rawdaGastan/farmerbot/internal/constants"
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
	"github.com/threefoldtech/substrate-client"
	"github.com/threefoldtech/zbus"
)

// RunServer for running farmerbot server
func RunServer(mnemonics, network, redisAddr, version string, logger zerolog.Logger) error {
	const module = "farmerbot"
	server, err := zbus.NewRedisServer(module, fmt.Sprintf("tcp://%s", redisAddr), 10)
	if err != nil {
		return err
	}

	substrateManager := substrate.NewManager(constants.SubstrateURLs[network]...)
	subConn, err := substrateManager.Substrate()
	if err != nil {
		return err
	}

	db := models.NewRedisDB(redisAddr)

	farmManager := manager.NewFarmManager(&db, logger)
	nodeManager, err := manager.NewNodeManager(mnemonics, subConn, &db, logger)
	if err != nil {
		return err
	}
	powerManager, err := manager.NewPowerManager(mnemonics, subConn, &db, logger)
	if err != nil {
		return err
	}

	err = server.Register(zbus.ObjectID{Name: "farmmanager", Version: zbus.Version(version)}, &farmManager)
	if err != nil {
		return err
	}

	err = server.Register(zbus.ObjectID{Name: "powermanager", Version: zbus.Version(version)}, &powerManager)
	if err != nil {
		return err
	}

	err = server.Register(zbus.ObjectID{Name: "nodemanager", Version: zbus.Version(version)}, &nodeManager)
	if err != nil {
		return err
	}

	return server.Run(context.Background())
}

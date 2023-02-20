// Package cmd for farmerbot commands
package cmd

import (
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/spf13/cobra"
)

var powerOFFCmd = &cobra.Command{
	Use:   "poweroff",
	Short: "power off a node",
	Run: func(cmd *cobra.Command, args []string) {
		network, mnemonics, redisAddr, logger, err := getDefaultFlags(cmd)
		if err != nil {
			logger.Error().Err(err)
			return
		}

		nodeID, err := cmd.Flags().GetUint32("node")
		if err != nil || nodeID == 0 {
			logger.Error().Err(err).Msgf("error in node ID input %d", nodeID)
			return
		}
		logger.Debug().Msgf("node ID is: %v", nodeID)

		powerManager, err := manager.NewPowerManager(network, mnemonics, redisAddr, logger)
		if err != nil {
			logger.Error().Err(err).Msg("node manager failed to start")
			return
		}

		if err := powerManager.PowerOff(nodeID); err != nil {
			logger.Error().Err(err).Msgf("failed to power off node %d", nodeID)
			return
		}

		logger.Info().Msgf("Node %d is OFF", nodeID)
	},
}

func init() {
	cobra.OnInitialize()

	powerOFFCmd.Flags().Uint32P("node", "x", 0, "Enter your node ID to power on")

	powerOFFCmd.Flags().StringP("network", "n", "dev", "The network to run on")
	powerOFFCmd.Flags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	powerOFFCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	powerOFFCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	powerOFFCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

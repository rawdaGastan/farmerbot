// Package cmd for farmerbot commands
package cmd

import (
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var defineNodeCmd = &cobra.Command{
	Use:   "define",
	Short: "define a new node",
	Run: func(cmd *cobra.Command, args []string) {
		network, mnemonics, redisAddr, logger, err := getDefaultFlags(cmd)
		if err != nil {
			logger.Error().Err(err)
			return
		}

		config, err := cmd.Flags().GetString("config")
		if err != nil {
			logger.Error().Err(err).Msgf("error in config file path input '%s'", config)
			return
		}
		logger.Debug().Msgf("config path is: %v", config)

		nodeManager, err := manager.NewNodeManager(network, mnemonics, redisAddr, logger)
		if err != nil {
			logger.Error().Err(err).Msg("node manager failed to start")
			return
		}

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read config file %s", config)
			return
		}

		err = nodeManager.Define(jsonContent)
		if err != nil {
			logger.Error().Err(err).Msg("failed to define node")
			return
		}

		logger.Info().Msgf("Node is defined successfully")
	},
}

func init() {
	cobra.OnInitialize()

	defineNodeCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")

	defineNodeCmd.Flags().StringP("network", "n", "dev", "The network to run on")
	defineNodeCmd.Flags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	defineNodeCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	defineNodeCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	defineNodeCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

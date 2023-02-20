// Package cmd for farmerbot commands
package cmd

import (
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var findNodeCmd = &cobra.Command{
	Use:   "findnode",
	Short: "find an available node",
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

		excludes, err := cmd.Flags().GetUintSlice("exclude")
		if err != nil {
			logger.Error().Err(err).Msgf("error in exclude nodes input %v", excludes)
			return
		}

		nodeManager, err := manager.NewNodeManager(network, mnemonics, redisAddr, logger)
		if err != nil {
			logger.Error().Err(err).Msg("node manager failed to start")
			return
		}

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			logger.Error().Err(err).Msg("failed to read config file")
			return
		}

		options, err := parser.ParseJSONIntoNodeOptions(jsonContent)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get node options from config file")
			return
		}

		node, err := nodeManager.FindNode(options, excludes)
		if err != nil {
			logger.Error().Err(err).Msg("failed to find as node")
			return
		}

		logger.Info().Msgf("Node is %d", node.ID)
	},
}

func init() {
	cobra.OnInitialize()

	findNodeCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")
	findNodeCmd.Flags().UintSliceP("exclude", "x", []uint{}, "Enter your nodes list to be excluded from result")

	findNodeCmd.Flags().StringP("network", "n", "dev", "The network to run on")
	findNodeCmd.Flags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	findNodeCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	findNodeCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	findNodeCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

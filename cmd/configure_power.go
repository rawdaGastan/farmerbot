// Package cmd for farmerbot commands
package cmd

import (
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var configurePowerCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure farmerbot power",
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

		powerManager, err := manager.NewPowerManager(network, mnemonics, redisAddr, logger)
		if err != nil {
			logger.Error().Err(err).Msg("node manager failed to start")
			return
		}

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read config file %s", config)
			return
		}

		err = powerManager.Configure(jsonContent)
		if err != nil {
			logger.Error().Err(err).Msg("failed to configure power")
			return
		}

		logger.Info().Msgf("Power is configured successfully")
	},
}

func init() {
	cobra.OnInitialize()

	configurePowerCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")

	configurePowerCmd.Flags().StringP("network", "n", "dev", "The network to run on")
	configurePowerCmd.Flags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	configurePowerCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	configurePowerCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	configurePowerCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

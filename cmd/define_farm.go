// Package cmd for farmerbot commands
package cmd

import (
	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var defineFarmCmd = &cobra.Command{
	Use:   "define",
	Short: "define farmerbot farm",
	Run: func(cmd *cobra.Command, args []string) {
		_, _, redisAddr, logger, err := getDefaultFlags(cmd)
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

		farmManager := manager.NewFarmManager(redisAddr, logger)

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read config file %s", config)
			return
		}

		err = farmManager.Define(jsonContent)
		if err != nil {
			logger.Error().Err(err).Msg("failed to define farm")
			return
		}

		logger.Info().Msgf("Farm is defined successfully")
	},
}

func init() {
	cobra.OnInitialize()

	defineFarmCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")

	defineFarmCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	defineFarmCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	defineFarmCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

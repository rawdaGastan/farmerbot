// Package cmd for farmerbot commands
package cmd

import (
	"fmt"

	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var defineFarmCmd = &cobra.Command{
	Use:   "define",
	Short: "define farmerbot farm",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, _, _, db, logger, err := getDefaultFlags(cmd)
		if err != nil {
			return err
		}

		config, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("error %w in config file path input '%s'", err, config)
		}
		logger.Debug().Msgf("config path is: %v", config)

		farmManager := manager.NewFarmManager(&db, logger)

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s' with error: %w", config, err)
		}

		err = farmManager.Define(jsonContent)
		if err != nil {
			return fmt.Errorf("failed to define farm with error: %w", err)
		}

		logger.Info().Msgf("Farm is defined successfully")
		return nil
	},
}

func init() {
	cobra.OnInitialize()
	defineFarmCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")
}

// Package cmd for farmerbot commands
package cmd

import (
	"fmt"

	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var defineNodeCmd = &cobra.Command{
	Use:   "define",
	Short: "define a new node",
	RunE: func(cmd *cobra.Command, args []string) error {
		subConn, _, mnemonics, db, logger, err := getDefaultFlags(cmd)
		if err != nil {
			return err
		}

		config, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("error %w in config file path input '%s'", err, config)
		}
		logger.Debug().Msgf("config path is: %v", config)

		nodeManager, err := manager.NewNodeManager(mnemonics, subConn, &db, logger)
		if err != nil {
			return fmt.Errorf("node manager failed to start with error: %w", err)
		}

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s' with error: %w", config, err)
		}

		err = nodeManager.Define(jsonContent)
		if err != nil {
			return fmt.Errorf("failed to define node with error: %w", err)
		}

		logger.Info().Msgf("Node is defined successfully")
		return nil
	},
}

func init() {
	cobra.OnInitialize()
	defineNodeCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")
}

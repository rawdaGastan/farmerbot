// Package cmd for farmerbot commands
package cmd

import (
	"fmt"

	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/rawdaGastan/farmerbot/internal/parser"
	"github.com/spf13/cobra"
)

var findNodeCmd = &cobra.Command{
	Use:   "findnode",
	Short: "find an available node",
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

		excludes, err := cmd.Flags().GetUintSlice("exclude")
		if err != nil {
			return fmt.Errorf("error %w in exclude nodes input %v", err, excludes)
		}

		nodeManager, err := manager.NewNodeManager(mnemonics, subConn, &db, logger)
		if err != nil {
			return fmt.Errorf("node manager failed to start with error: %w", err)
		}

		jsonContent, err := parser.ReadFile(config)
		if err != nil {
			return fmt.Errorf("failed to read config file '%s' with error: %w", config, err)
		}

		options, err := parser.ParseJSONIntoNodeOptions(jsonContent)
		if err != nil {
			return fmt.Errorf("failed to get node options from config file with error: %w", err)
		}

		node, err := nodeManager.FindNode(options, excludes)
		if err != nil {
			return fmt.Errorf("failed to find a node with error: %w", err)
		}

		logger.Info().Msgf("Node is %d", node.ID)
		return nil
	},
}

func init() {
	cobra.OnInitialize()

	findNodeCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")
	findNodeCmd.Flags().UintSliceP("exclude", "x", []uint{}, "Enter your nodes list to be excluded from result")
}

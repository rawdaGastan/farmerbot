// Package cmd for farmerbot commands
package cmd

import (
	"fmt"

	manager "github.com/rawdaGastan/farmerbot/internal/managers"
	"github.com/spf13/cobra"
)

var powerONCmd = &cobra.Command{
	Use:   "poweron",
	Short: "power on a node",
	RunE: func(cmd *cobra.Command, args []string) error {
		subConn, _, mnemonics, db, logger, err := getDefaultFlags(cmd)
		if err != nil {
			return err
		}

		nodeID, err := cmd.Flags().GetUint32("node")
		if err != nil || nodeID == 0 {
			return fmt.Errorf("error %w in node ID input %d", err, nodeID)
		}
		logger.Debug().Msgf("node ID is: %v", nodeID)

		powerManager, err := manager.NewPowerManager(mnemonics, subConn, &db, logger)
		if err != nil {
			return fmt.Errorf("power manager failed to start with error %w", err)
		}

		if err := powerManager.PowerOn(nodeID); err != nil {
			return fmt.Errorf("failed to power on node %d with error: %w", nodeID, err)
		}

		logger.Info().Msgf("Node %d is ON", nodeID)
		return nil
	},
}

func init() {
	cobra.OnInitialize()
	powerONCmd.Flags().Uint32P("node", "x", 0, "Enter your node ID to power on")
}

// Package cmd for farmerbot commands
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Run farmerbot version to get it",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("Version is %s", version)
	},
}

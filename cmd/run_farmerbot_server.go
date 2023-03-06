// Package cmd for farmerbot commands
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/rawdaGastan/farmerbot/internal"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run farmerbot server to manage commands",
	Long:  `Welcome to the farmerbot (v0.0.0). The farmerbot is a service that a farmer can run allowing him to automatically manage the nodes of his farm.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, network, mnemonics, redisAddr, logger, err := getDefaultFlags(cmd)
		if err != nil {
			return err
		}

		err = internal.RunServer(mnemonics, network, redisAddr, version, logger)
		if err != nil {
			return err
		}

		return nil
	},
}

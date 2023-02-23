// Package cmd for farmerbot commands
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rawdaGastan/farmerbot/internal"
	"github.com/rawdaGastan/farmerbot/internal/constants"
	"github.com/rawdaGastan/farmerbot/internal/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/threefoldtech/substrate-client"
)

// farmerBotCmd represents the root base command when called without any subcommands
var farmerBotCmd = &cobra.Command{
	Use:   "farmerbot",
	Short: "Run farmerbot to manage your farms",
	Long:  `Farmerbot is a service that a farmer can run allowing him to automatically manage the nodes of his farm.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		subConn, network, mnemonics, db, logger, err := getDefaultFlags(cmd)
		if err != nil {
			return err
		}

		config, err := cmd.Flags().GetString("config")
		if err != nil {
			return fmt.Errorf("error in config file path input '%s'", config)
		}
		logger.Debug().Msgf("config path is: %v", config)

		farmerBot, err := internal.NewFarmerBot(config, network, mnemonics, subConn, db, logger)
		if err != nil {
			return fmt.Errorf("farmerbot failed to start")
		}

		farmerBot.Run(cmd.Context())
		return nil
	},
}

var nodeManagerCmd = &cobra.Command{
	Use:   "nodemanager",
	Short: "node manager command",
}

var farmManagerCmd = &cobra.Command{
	Use:   "farmmanager",
	Short: "farm manager command",
}

var powerManagerCmd = &cobra.Command{
	Use:   "powermanager",
	Short: "power manager command",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	farmerBotCmd.AddCommand(nodeManagerCmd)
	farmerBotCmd.AddCommand(farmManagerCmd)
	farmerBotCmd.AddCommand(powerManagerCmd)

	nodeManagerCmd.AddCommand(findNodeCmd)
	powerManagerCmd.AddCommand(powerONCmd)
	powerManagerCmd.AddCommand(powerOFFCmd)

	nodeManagerCmd.AddCommand(defineNodeCmd)
	farmManagerCmd.AddCommand(defineFarmCmd)
	powerManagerCmd.AddCommand(configurePowerCmd)

	err := farmerBotCmd.Execute()
	if err != nil {
		log.Err(err).Send()
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	farmerBotCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")

	farmerBotCmd.PersistentFlags().StringP("network", "n", "dev", "The network to run on")
	farmerBotCmd.PersistentFlags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	farmerBotCmd.PersistentFlags().StringP("redis", "r", "", "The address of the redis db")
	farmerBotCmd.PersistentFlags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	farmerBotCmd.PersistentFlags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

func getDefaultFlags(cmd *cobra.Command) (subConn *substrate.Substrate, network string, mnemonics string, db models.RedisDB, logger zerolog.Logger, err error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var debug bool
	debug, err = cmd.Flags().GetBool("debug")
	if err != nil {
		log.Error().Err(err).Msgf("error in debug mode input '%v'", debug)
		return
	}

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	log.Debug().Msgf("debug mode is: %v", debug)

	var logPath string
	logPath, err = cmd.Flags().GetString("log")
	if err != nil {
		log.Error().Err(err).Msgf("error in log file path input '%s'", logPath)
		return
	}
	log.Debug().Msgf("log path is: %v", logPath)

	var logFile *os.File
	logFile, err = os.OpenFile(
		logPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		log.Error().Err(err).Msgf("error in log file %v", err)
		return
	}

	multiWriter := zerolog.MultiLevelWriter(os.Stdout, logFile)
	logger = zerolog.New(multiWriter).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})

	redisAddr, err := cmd.Flags().GetString("redis")
	if err != nil {
		logger.Error().Err(err).Msgf("error in redis address input '%s'", redisAddr)
		return
	}

	if len(strings.TrimSpace(redisAddr)) == 0 {
		logger.Error().Msg("redis address is required")
		return
	}
	logger.Debug().Msgf("redis address is: %v", redisAddr)

	db = models.NewRedisDB(redisAddr)

	network, err = cmd.Flags().GetString("network")
	if err != nil {
		// we use it for farm manager and it doesn't use network nor mnemonics so we return
		if err.Error() == "flag accessed but not defined: network" {
			err = nil
			return
		}

		logger.Error().Err(err).Msgf("error in network input '%s'", network)
		return
	}
	logger.Debug().Msgf("network is: %v", strings.ToUpper(network))

	substrateManager := substrate.NewManager(constants.SubstrateURLs[network]...)
	subConn, err = substrateManager.Substrate()
	if err != nil {
		err = fmt.Errorf("error %w getting substrate connection using %v", err, constants.SubstrateURLs[network])
		return
	}

	mnemonics, err = cmd.Flags().GetString("mnemonics")
	if err != nil {
		logger.Error().Err(err).Msgf("error in mnemonics input '%s'", mnemonics)
		return
	}

	if len(strings.TrimSpace(mnemonics)) == 0 {
		logger.Error().Msg("mnemonics is required")
		return
	}
	logger.Debug().Msgf("mnemonics is: %v", mnemonics)

	return
}

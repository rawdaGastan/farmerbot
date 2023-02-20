// Package cmd for farmerbot commands
/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/rawdaGastan/farmerbot/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// farmerBotCmd represents the root base command when called without any subcommands
var farmerBotCmd = &cobra.Command{
	Use:   "farmerbot",
	Short: "Run farmerbot to manage your farms",
	Long:  `Farmerbot is a service that a farmer can run allowing him to automatically manage the nodes of his farm.`,
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

		farmerBot, err := internal.NewFarmerBot(config, network, mnemonics, redisAddr, logger)
		if err != nil {
			logger.Error().Err(err).Msg("farmerbot failed to start")
			return
		}

		farmerBot.Run(cmd.Context())
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
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	farmerBotCmd.Flags().StringP("config", "c", "config.json", "Enter your config json file path")

	farmerBotCmd.Flags().StringP("network", "n", "dev", "The network to run on")
	farmerBotCmd.Flags().StringP("mnemonics", "m", "", "The mnemonics of the farmer")
	farmerBotCmd.Flags().StringP("redis", "r", "", "The address of the redis db")
	farmerBotCmd.Flags().BoolP("debug", "d", false, "By setting this flag the farmerbot will print debug logs too")
	farmerBotCmd.Flags().StringP("log", "l", "farmerbot.log", "Enter your log file path to debug")
}

func getDefaultFlags(cmd *cobra.Command) (network string, mnemonics string, redisAddr string, logger zerolog.Logger, err error) {
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

	redisAddr, err = cmd.Flags().GetString("redis")
	if err != nil {
		logger.Error().Err(err).Msgf("error in redis address input '%s'", redisAddr)
		return
	}

	if len(strings.Trim(redisAddr, " ")) == 0 {
		logger.Error().Msg("redis address is required")
		return
	}
	logger.Debug().Msgf("redis address is: %v", redisAddr)

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

	mnemonics, err = cmd.Flags().GetString("mnemonics")
	if err != nil {
		logger.Error().Err(err).Msgf("error in mnemonics input '%s'", mnemonics)
		return
	}

	if len(strings.Trim(mnemonics, " ")) == 0 {
		logger.Error().Msg("mnemonics is required")
		return
	}
	logger.Debug().Msgf("mnemonics is: %v", mnemonics)

	return
}

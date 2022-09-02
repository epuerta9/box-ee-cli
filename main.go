package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var BuildVersion = "development"

const (
	configFile     string = ".box-ee.yaml"
	configFileType string = "yaml"
)

var (
	ErrorConfigNotFound = errors.New("config file not found in home directory or current directory. Run boxee init to get started")
	ErrorEmptyFlag      = errors.New("flag cannot be empty string")
)

func main() {
	setUpViper()
	var rootCmd = &cobra.Command{
		Use:   "boxee",
		Short: "Boxee Cli is a cli client for the Box-ee platform api",
		Long: `
		To learn more about usage and managing your box-ee account with cli visit the docs on the website
			`,
	}
	//registering all subcommands
	rootCmd.AddCommand(getInitCommand())
	rootCmd.AddCommand(getDeviceCmd())
	rootCmd.AddCommand(getTrackingCmd())
	rootCmd.AddCommand(getLoginCmd())
	rootCmd.AddCommand(getRegisterCmd())
	rootCmd.AddCommand(getRecoverCmd())
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setUpViper() {
	viper.SetConfigName(configFile)
	viper.SetConfigFile(configFile)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(".")
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "version of tool",
		Long: `
		check the current version of box-ee's boxee cli `,
		RunE: func(cmd *cobra.Command, args []string) error {
			version := os.Getenv("VERSION")
			if version == "" {
				version = "0.0.1-beta"
			}

			fmt.Println(version)
			return nil
		},
	}
}

func readConfig() error {
	file, err := os.Open(fmt.Sprintf("%v", configFile))
	if err != nil {
		return err
	}
	if err := viper.ReadConfig(file); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			//config file was not found
			return ErrorConfigNotFound

		} else {
			// Config file was found but another error was produced
			return err
		}
	}
	return nil

}

func CheckEmptyFlag(flag string) error {
	if flag == "" {
		return ErrorEmptyFlag
	} else {
		return nil
	}
}

func checkEmptyFlags(flags []string) error {
	for _, f := range flags {
		if f == "" {
			return ErrorEmptyFlag
		}
	}
	return nil

}

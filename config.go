package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Email    string
	Password string
	Address  string
)

func init() {
	adrs := os.Getenv("BOXEE_ADDRESS")
	if adrs == "" {
		return
	}
	Address = adrs

}

func getInitCommand() *cobra.Command {

	initCommand := &cobra.Command{
		Use:   "init",
		Short: "init boxee config in current directory",
		Long: `
			init creates a config file is it does not exist already. When creating a config file, passing certain required flags will generate a box-ee.yaml
			`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("running init config")
			//config file does not exist so create one with current flags set
			setUpViper()
			return viper.WriteConfig()

		},
	}
	initCommand.Flags().StringVarP(&Email, "email", "e", "", "email used to register")
	initCommand.Flags().StringVarP(&Address, "address", "a", "http://localhost:3000", "package place url address")

	//make email and password required
	initCommand.MarkFlagRequired("email")
	//bind viper to the flags passed in
	viper.BindPFlag("email", initCommand.Flags().Lookup("email"))
	viper.BindPFlag("address", initCommand.Flags().Lookup("address"))
	return initCommand
}

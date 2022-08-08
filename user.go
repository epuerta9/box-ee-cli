package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getLoginCmd() *cobra.Command {
	var password string
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login to package place",
		Long:  `login to the package place server using the config credentials`,
		RunE: func(cmd *cobra.Command, args []string) error {

			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}

			client.RequestEditors = append(client.RequestEditors, setRequestHeaders())
			ctx := context.TODO()
			resp, err := client.AdminLogin(ctx, AdminLoginRequest{
				Email:    cParams.Email,
				Password: password,
			})
			if err != nil {
				json.NewEncoder(os.Stdout).Encode(AdminLoginResponseItem{
					Msg:          err.Error(),
					SessionToken: "",
					StatusCode:   resp.StatusCode,
				})
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			var loginResponse AdminLoginResponseItem
			json.Unmarshal(body, &loginResponse)
			json.NewEncoder(os.Stdout).Encode(loginResponse)

			//write to config
			viper.Set("session_token", loginResponse.SessionToken)

			return viper.WriteConfig()

		},
	}
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password to login to package place api")
	loginCmd.MarkFlagRequired("password")
	return loginCmd

}
func getRegisterCmd() *cobra.Command {
	var password string
	var email string
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "register to package place",
		Long:  `register to the package place servers`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := readConfig(); err != nil {
				return err
			}
			//check email in config. If not equal then replace config email with new one
			if viper.Get("email") != email {
				viper.Set("email", email)
			}
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			//todo register needs an admin type request
			token, err := client.AdminRegister(ctx, AdminLoginRequest{
				Email:    cParams.Email,
				Password: password,
			})
			if err != nil {
				fmt.Println("error ", err)
				return err
			}
			viper.Set("session_token", token)
			fmt.Println("[+] session token ", token)
			return viper.WriteConfig()

		},
	}
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "password to login to package place api")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "email to login to package place api")
	registerCmd.MarkFlagRequired("password")
	registerCmd.MarkFlagRequired("email")
	return registerCmd

}

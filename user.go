package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getLoginCmd() *cobra.Command {
	var password string
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login to box-ee",
		Long:  `login to the box-ee server using the config credentials`,
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
		Short: "register to box-ee",
		Long:  `register to the box-ee servers`,
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
			_, err = client.AdminRegister(ctx, AdminLoginRequest{
				Email:    cParams.Email,
				Password: password,
			})
			if err != nil {
				fmt.Println("error ", err)
				return err
			}
			return nil

		},
	}
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "password to login to box-ee api")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "email to login to box-ee api")
	registerCmd.MarkFlagRequired("password")
	registerCmd.MarkFlagRequired("email")
	return registerCmd

}
func getRecoverCmd() *cobra.Command {
	var email string
	recoverCmd := &cobra.Command{
		Use:   "recover",
		Short: "recover account for box-ee",
		Long:  `recover password for box-ee servers`,
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
			resp, err := client.AdminRecover(ctx, AdminRecoverRequest{
				Email: cParams.Email,
			})
			if err != nil {
				fmt.Println("error ", err)
				return err
			}
			var r map[string]interface{}
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)
			json.Unmarshal(b, &r)
			fmt.Println(r)
			return nil

		},
	}
	recoverCmd.Flags().StringVarP(&email, "email", "e", "", "email to send recovery link")
	recoverCmd.MarkFlagRequired("email")
	return recoverCmd

}

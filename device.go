package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	deviceName string
	deviceType string
	deviceId   string
)

func getDeviceCmd() *cobra.Command {
	//device root command. Hang all sub commands related to device off of this one
	deviceCmd := &cobra.Command{
		Use:   "device",
		Short: "device actions command",
		Long: `
			The root command for device. Possible subcommands include add/list/delete/update`,
	}

	//add sub commands
	deviceCmd.AddCommand(deviceAdd())
	deviceCmd.AddCommand(deviceUpdate())
	deviceCmd.AddCommand(deviceList())
	deviceCmd.AddCommand(deviceDelete())
	deviceCmd.AddCommand(deviceGenerateKeys())
	return deviceCmd
}

func deviceAdd() *cobra.Command {
	deviceAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add a device",
		Long: `
		add a device. Required flags include device name and device type`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var createdResponse DeviceCreatedResponse
			//add a device
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.AddDevice(ctx, DeviceRequestAdd{
				DeviceName: deviceName,
				DeviceType: deviceType,
			})
			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &createdResponse)

			json.NewEncoder(os.Stdout).Encode(createdResponse)

			return nil
		},
	}
	deviceAddCmd.Flags().StringVarP(&deviceName, "name", "n", "default", "specify a device name")
	deviceAddCmd.Flags().StringVarP(&deviceType, "type", "t", "main", "specify a device type")
	return deviceAddCmd
}
func deviceGenerateKeys() *cobra.Command {
	deviceGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a device api key",
		Long: `
		generate a device client key used to setup the box-ee device. Required flags include device id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//generate device api key
			var keyGenResp DeviceKeyGenResponse
			if err := checkEmptyFlags([]string{deviceId}); err != nil {
				return err
			}
			if err := readConfig(); err != nil {
				return err
			}

			cParams := readValuesFromConfig()
			client, err := NewClient(cParams.Address)
			if err != nil {
				return err
			}
			ctx := context.TODO()
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))

			resp, err := client.GenKey(ctx, DeviceRequestKeyGen{
				DeviceId: deviceId,
			})

			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &keyGenResp)
			json.NewEncoder(os.Stdout).Encode(keyGenResp)
			return nil
		},
	}
	deviceGenerateCmd.Flags().StringVarP(&deviceId, "id", "i", "", "specify a device id")
	deviceGenerateCmd.MarkFlagRequired("id")
	return deviceGenerateCmd
}
func deviceDelete() *cobra.Command {
	deviceDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a device",
		Long: `
		delete a device. Required flags include device id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//update a device
			var deleteResp StandardResponse
			if err := checkEmptyFlags([]string{deviceId}); err != nil {
				return err
			}
			if err := readConfig(); err != nil {
				return err
			}
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))

			resp, err := client.DeleteDevice(ctx, &DeleteDeviceParams{
				DeviceId: deviceId,
			})
			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &deleteResp)
			json.NewEncoder(os.Stdout).Encode(deleteResp)
			return nil
		},
	}
	deviceDeleteCmd.Flags().StringVarP(&deviceId, "id", "i", "", "specify a device id")
	deviceDeleteCmd.MarkFlagRequired("id")
	return deviceDeleteCmd
}
func deviceUpdate() *cobra.Command {
	var toName string
	deviceUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "update a device",
		Long: `
		update a device. Required flags include device id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//update a device
			var updateResp StandardResponse
			if err := checkEmptyFlags([]string{toName, deviceId}); err != nil {
				return err
			}
			if err := readConfig(); err != nil {
				return err
			}
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.UpdateDevice(ctx, DeviceRequestPatch{
				DeviceId: deviceId,
				ToName:   toName,
			})
			if err != nil {
				json.NewEncoder(os.Stdout).Encode(err)
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &updateResp)
			json.NewEncoder(os.Stdout).Encode(updateResp)

			return nil
		},
	}
	deviceUpdateCmd.Flags().StringVarP(&deviceId, "id", "i", "", "specify a device id")
	deviceUpdateCmd.Flags().StringVarP(&toName, "to-name", "", "", "specify a device name")
	deviceUpdateCmd.MarkFlagRequired("id")
	deviceUpdateCmd.MarkFlagRequired("to-name")
	return deviceUpdateCmd
}
func deviceList() *cobra.Command {
	deviceListCmd := &cobra.Command{
		Use:   "list",
		Short: "list all devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			//listing out all devices
			var listResp ListDevices
			if err := readConfig(); err != nil {
				return err
			}
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.ListDevices(ctx)
			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &listResp)
			json.NewEncoder(os.Stdout).Encode(listResp)
			return nil
		},
	}

	return deviceListCmd
}

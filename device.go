package main

import (
	"context"
	"encoding/json"
	"fmt"
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
			//add a device
			if err := checkEmptyFlags([]string{deviceName, deviceType}); err != nil {
				return err
			}

			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()

			resp, err := client.AddDevice(ctx, DeviceRequestAdd{
				DeviceName: deviceName,
				DeviceType: deviceType,
			})
			if err != nil {
				return err
			}

			//			resp, err := httpclient.NewPClient().DevicePost(httpclient.DeviceRequest{
			//				DeviceName: deviceName,
			//				DeviceType: deviceType,
			//			})
			//			if err != nil {
			//				return err
			//			}
			fmt.Println(resp)
			json.NewEncoder(os.Stdout).Encode(resp)

			return nil
		},
	}
	deviceAddCmd.Flags().StringVarP(&deviceName, "name", "n", "", "specify a device name")
	deviceAddCmd.Flags().StringVarP(&deviceType, "type", "t", "", "specify a device type")
	deviceAddCmd.MarkFlagRequired("name")
	deviceAddCmd.MarkFlagRequired("type")
	return deviceAddCmd
}
func deviceGenerateKeys() *cobra.Command {
	deviceGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "generate a client key for device",
		Long: `
		generate a device client key used to setup the package place device. Required flags include device id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//update a device
			if err := checkEmptyFlags([]string{deviceId}); err != nil {
				return err
			}
			if err := readConfig(); err != nil {
				return err
			}

			cParams := readValuesFromConfig()
			client, err := NewClient(cParams.Address)
			ctx := context.TODO()
			resp, err := client.GenKey(ctx, DeviceRequestKeyGen{
				DeviceId: deviceId,
			})

			//			resp, err := httpclient.NewPClient().DeviceGenerate(httpclient.DeviceRequest{
			//				DeviceID: deviceId,
			//			})
			//			if err != nil {
			//				return err
			//}
			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)
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
			resp, err := client.DeleteDevice(ctx, &DeleteDeviceParams{
				DeviceId: deviceId,
			})
			//			resp, err := httpclient.NewPClient().DeviceDelete(httpclient.DeviceRequest{
			//				DeviceID: deviceId,
			//			})
			//			if err != nil {
			//				return err
			//			}

			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)
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
			resp, err := client.UpdateDevice(ctx, DeviceRequestPatch{
				DeviceId: deviceId,
				ToName:   toName,
			})
			//			resp, err := httpclient.NewPClient().DeviceUpdate(httpclient.DeviceRequest{
			//				DeviceID: deviceId,
			//				ToName:   toName,
			//			})
			if err != nil {
				json.NewEncoder(os.Stdout).Encode(err)
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)

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
			if err := readConfig(); err != nil {
				return err
			}
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()
			//pp := httpclient.NewPClient()
			resp, err := client.ListDevices(ctx)

			//resp, err := pp.DeviceList()
			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)
			return nil
		},
	}

	return deviceListCmd
}

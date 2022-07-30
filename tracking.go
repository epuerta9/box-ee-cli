package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
)

var (
	trackingNumber string
	trackingID     string
)

func getTrackingCmd() *cobra.Command {
	//tracking root command. Hang all sub commands related to tracking off of this one
	trackingCmd := &cobra.Command{
		Use:   "tracking",
		Short: "tracking actions command",
		Long: `
			The root command for tracking. Possible subcommands include add/list/delete/update`,
	}

	//add sub commands
	trackingCmd.AddCommand(trackingAdd())
	trackingCmd.AddCommand(trackingList())
	trackingCmd.AddCommand(trackingDelete())
	return trackingCmd
}

func trackingAdd() *cobra.Command {
	trackingAddCmd := &cobra.Command{
		Use:   "add",
		Short: "add a tracking",
		Long: `
		add a tracking. Required flags include tracking name and tracking type`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//add a tracking
			if err := checkEmptyFlags([]string{trackingNumber, deviceId}); err != nil {
				return err
			}
			readConfig()
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()

			resp, err := client.AddTracking(ctx, TrackingRequestItem{
				DeviceId:       deviceId,
				TrackingNumber: trackingNumber,
			})
			//resp, err := httpclient.NewPClient().TrackingPost(httpclient.TrackingRequest{
			//	TrackingNumber: trackingNumber,
			//	DeviceID:       deviceId,
			//})
			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)

			return nil
		},
	}
	trackingAddCmd.Flags().StringVarP(&trackingNumber, "tracking-number", "", "", "specify a tracking number")
	trackingAddCmd.Flags().StringVarP(&deviceId, "device-id", "", "", "specify a device id")
	trackingAddCmd.MarkFlagRequired("tracking-number")
	trackingAddCmd.MarkFlagRequired("device-id")
	return trackingAddCmd
}
func trackingDelete() *cobra.Command {
	var trackingID string
	trackingDeleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete a tracking",
		Long: `
		delete a tracking. Required flags include tracking id`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//update a tracking
			if err := checkEmptyFlags([]string{trackingID}); err != nil {
				return err
			}
			readConfig()
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()

			res, err := client.DeleteTracking(ctx, &DeleteTrackingParams{
				TrackingId: trackingID,
			})
			//res, err := httpclient.NewPClient().TrackingDelete(httpclient.TrackingRequest{
			//	TrackingID: trackingID,
			//})
			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(res)

			return nil
		},
	}
	trackingDeleteCmd.Flags().StringVarP(&trackingID, "id", "i", "", "specify a tracking id")
	trackingDeleteCmd.MarkFlagRequired("id")
	return trackingDeleteCmd
}

func trackingList() *cobra.Command {
	trackingListCmd := &cobra.Command{
		Use:   "list",
		Short: "list all trackings",
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

			resp, err := client.ListTrackings(ctx, &ListTrackingsParams{
				DeviceId: deviceId,
			})
			//	pp := httpclient.NewPClient()
			//	resp, err := pp.TrackingList(httpclient.TrackingRequest{
			//		DeviceID: deviceId,
			//	})
			if err != nil {
				return err
			}
			json.NewEncoder(os.Stdout).Encode(resp)
			return nil
		},
	}
	trackingListCmd.Flags().StringVarP(&deviceId, "device-id", "", "", "specify a device id")
	trackingListCmd.MarkFlagRequired("device-id")

	return trackingListCmd
}

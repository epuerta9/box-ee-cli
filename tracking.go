package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
			var addResp StandardResponse
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

			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.AddTracking(ctx, TrackingRequestItem{
				DeviceId:       deviceId,
				TrackingNumber: trackingNumber,
			})
			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &addResp)

			json.NewEncoder(os.Stdout).Encode(addResp)
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
			var deleteResp StandardResponse
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

			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.DeleteTracking(ctx, &DeleteTrackingParams{
				TrackingId: trackingID,
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
	trackingDeleteCmd.Flags().StringVarP(&trackingID, "id", "i", "", "specify a tracking id")
	trackingDeleteCmd.MarkFlagRequired("id")
	return trackingDeleteCmd
}

func trackingList() *cobra.Command {
	trackingListCmd := &cobra.Command{
		Use:   "list",
		Short: "list all trackings",
		RunE: func(cmd *cobra.Command, args []string) error {
			var listResp ListTrackings
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
			resp, err := client.ListTrackings(ctx, &ListTrackingsParams{
				DeviceId: deviceId,
			})
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
	trackingListCmd.Flags().StringVarP(&deviceId, "device-id", "", "", "specify a device id")
	trackingListCmd.MarkFlagRequired("device-id")

	return trackingListCmd
}

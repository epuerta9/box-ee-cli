package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	trackingNumber string
	trackingID     string
	file           string
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
	trackingCmd.AddCommand(trackingAddFile())

	return trackingCmd
}

func trackingAddFile() *cobra.Command {
	trackingfileCmd := &cobra.Command{
		Use:   "file",
		Short: "select tracking numbers file",
		Long: `
		add a tracking via file. Required flags include tracking name and tracking type`,
		RunE: func(cmd *cobra.Command, args []string) error {
			//add a tracking
			if err := checkEmptyFlags([]string{file}); err != nil {
				return err
			}
			readConfig()
			cParams := readValuesFromConfig()

			readFile, err := os.Open(file)

			if err != nil {
				fmt.Println(err)
			}
			fileScanner := bufio.NewScanner(readFile)
			fileScanner.Split(bufio.ScanLines)
			var fileLines []string

			for fileScanner.Scan() {
				fileLines = append(fileLines, fileScanner.Text())
			}

			readFile.Close()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()

			//check if device id is passed
			var allResponses []StandardResponse

			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			for _, tn := range fileLines {
				var payload TrackingRequestItem
				var addResp StandardResponse
				if deviceId == "" {
					payload = TrackingRequestItem{
						TrackingNumber: tn,
					}
				} else {
					payload = TrackingRequestItem{
						DeviceId:       &deviceId,
						TrackingNumber: tn,
					}
				}
				resp, err := client.AddTracking(ctx, payload)
				if err != nil {
					return err
				}
				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				json.Unmarshal(body, &addResp)
				allResponses = append(allResponses, addResp)
			}

			json.NewEncoder(os.Stdout).Encode(allResponses)
			return nil
		},
	}
	trackingfileCmd.Flags().StringVarP(&file, "file", "f", "", "specify a file")
	trackingfileCmd.Flags().StringVarP(&deviceId, "device-id", "", "", "specify a device id")
	trackingfileCmd.MarkFlagRequired("file")
	return trackingfileCmd
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
			if err := checkEmptyFlags([]string{trackingNumber}); err != nil {
				return err
			}
			readConfig()
			cParams := readValuesFromConfig()

			client, err := NewClient(cParams.Address)
			if err != nil {
				return err

			}
			ctx := context.TODO()

			//check if device id is passed
			var payload TrackingRequestItem
			if deviceId == "" {
				payload = TrackingRequestItem{
					TrackingNumber: trackingNumber,
				}
			} else {
				payload = TrackingRequestItem{
					DeviceId:       &deviceId,
					TrackingNumber: trackingNumber,
				}
			}
			client.RequestEditors = append(client.RequestEditors, setBoxeeAuthHeaders(cParams.SessionToken))
			resp, err := client.AddTracking(ctx, payload)
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

			var payload ListTrackingsParams
			if deviceId == "" {
				payload = ListTrackingsParams{}
			} else {
				payload = ListTrackingsParams{
					DeviceId: &deviceId,
				}
			}

			resp, err := client.ListTrackings(ctx, &payload)
			if err != nil {
				return err
			}
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			switch resp.StatusCode {
			case http.StatusOK:
				var listResp ListTrackings
				json.Unmarshal(body, &listResp)
				json.NewEncoder(os.Stdout).Encode(listResp)
			default:
				var defaultResp StandardResponse
				json.Unmarshal(body, &defaultResp)
				json.NewEncoder(os.Stdout).Encode(defaultResp)
			}

			return nil
		},
	}
	trackingListCmd.Flags().StringVarP(&deviceId, "device-id", "", "", "specify a device id")

	return trackingListCmd
}

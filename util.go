package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/spf13/viper"
)

type ConfigParams struct {
	Email        string
	SessionToken string
	Address      string
}

//
//func SerializeDeviceResponse(res httpclient.DeviceResponse) httpclient.DeviceResponse {
//	p, _ := json.Marshal(res)
//	var deviceResponse httpclient.DeviceResponse
//	json.Unmarshal(p, &deviceResponse)
//	return deviceResponse
//}
//
//func SerializeTrackingResponse(res httpclient.TrackingResponse) httpclient.TrackingResponse {
//	p, _ := json.Marshal(res)
//	var trackingResponse httpclient.TrackingResponse
//	json.Unmarshal(p, &trackingResponse)
//	return trackingResponse
//}

func checkAddressSet() bool {
	if Address == "" {
		return false

	}
	return true
}

func readValuesFromConfig() ConfigParams {
	readConfig()
	email := viper.GetString("email")
	address := viper.GetString("address")

	sessionToken := viper.GetString("session_token")

	return ConfigParams{
		Email:        email,
		Address:      address,
		SessionToken: sessionToken,
	}
}

func setRequestHeaders() RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Content-Type", "application/json")
		return nil
	}
}

func setBoxeeAuthHeaders(token string) RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		if token == "" {
			return errors.New("empty token in X-Boxee-Auth header")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Boxee-Auth", token)
		return nil
	}
}

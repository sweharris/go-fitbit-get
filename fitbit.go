package main

import (
	"context"
	"fmt"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/fitbit"
)

// OAuth2 configuration for Fitbit
func fitbit_config() oauth2.Config {
	return oauth2.Config{
		ClientID:     configuration.ClientID,
		ClientSecret: configuration.ClientSecret,
		RedirectURL:  fmt.Sprintf("http://localhost:%d", configuration.Port),
		Endpoint:     fitbit.Endpoint,

		// These are all the scopes according to
		// https://dev.fitbit.com/build/reference/web-api/developer-guide/application-design/#Scopes
		// as of 2025/05/05
		Scopes: []string{"activity", "cardio_fitness", "electrocardiogram", "heartrate", "irregular_rhythm_notifications", "location", "nutrition", "oxygen_saturation", "profile", "respiratory_rate", "settings", "sleep", "social", "temperature", "weight"},
	}
}

// Our access token
func fitbit_token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  configuration.AccessToken,
		RefreshToken: configuration.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       configuration.Expiry,
	}
}

// Call the FitBit API, check the token has been updated, save if needed
func call_fitbit_api(url string) []byte {
	// Get the oauth2 config and token
	conf := fitbit_config()
	token := fitbit_token()

	client := conf.Client(context.Background(), token)
	response, err := client.Get(url)
	check_or_die(err)

	check_refresh_token(client)

	body, err := io.ReadAll(response.Body)
	check_or_die(err)

	return body
}

package main

import (
	"fmt"
	"os"
)

func main() {
	// This will abort if the clientID/secret aren't specified
	// so we can assume they've been presented
	load_config()

	// If we don't have an access token we need to start our oauth
	// initial login process
	if configuration.AccessToken == "" {
		fmt.Fprintln(os.Stderr, "No access token found so we need to authenticate this application")
		fmt.Fprintln(os.Stderr, "")

		get_oauth2_token()
	}

	if len(os.Args) != 2 {
		die("Need a fitbit URL as the only parameter")
	}

	// body := call_fitbit_api("https://api.fitbit.com/1/user/-/activities/heart/date/2025-05-29/2025-05-29/1sec/time/00:00/23:59.json")
	body := call_fitbit_api(os.Args[1])

	fmt.Printf("%s", body)
}

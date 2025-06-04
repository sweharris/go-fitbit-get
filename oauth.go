package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// Function to handle the OAuth2 flow
// This is messy; we spin up a listener on the specified port and that
// handles the redirect from the Fitbit web page.  I'm sure there's a way
// of getting the code without this redirect process, but I dunno...
// It's possible the fitbit endpoints don't support OOB properly; at
// least I can't configure `urn:ietf:wg:oauth:2.0:oob` as an OOB URL
// (this appears to be deprecated across multiple services and requires
// this kludgy server instead.  Hmm)
//
// We need to pick a port that likely won't be used ephemeral but it seems
// that there isn't a guaranteed range
//    https://en.wikipedia.org/wiki/Ephemeral_port
// So we default to 16601  ("printf fitbit | sum") but let it be changed
// in the config.

var ouath2_state string

func get_oauth2_token() {
	// OAuth2 configuration for Fitbit
	var oauth2Config = fitbit_config()

	// This will let the function pause until the http handshake
	// completes
	channel := make(chan string)

	// We need a random string for "state" which is returned to the
	// web listener so the Oauth2 process can detect replays
	oauth2_state := uuid.NewString()
	authURL := oauth2Config.AuthCodeURL(oauth2_state, oauth2.AccessTypeOffline)
	fmt.Fprintln(os.Stderr, "Go to the following URL and authorize the application:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, authURL)

	// This is the web listener handler.  The response should
	// be .../?code=...  so we need to look for that value
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// First get the "code=" value
		code := r.URL.Query().Get("code")
		if code == "" {
			channel <- "Did not get code from redirect URL"
			http.Error(w, "Code is missing", http.StatusBadRequest)
			return
		}

		state := r.URL.Query().Get("state")
		if state != oauth2_state {
			channel <- "State from redirect URL does not match expected value"
			http.Error(w, "State is wrong", http.StatusBadRequest)
			return
		}

		// Now take that code and swap it for an access token
		token, err := oauth2Config.Exchange(context.Background(), code)
		if err != nil {
			channel <- "Failed to exchange code for token: " + err.Error()
			http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
			return
		}

		// Yay!  Save the results in the configuration
		configuration.AccessToken = token.AccessToken
		configuration.RefreshToken = token.RefreshToken
		configuration.Expiry = token.Expiry
		channel <- ""
	})

	// Start an HTTP server to handle the callback
	go http.ListenAndServe(fmt.Sprintf(":%d", configuration.Port), nil)

	// Wait for the handler to return a value
	msg := <-channel

	// This leaves the listener running until the program finishes;
	// I think this will be fine 'cos it's just for the duration of the
	// short CLI call.

	if msg != "" {
		fmt.Fprintf(os.Stderr, "Error during authentication process: %s\n", msg)
		os.Exit(255)
	}

	// Save this configuration
	save_config()
}

// After a call to Client.Get(...) the token may be refreshed, so we
// check this and if it's changed then update our saved value
func check_refresh_token(client *http.Client) {
	// There's gotta be a better way to get a new token!
	nt, err := client.Transport.(*oauth2.Transport).Source.Token()
	check_or_die(err)

	if nt.AccessToken != configuration.AccessToken ||
		nt.RefreshToken != configuration.RefreshToken ||
		nt.Expiry != configuration.Expiry {

		configuration.AccessToken = nt.AccessToken
		configuration.RefreshToken = nt.RefreshToken
		configuration.Expiry = nt.Expiry
		save_config()
		// fmt.Fprintf(os.Stderr, "Token refreshed")
	}
}

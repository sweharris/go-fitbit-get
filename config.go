package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/tkanos/gonfig"
)

// This is what we store in the config file
type Configuration struct {
	ClientID     string
	ClientSecret string
	Port         int
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// Global variables based on that config
// (Yeah, no pretense at data hiding or Objects here, just a global!)
var configuration Configuration

// Pretty Print a structure as JSON
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// Where we want the config file to be saved
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home + "\\"
	}
	return os.Getenv("HOME") + "/"
}

func config_file() string {
	return UserHomeDir() + ".fitbit_get"
}

func load_config() {
	parse := gonfig.GetConf(config_file(), &configuration)
	if parse != nil {
		fmt.Fprintln(os.Stderr, "Error parsing "+config_file()+"\n  ", parse)
		os.Exit(255)
	}

	if configuration.ClientID == "" {
		fmt.Fprintln(os.Stderr, "ClientID is not defined.  Aborted")
		os.Exit(255)
	}

	if configuration.ClientSecret == "" {
		fmt.Fprintln(os.Stderr, "ClientSecret is not defined.  Aborted")
		os.Exit(255)
	}

	if configuration.Port == 0 {
		configuration.Port = 16601
	}
}

func save_config() {
	data := []byte(prettyPrint(configuration))
	// fmt.Fprintln(os.Stderr, "Saving new config")
	err := os.WriteFile(config_file(), data, 0666)
	check_or_die(err)
}

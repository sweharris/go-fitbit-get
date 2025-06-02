package main

import (
	"fmt"
	"os"
)

func die(s string) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(255)
}

func check_or_die(err error) {
	if err != nil {
		die(err.Error())
	}
}

package main

import (
	"fmt"
	"os"
)

// Config represents the configuration of the gokku CLI client.
type Config struct {
	Username    string
	Hostname    string
	Port        int
	KeyFile     string `yaml:"keyfile,omitempty"`
	IgnoreAgent bool   `yaml:"ignore-agent,omitempty"`
}

var (
	gokkuConfig Config

	configFile = ".gokku.yml"
)

func check(err error, msg string, code int) {
	if err != nil {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(code)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "No arguments were passed to gokku, either run 'init' or a remote command.")
		os.Exit(1)
	}

	var ret int

	if args[0] == "init" {
		ret = runInit(args[1:])
	} else {
		if (args[0] == "-f" || args[0] == "--file") && len(args) > 2 {
			configFile = args[1]
			args = args[2:]
		}
		if args[0] == "--" {
			args = args[1:]
		}
		ret = runRemote(args)
	}

	os.Exit(ret)
}

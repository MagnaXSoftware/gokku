package main

import (
	"fmt"
	"os"
)

// ClientConfig represents the configuration of the gokku CLI client.
//
// The type gets marshalled as the "Gokku" key of an anonymous struct
// for the yaml .gokku.yml file.
type ClientConfig struct {
	Username    string
	Hostname    string
	Port        int
	KeyFile     string `yaml:"keyfile,omitempty"`
	IgnoreAgent bool   `yaml:"ignore-agent,omitempty"`
}

var CurrentConfig ClientConfig

var configFile = "./.gokku.yml"

const configFileName = ".gokku.yml"

// Command represents a Gokku CLI command.
type Command interface {
	Execute([]string) int
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
		ret = initCmd.Execute(args[1:])
	} else {
		ret = remoteCmd.Execute(args)
	}

	os.Exit(ret)
}

func check(err error, msg string, code int) {
	if err != nil {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, msg, err)
		os.Exit(code)
	}
}

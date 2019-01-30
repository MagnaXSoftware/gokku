package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var initLine = flag.NewFlagSet("init", flag.ExitOnError)

func init() {
	initLine.SetOutput(os.Stderr)
	initLine.StringVarP(&gokkuConfig.Username, "user", "u", "gokku", "Username used for the remote ssh connection, often gokku or dokku.")
	initLine.StringVarP(&gokkuConfig.Hostname, "hostname", "h", "", "The hostname of the gokku/dokku server.")
	initLine.IntVarP(&gokkuConfig.Port, "port", "p", 22, "The port on which the gokku/dokku binary can be reached (over ssh).")
	initLine.StringVarP(&gokkuConfig.KeyFile, "key", "k", "", "The path to the ssh key to use for the connection.")
	initLine.StringVarP(&configFile, "file", "f", ".gokku.yml", "The file that contains the configuration.")
	err := initLine.MarkHidden("file")
	if err != nil {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, "Could not configure flags: %v.\n", err)
		os.Exit(1)
	}
}

func runInit(args []string) int {
	// We ignore errors as initLine is set to ExitOnError
	//noinspection GoUnhandledErrorResult
	initLine.Parse(args)

	if gokkuConfig.Hostname == "" {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "A hostname is required.")
		initLine.PrintDefaults()
		os.Exit(3)
	}

	data, err := yaml.Marshal(struct{ Gokku *Config }{&gokkuConfig})
	check(err, "Could not export configuration file: %v.\n", 3)
	err = ioutil.WriteFile(configFile, data, 0664)
	check(err, "Could not write configuration file: %v.\n", 3)
	return 0
}

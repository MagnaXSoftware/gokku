package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)


type initCommand struct {
	FlagSet *flag.FlagSet
}

var initCmd = func() *initCommand{
	cmd := new(initCommand)

	cmd.FlagSet = flag.NewFlagSet("init", flag.ExitOnError)
	cmd.FlagSet.SetOutput(os.Stderr)
	cmd.FlagSet.StringVarP(&GokkuConfig.Username, "user", "u", "gokku", "Username used for the remote ssh connection, often gokku or dokku.")
	cmd.FlagSet.StringVarP(&GokkuConfig.Hostname, "hostname", "h", "", "The hostname of the gokku/dokku server.")
	cmd.FlagSet.IntVarP(&GokkuConfig.Port, "port", "p", 22, "The port on which the gokku/dokku binary can be reached (over ssh).")
	cmd.FlagSet.StringVarP(&GokkuConfig.KeyFile, "key", "k", "", "The path to the ssh key to use for the connection.")
	cmd.FlagSet.StringVarP(&configFile, "file", "f", ".gokku.yml", "The file that contains the configuration.")
	err := cmd.FlagSet.MarkHidden("file")
	if err != nil {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintf(os.Stderr, "Could not configure flags: %v.\n", err)
		os.Exit(1)
	}

	return cmd
}()

func (cmd *initCommand) Execute(args []string) int {
	// We ignore errors as cmd.FlagSet is set to ExitOnError
	//noinspection GoUnhandledErrorResult
	cmd.FlagSet.Parse(args)

	if GokkuConfig.Hostname == "" {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "A hostname is required.")
		cmd.FlagSet.PrintDefaults()
		os.Exit(3)
	}

	data, err := yaml.Marshal(struct{ Gokku *Config }{&GokkuConfig})
	check(err, "Could not export configuration file: %v.\n", 3)
	err = ioutil.WriteFile(configFile, data, 0664)
	check(err, "Could not write configuration file: %v.\n", 3)
	return 0
}

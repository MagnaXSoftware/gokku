package gokku

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var initCmd = newInitCommand()

type initCommand struct {
	FlagSet *flag.FlagSet

	NewConfig  ClientConfig
	OutputFile string
}

func (cmd *initCommand) Name() string {
	return "init"
}

func (cmd *initCommand) Execute(app *AppEnv, args []string) int {
	fmt.Printf("AppConfig:\n%#v\n\n", app.Config)
	fmt.Printf("NewConfig:\n%#v\n\n", cmd.NewConfig)

	// We ignore errors as cmd.FlagSet is set to ExitOnError
	//noinspection GoUnhandledErrorResult
	cmd.FlagSet.Parse(args)
	// Copy the config's value if a config exists and we haven't specified it on the CLI
	for _, f := range []string{"username", "hostname", "port", "key"} {
		changed := cmd.FlagSet.Changed(f)
		switch f {
		case "username":
			if !changed && app.HasConfig {
				cmd.NewConfig.Username = app.Config.Username
			}
		case "hostname":
			if !changed && app.HasConfig {
				cmd.NewConfig.Hostname = app.Config.Hostname
			}
		case "port":
			if !changed && app.HasConfig {
				cmd.NewConfig.Port = app.Config.Port
			}
		case "key":
			if !changed && app.HasConfig {
				cmd.NewConfig.KeyFile = app.Config.KeyFile
			}
		}
	}
	fmt.Printf("NewConfig (after parse):\n%#v\n\n", cmd.NewConfig)

	if cmd.NewConfig.Hostname == "" {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "A hostname is required.")
		cmd.FlagSet.PrintDefaults()
		return 3
	}

	data, err := yaml.Marshal(struct{ Gokku ClientConfig }{cmd.NewConfig})
	if err != nil {
		log.Printf("Could not export configuration file: %v.\n", err)
		return 3
	}

	if cmd.OutputFile == "" {
		cmd.OutputFile = app.ConfigFile
	}

	err = ioutil.WriteFile(cmd.OutputFile, data, 0664)
	if err != nil {
		log.Printf("Could not write configuration file: %v.\n", err)
		return 3
	}
	return 0
}

func newInitCommand() *initCommand {
	cmd := new(initCommand)

	cmd.FlagSet = flag.NewFlagSet("init", flag.ExitOnError)
	cmd.FlagSet.SetOutput(os.Stderr)
	cmd.FlagSet.StringVarP(&cmd.NewConfig.Username, "user", "u", "gokku", "Username used for the remote ssh connection, often gokku or dokku.")
	cmd.FlagSet.StringVarP(&cmd.NewConfig.Hostname, "hostname", "h", "", "The hostname of the gokku/dokku server.")
	cmd.FlagSet.IntVarP(&cmd.NewConfig.Port, "port", "p", 22, "The port on which the gokku/dokku binary can be reached (over ssh).")
	cmd.FlagSet.StringVarP(&cmd.NewConfig.KeyFile, "key", "k", "", "The path to the ssh key to use for the connection.")

	cmd.FlagSet.StringVar(&cmd.OutputFile, "output-file", "", "The path to the output file.")
	err := cmd.FlagSet.MarkHidden("output-file")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not configure flags: %v.\n", err)
		os.Exit(1)
	}

	return cmd
}

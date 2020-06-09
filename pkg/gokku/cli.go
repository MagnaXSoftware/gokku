package gokku

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func Cli(args []string) int {

	if len(args) < 2 {
		//noinspection GoUnhandledErrorResult
		fmt.Fprintln(os.Stderr, "No arguments were passed to gokku, either run 'init' or a remote command.")
		return 1
	}

	appEnv := NewAppEnv(args, remoteCmd, initCmd)

	return appEnv.Run()
}

// Command represents a CLI subcommand.
type Command interface {
	Name() string
	Execute(app *AppEnv, args []string) int
}

type AppEnv struct {
	AppName string
	Args    []string

	DefaultCmd Command
	Cmds       map[string]Command

	ConfigFile string
	HasConfig  bool
	Config     *ClientConfig
}

func NewAppEnv(args []string, defaultCmd Command, cmds ...Command) *AppEnv {
	a := &AppEnv{
		AppName:    args[0],
		Args:       args[1:],
		ConfigFile: DefaultConfigFile,
		Config:     new(ClientConfig),
		DefaultCmd: defaultCmd,
		Cmds:       make(map[string]Command),
	}

	for _, cmd := range cmds {
		a.Cmds[cmd.Name()] = cmd
	}

	return a
}

func (app *AppEnv) ParseConfig() error {
	fs := flag.NewFlagSet("cli", flag.ContinueOnError)

	fs.StringVarP(&app.ConfigFile, "file", "f", app.ConfigFile, "The file that contains the configuration.")

	err := fs.Parse(app.Args)
	if err != nil {
		return err
	}
	app.Args = fs.Args()

	if !fs.Changed("file") {
		app.ConfigFile, err = findConfigFile(app.ConfigFile)
		if err != nil {
			return nil
		}
	}

	data, err := ioutil.ReadFile(app.ConfigFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "config: unable to load config file %v: %v\n", app.ConfigFile, err)
		return nil
	}

	err = yaml.Unmarshal(data, struct{ Gokku *ClientConfig }{app.Config})
	if err != nil {
		return err
	}
	app.HasConfig = true

	if app.Config.Port == 0 {
		app.Config.Port = 22
	}

	return nil
}

func (app *AppEnv) Run() int {
	if err := app.ParseConfig(); err != nil {
		log.Printf("unable to load configuration: %v", err)
		return 2
	}

	if len(app.Args) >= 1 {
		if cmd, ok := app.Cmds[app.Args[0]]; ok {
			return cmd.Execute(app, app.Args[1:])
		}
	}

	return app.DefaultCmd.Execute(app, app.Args)
}

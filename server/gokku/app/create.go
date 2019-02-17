package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"magnax.ca/gokku/server/gokku"
	"os"
	"path"
)

var createCmd = &cobra.Command{
	Use: "create [app-name]",
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		return NameIsValid(args[0])
	},
	Run: create,
}

func init() {
	Plugin.Cmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) {
	appPath := path.Join(gokku.CurrentConfig.AppDirectory, args[0])
	_, err := os.Stat(appPath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("app: create: could not determine if app exists: %v\n", err)
			return
		}
	}

	fmt.Printf("app: create: initializing new app: %s\n", args[0])
	app := NewApp(args[0])
	err = app.Create()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

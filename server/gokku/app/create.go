package app

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"magnax.ca/gokku/server/gokku"
)

var createCmd = &cobra.Command{
	Use: "create [app-name]",
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		// check valid characters. We first check the whitelist for allowed characters,
		// then we extract all invalid characters if there are any.
		validRe := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
		if !validRe.MatchString(args[0]) {
			invalidRe := regexp.MustCompile("([^a-zA-Z0-9_-]+)")
			invalidRunes := strings.Join(invalidRe.FindAllString(args[0], -1), "")
			return fmt.Errorf("app: create: given app name contains invalid characters: %s", invalidRunes)
		}
		return nil
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

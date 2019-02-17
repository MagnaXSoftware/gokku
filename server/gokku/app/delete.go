package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var deleteCmd = &cobra.Command{
	Use:     "delete [app-name]",
	Aliases: []string{"remove"},
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.ExactArgs(1)(cmd, args)
		if err != nil {
			return err
		}
		return NameIsValid(args[0])
	},
	Run: remove,
}

func init() {
	Plugin.Cmd.AddCommand(deleteCmd)
}

func remove(cmd *cobra.Command, args []string) {
	app, err := LoadApp(args[0])
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("app: delete: app doesn't exist")
			return
		}
		fmt.Printf("app: delete: couldn't load app %s: %v\n", args[0], err)
		return
	}

	err = app.Destroy()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

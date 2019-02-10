package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"magnax.ca/gokku/server/gokku"
	"os"

	// core commands are added here. The plugin semantic, which commands must implement, handles registration.
	_ "magnax.ca/gokku/server/gokku/shell"
)

var RootCmd = &cobra.Command{
	Use: "gokku",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Root().Help()
	},
}

func Init() {
	gokku.InitPlugins(RootCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

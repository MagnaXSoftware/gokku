package app

import (
	"github.com/spf13/cobra"
	"magnax.ca/gokku/server/gokku"
)

// Implements the gokku.Plugin interface.
type appPlugin struct {
	Cmd *cobra.Command
}

// Plugin represents the app gokku.Plugin.
var Plugin = NewAppPlugin()

func init() {
	gokku.PrependPlugin(Plugin)
}

func NewAppPlugin() *appPlugin {
	plugin := new(appPlugin)
	plugin.Cmd = &cobra.Command{
		Use: "app",
	}
	return plugin
}

func (p *appPlugin) Name() string {
	return "app"
}

func (p *appPlugin) Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(p.Cmd)
}

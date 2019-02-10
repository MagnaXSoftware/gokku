package gokku

import "github.com/spf13/cobra"

// PluginList contains the list of all enabled plugins.
var PluginList []Plugin

// Plugin is the interface that gokku plugins must implement.
type Plugin interface {
	Name() string
	Init(rootCmd *cobra.Command)
}

// AppendPlugin adds a plugin to the end of the plugin list.
func AppendPlugin(plugin Plugin) {
	PluginList = append(PluginList, plugin)
}

// PrependPlugin adds a plugin to the front of the plugin list.
func PrependPlugin(plugin Plugin) {
	PluginList = append([]Plugin{plugin}, PluginList...)
}

// InitPlugins performs the initialization of all registered plugins, prior to call dispatch.
func InitPlugins(rootCmd *cobra.Command) {
	for _, plugin := range PluginList {
		plugin.Init(rootCmd)
	}
}

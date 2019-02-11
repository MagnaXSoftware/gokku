package main

import (
	"fmt"
	"io/ioutil"
	"magnax.ca/gokku/server/gokku"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const configDirectory = "/etc/gokku"
const pluginsFile = "plugins.yml"
const configFile = "gokku.yml"

// pluginFileList contains the list of all plugins after initialization.
//
// Each plugin must either be the full path to the plugin or a
// shortname for a plugin located in {$ConfigDir}/plugins.
var pluginFileList []string

func setupPlugins() {
	data, err := ioutil.ReadFile(path.Join(configDirectory, pluginsFile))
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("plugins: could not load plugin list: %v\n", err)
		}
		return
	}
	anon := struct{ Plugins []string }{pluginFileList}
	if yaml.Unmarshal(data, &anon) != nil {
		fmt.Printf("plugins: could not unmarshal plugin list: %v\n", err)
		os.Exit(1)
	}

	loadDynamicPlugins()
}

func loadDynamicPlugins() {

}

func loadConfig() {
	data, err := ioutil.ReadFile(path.Join(configDirectory, configFile))
	if err != nil {
		fmt.Printf("config: could not load config file: %v\n", err)
		os.Exit(1)
	}
	newCfg, err := gokku.CurrentConfig.UpdateFromMarshaled(data)
	if err != nil {
		fmt.Printf("config: could not parse config file: %v\n", err)
		os.Exit(1)
	}
	gokku.CurrentConfig = newCfg
}

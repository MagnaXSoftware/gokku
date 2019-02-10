package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// pluginsFile is the name of the file containing the plugin list.
var pluginsFile = "plugins.yaml"

// pluginFileList contains the list of all plugins after initialization.
//
// Each plugin must either be the full path to the plugin or a
// shortname for a plugin located in {$ConfigDir}/plugins.
var pluginFileList []string

func setupPlugins() {
	data, err := ioutil.ReadFile(path.Join(configDirectory, pluginsFile))
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("plugins: could not load plugin list: %v", err)
		}
		return
	}
	anon := struct{ Plugins []string }{pluginFileList}
	if yaml.Unmarshal(data, &anon) != nil {
		log.Fatalf("plugins: could not unmarshal plugin list: %v", err)
	}

	loadDynamicPlugins()
}

func loadDynamicPlugins() {

}

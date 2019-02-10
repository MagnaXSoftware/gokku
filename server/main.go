package main

import (
	"magnax.ca/gokku/server/cmd"
)

var (
	configDirectory = "/etc/gokku"
	configFile      = "gokku.yaml"
)

func main() {
	setupPlugins()

	cmd.Init()
	cmd.Execute()
}

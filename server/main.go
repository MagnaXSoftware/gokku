package main

import (
	"magnax.ca/gokku/server/cmd"
)


func main() {
	setupPlugins()
	loadConfig()

	cmd.Init()
	cmd.Execute()
}

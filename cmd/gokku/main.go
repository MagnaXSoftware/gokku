package main

import (
	"os"

	"magnax.ca/gokku/pkg/gokku"
)

func main() {
	os.Exit(gokku.Cli(os.Args))
}

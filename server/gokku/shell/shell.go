package shell

import (
	"github.com/spf13/cobra"
	"log"
	"magnax.ca/gokku/server/gokku"
	"os"
	"strings"
)

// Implements the gokku.Plugin interface.
type shellPlugin struct {
	Cmd *cobra.Command
}

// Plugin represents the shell gokku.Plugin.
var Plugin = NewAppPlugin()

func init() {
	gokku.AppendPlugin(Plugin)
}

// Executes the shell functionality.
//
// This bypasses the initial dispatch through cobra, however,
// a valid root command is still required.
func Exec() {
	_ = Plugin.exec()
}

func NewAppPlugin() *shellPlugin {
	plugin := new(shellPlugin)
	plugin.Cmd = &cobra.Command{
		Use: "shell",
		RunE: func(cmd *cobra.Command, args []string) error {
			return plugin.exec()
		},
	}
	return plugin
}

func (p *shellPlugin) Name() string {
	return "shell"
}

func (p *shellPlugin) Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(p.Cmd)
}

func (p *shellPlugin) exec() error {
	rootCmd := p.Cmd.Root()
	origCmd := cleanupCommand(os.Getenv("SSH_ORIGINAL_COMMAND"))

	if len(origCmd) == 0 {
		_ = rootCmd.Help()
		return nil
	}

	if IsGitCommand(origCmd[0]) {
		return RunGitCommand(origCmd)
	}

	// We replace the args with the "original" command.
	rootCmd.SetArgs(origCmd)
	// We absorb errors here, otherwise cobra thinks that this command failed, not the wrapped command.
	err := rootCmd.Execute()
	if err != nil {
		log.Printf("shell: %v", err)
	}
	return nil
}

func cleanupCommand(cmd string) []string {
	parts := strings.Split(cmd, " ")
	for {
		if len(parts) == 0 || parts[0] != "shell" {
			break
		}
		parts = parts[1:]
	}
	return parts
}

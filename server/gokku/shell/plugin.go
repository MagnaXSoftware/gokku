package shell

import (
	"fmt"
	"github.com/spf13/cobra"
	"magnax.ca/gokku/server/gokku"
	"os"
	"strings"
)

// Implements the gokku.Plugin interface.
type shellPlugin struct {
	Cmd *cobra.Command
}

// Plugin represents the shell gokku.Plugin.
var Plugin = NewShellPlugin()

func init() {
	gokku.AppendPlugin(Plugin)
}

// Exec runs the ssh shell.
//
// This bypasses the initial dispatch through cobra, though
// a valid root command is still required.
func Exec(user string) {
	_ = Plugin.exec(user)
}

func NewShellPlugin() *shellPlugin {
	plugin := new(shellPlugin)
	plugin.Cmd = &cobra.Command{
		Use: "shell [username]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			_ = plugin.exec(args[0])
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

func (p *shellPlugin) exec(user string) error {
	rootCmd := p.Cmd.Root()
	origCmd := cleanupCommand(os.Getenv("SSH_ORIGINAL_COMMAND"))
	err := os.Setenv("GOKKU_USER", user)
	if err != nil {
		fmt.Printf("shell: could not set gokku user: %v\n", err)
		return nil
	}

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
	err = rootCmd.Execute()
	if err != nil {
		fmt.Printf("shell: %v\n", err)
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

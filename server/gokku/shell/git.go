package shell

import (
	"fmt"
	"strings"
)

var GitCommands = []string{
	"git-upload-pack",
	"git-receive-pack",
}

func IsGitCommand(cmd string) bool {
	for _, gitCmd := range GitCommands {
		if gitCmd == cmd {
			return true
		}
	}
	return false
}

func RunGitCommand(cmd []string) error {
	return fmt.Errorf("shell: currently unable to handle %v\n", strings.Join(cmd, " "))
}

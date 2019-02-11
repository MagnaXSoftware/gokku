package shell

import (
	"fmt"
	"magnax.ca/gokku/server/gokku/app"
	"os"
	"path"
	"strings"
)

// GitError represents all git-related errors.
type GitError struct {
	Message string
}

// NewGitError creates a new GitError.
func NewGitError(message string) *GitError {
	return &GitError{
		Message: message,
	}
}

// Error implements the error interface on *GitError.
func (e *GitError) Error() string {
	return fmt.Sprintf("git: %s", e.Message)
}

// GitCommands represents the list of commands that are known to Gokku.
var GitCommands = []string{
	"git-upload-pack", // While we recognized git-upload-pack, we don't support it.
	"git-receive-pack",
}

// IsGitCommand determines if a given command is a known git command.
//
// The argument should be the first string in the args slice (args[0]).
func IsGitCommand(cmd string) bool {
	for _, gitCmd := range GitCommands {
		gitCmd = path.Base(gitCmd)
		if gitCmd == cmd {
			return true
		}
	}
	return false
}

// RunGitCommand executes the known git commands.
//
// The entire arguments, including the program name must be passed in.
// Giving a slice of less than 1 elements will result in an out-of-bounds panic.
func RunGitCommand(args []string) error {
	cmd := path.Base(args[0])
	switch cmd {
	case "git-upload-pack":
		writeString(PktLineString("ERR fetching is not supported"))
		return nil
	case "git-receive-pack":
		if len(args) < 2 {
			writeString(PktLineString("ERR no repository given"))
			return NewGitError("no repository given")
		}
		appName := strings.Split(path.Base(args[1]), ".")[0]
		currentApp, err := app.LoadApp(appName)
		if err != nil {
			writeString(PktLineString(fmt.Sprintf("ERR could not open app: %v", err)))
			return NewGitError("could not open app: " + err.Error())
		}
		writeString(PktLineString("ERR " + currentApp.RepositoryPath()))

		return nil
	}
	return NewGitError("unknown git command: " + strings.Join(args, " "))
}

func writeString(pktline string) {
	_, _ = os.Stdout.WriteString(pktline)
}

// PktLineString formats a given string into the pkt_line format expected by the git protocol.
func PktLineString(line string) string {
	return fmt.Sprintf("%04x%s\n", len(line)+5, line)
}

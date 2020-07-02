package gokku

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var remoteCmd = newRemoteCommand()

type remoteCommand struct{}

func (cmd *remoteCommand) Name() string {
	return "remote"
}

func (cmd *remoteCommand) Execute(app *AppEnv, args []string) int {
	if !app.HasConfig {
		_, _ = fmt.Fprintln(os.Stderr, "Gokku requires that a config file be present. Try running `gokku init`")
		return 1
	}

	config := &ssh.ClientConfig{
		User:            app.Config.Username,
		Auth:            buildAuthMethods(app.Config),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", app.Config.Hostname+":"+strconv.Itoa(app.Config.Port), config)
	if err != nil {
		log.Printf("remote: unable to connect to remote host: %v\n", err)
		return 1
	}

	session, err := client.NewSession()
	if err != nil {
		log.Printf("remote: unable to open new session: %v\n", err)
		return 1
	}
	//noinspection GoUnhandledErrorResult
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(strings.Join(args, " "))
	if eerr, ok := err.(*ssh.ExitError); ok {
		return eerr.ExitStatus()
	}
	if err != nil {
		log.Printf("remote: unable to execute given command: %v\n", err)
		return 1
	}

	return 0
}

func newRemoteCommand() *remoteCommand {
	cmd := new(remoteCommand)
	return cmd
}

func buildAuthMethods(conf *ClientConfig) []ssh.AuthMethod {
	var methods []ssh.AuthMethod

	if conf.KeyFile != "" {
		methods = append(methods, ssh.PublicKeys(conf.Key()))
	}

	socket := os.Getenv("SSH_AUTH_SOCK")
	if !conf.IgnoreAgent && socket != "" {
		conn, err := net.Dial("unix", socket)
		if err != nil {
			//noinspection GoUnhandledErrorResult
			fmt.Fprintf(os.Stderr, "remote: unable to connect with ssh-agent: %v, ignoring it\n", err)
			return methods
		}
		agentClient := agent.NewClient(conn)
		methods = append(methods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	return methods
}

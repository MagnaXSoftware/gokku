package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

type remoteCommand struct{}

var remoteCmd = func() *remoteCommand {
	cmd := new(remoteCommand)
	return cmd
}()

func (cmd *remoteCommand) Execute(args []string) int {
	args = findConfigFile(args)
	data, err := ioutil.ReadFile(configFile)
	check(err, "config: unable to open .gokku.yml file: %v\n", 1)

	anon := struct{ Gokku *GokkuConfig }{&Config}
	err = yaml.Unmarshal(data, &anon)
	check(err, "config: unable to import configuration data: %v\n", 1)

	if Config.Port == 0 {
		Config.Port = 22
	}

	return embeddedSSH(args)
}

func findConfigFile(args []string) []string {
	found := false
	if (args[0] == "-f" || args[0] == "--file") && len(args) > 2 {
		configFile = args[1]
		args = args[2:]
		found = true
	}
	if args[0] == "--" {
		args = args[1:]
	}
	if found {
		return args
	}

	// config file location was not passed in, so we recurse up to try to find one.
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("config: unable to get current working directory: %v", err)
		return args
	}
	for !found {
		if _, err = os.Stat(path.Join(cwd, configFileName)); err != nil {
			if cwd == path.Dir(cwd) {
				log.Fatalf("config: unable to find %v in the directory tree", configFileName)
			}
			if os.IsNotExist(err) {
				cwd = path.Dir(cwd)
				continue
			}
			log.Fatalf("config: unexpected error: %v", err)
		}
		found = true
		configFile = path.Join(cwd, configFileName)
	}
	return args
}

func readKey() ssh.Signer {
	key, err := ioutil.ReadFile(Config.KeyFile)
	check(err, "remote: unable to read private key: %v\n", 1)

	signer, err := ssh.ParsePrivateKey(key)
	check(err, "remote: unable to parse private key: %v\n", 1)
	return signer
}

func buildAuthMethods() []ssh.AuthMethod {
	var methods []ssh.AuthMethod

	if Config.KeyFile != "" {
		methods = append(methods, ssh.PublicKeys(readKey()))
	}

	socket := os.Getenv("SSH_AUTH_SOCK")
	if !Config.IgnoreAgent && socket != "" {
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

func embeddedSSH(args []string) int {
	config := &ssh.ClientConfig{
		User:            Config.Username,
		Auth:            buildAuthMethods(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", Config.Hostname+":"+strconv.Itoa(Config.Port), config)
	check(err, "remote: unable to connect to remote host: %v\n", 1)

	session, err := client.NewSession()
	check(err, "remote: unable to open new session: %v\n", 1)
	//noinspection GoUnhandledErrorResult
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(strings.Join(args, " "))
	if eerr, ok := err.(*ssh.ExitError); ok {
		return eerr.ExitStatus()
	}
	check(err, "remote: unable to execute given command: %v\n", 1)

	return 0
}

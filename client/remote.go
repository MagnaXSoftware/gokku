package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

func runRemote(args []string) int {
	data, err := ioutil.ReadFile(configFile)
	check(err, "unable to open .gokku.yml file: %v\n", 4)

	anon := struct{ Gokku *Config }{&gokkuConfig}
	err = yaml.Unmarshal(data, &anon)
	check(err, "unable to import configuration data: %v\n", 4)

	if gokkuConfig.Port == 0 {
		gokkuConfig.Port = 22
	}

	return embeddedSSH(args)
}

func readKey() ssh.Signer {
	key, err := ioutil.ReadFile(gokkuConfig.KeyFile)
	check(err, "unable to read private key: %v\n", 4)

	signer, err := ssh.ParsePrivateKey(key)
	check(err, "unable to parse private key: %v\n", 4)
	return signer
}

func buildAuthMethods() []ssh.AuthMethod {
	var methods []ssh.AuthMethod

	if gokkuConfig.KeyFile != "" {
		methods = append(methods, ssh.PublicKeys(readKey()))
	}

	socket := os.Getenv("SSH_AUTH_SOCK")
	if !gokkuConfig.IgnoreAgent && socket != "" {
		conn, err := net.Dial("unix", socket)
		if err != nil {
			//noinspection GoUnhandledErrorResult
			fmt.Fprintf(os.Stderr, "unable to connect with ssh-agent: %v, ignoring it\n", err)
			return methods
		}
		agentClient := agent.NewClient(conn)
		methods = append(methods, ssh.PublicKeysCallback(agentClient.Signers))
	}

	return methods
}

func embeddedSSH(args []string) int {
	config := &ssh.ClientConfig{
		User:            gokkuConfig.Username,
		Auth:            buildAuthMethods(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", gokkuConfig.Hostname+":"+strconv.Itoa(gokkuConfig.Port), config)
	check(err, "unable to connect to remote host: %v\n", 5)

	session, err := client.NewSession()
	check(err, "unable to open new session: %v\n", 5)
	//noinspection GoUnhandledErrorResult
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.Run(strings.Join(args, " "))
	check(err, "unable to execute given command: %v\n", 5)

	return 0
}

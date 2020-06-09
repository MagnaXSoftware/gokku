package gokku

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const DefaultConfigFile = ".gokku.yml"

// ClientConfig represents the configuration of the gokku CLI client.
//
// The type gets marshalled as the "Gokku" key of an anonymous struct
// for the yaml .gokku.yml file.
type ClientConfig struct {
	Username string
	Hostname string
	Port     int `yaml:"port,omitempty"`

	KeyFile     string `yaml:"keyfile,omitempty"`
	key         ssh.Signer
	IgnoreAgent bool `yaml:"ignore-agent,omitempty"`
}

func (c *ClientConfig) Key() ssh.Signer {
	if c.key == nil {
		key, err := ioutil.ReadFile(c.KeyFile)
		if err != nil {
			log.Fatalf("config: unable to read private key: %v\n", err)
		}

		c.key, err = ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("config: unable to parse private key: %v\n", err)
		}
	}

	return c.key
}

func findConfigFile(configFileName string) (string, error) {
	// config file location was not passed in, so we recurse up to try to find one.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err = os.Stat(path.Join(cwd, configFileName)); err != nil {
			if cwd == path.Dir(cwd) {
				// this works because path.Dir will return "/" when the input is "/"
				return "", fmt.Errorf("config: unable to find %v in the directory tree", configFileName)
			}
			if os.IsNotExist(err) {
				cwd = path.Dir(cwd)
				continue
			}
			return "", err
		}
		return path.Join(cwd, configFileName), nil
	}
}

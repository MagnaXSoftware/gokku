package gokku

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// GokkuConfig represents the configuration of Gokku.
type GokkuConfig struct {
	AppDirectory string
}

// Config is the current running configuration.
var Config = NewDefaultGokkuConfig()

func NewDefaultGokkuConfig() GokkuConfig {
	return GokkuConfig{
		AppDirectory: "/var/lib/gokku",
	}
}

func NewGokkuConfig() GokkuConfig {
	return GokkuConfig{}
}

func MarshalNewGokkuConfig(data []byte) (GokkuConfig, error) {
	cfg := NewGokkuConfig()
	anon := struct {
		Gokku *GokkuConfig
	}{&cfg}
	err := yaml.Unmarshal(data, &anon)
	if err != nil {
		return GokkuConfig{}, err
	}
	return cfg, nil
}

// UpdateFromMarshaled returns a copy of the GokkuConfig, updating values as per the given marshaled data.
func (c *GokkuConfig) UpdateFromMarshaled(data []byte) (GokkuConfig, error) {
	cfg := *c
	anon := struct {
		Gokku *GokkuConfig
	}{&cfg}
	err := yaml.Unmarshal(data, &anon)
	if err != nil {
		return GokkuConfig{}, err
	}
	return cfg, nil
}

func (c *GokkuConfig) Marshal() []byte {
	anon := struct {
		Gokku *GokkuConfig
	}{c}
	data, err := yaml.Marshal(anon)
	if err != nil {
		panic(fmt.Errorf("config: couldn't marshal config: %v", err))
	}
	return data
}

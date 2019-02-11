package gokku

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration of Gokku.
type Config struct {
	AppDirectory string
}

// CurrentConfig is the current running configuration.
var CurrentConfig = NewDefaultGokkuConfig()

func NewDefaultGokkuConfig() Config {
	return Config{
		AppDirectory: "/var/lib/gokku",
	}
}

func NewConfig() Config {
	return Config{}
}

func UnmarshalNewConfig(data []byte) (Config, error) {
	cfg := NewConfig()
	anon := struct {
		Gokku *Config
	}{&cfg}
	err := yaml.Unmarshal(data, &anon)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// UpdateFromMarshaled returns a copy of the Config, updating values as per the given marshaled data.
func (c *Config) UpdateFromMarshaled(data []byte) (Config, error) {
	cfg := *c
	anon := struct {
		Gokku *Config
	}{&cfg}
	err := yaml.Unmarshal(data, &anon)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) Marshal() []byte {
	anon := struct {
		Gokku *Config
	}{c}
	data, err := yaml.Marshal(anon)
	if err != nil {
		panic(fmt.Errorf("config: couldn't marshal config: %v", err))
	}
	return data
}

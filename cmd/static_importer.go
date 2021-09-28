package cmd

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"gopkg.in/yaml.v2"
)

// ImportToml import toml conf to static configuration.
func ImportToml(input io.Reader) (*static.Configuration, error) {
	conf := &static.Configuration{}
	_, err := toml.DecodeReader(input, conf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode static configuration from toml file: %w", err)
	}

	return conf, nil
}

// ImportYaml import yaml conf to static configuration.
func ImportYaml(input io.Reader) (*static.Configuration, error) {
	conf := &static.Configuration{}
	err := yaml.NewDecoder(input).Decode(conf)
	if err != nil {
		return nil, fmt.Errorf("cannot decode static configuration from yaml file: %w", err)
	}

	return conf, nil
}

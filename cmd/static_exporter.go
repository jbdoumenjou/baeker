package cmd

import (
	"fmt"

	"github.com/traefik/traefik/v2/pkg/provider/docker"

	"github.com/traefik/traefik/v2/pkg/config/static"
)

type MyConf static.Configuration

func NewMyConf() *MyConf {
	return &MyConf{Providers: &static.Providers{}}
}

func (m *MyConf) AddProvider(name string) {
	switch name {
	case "docker":

		if m.Providers.Docker != nil {
			fmt.Println("docker provider already exists")
			return
		}
		m.Providers.Docker = &docker.Provider{}

		return
	}
}

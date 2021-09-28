package cmd

import (
	"errors"
	"fmt"

	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/provider/docker"
	"github.com/traefik/traefik/v2/pkg/provider/file"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd"
)

// StaticConfBuilder store build configuration.
type StaticConfBuilder struct {
	conf *static.Configuration
}

// NewStaticConfBuilder creates a StaticConfBuilder.
func NewStaticConfBuilder() *StaticConfBuilder {
	return &StaticConfBuilder{
		conf: &static.Configuration{
			Providers:   &static.Providers{},
			EntryPoints: static.EntryPoints{},
		},
	}
}

// GetConfiguration returns the current configuration.
func (s StaticConfBuilder) GetConfiguration() *static.Configuration {
	return s.conf
}

// AddKubernetesProvider adds Kubernetes CRD Provider to the current configuration.
func (s StaticConfBuilder) AddKubernetesProvider() (*StaticConfBuilder, error) {
	if s.conf.Providers != nil && s.conf.Providers.KubernetesCRD != nil {
		return nil, errors.New("the KubernetesCRD provider already exists")
	}

	s.conf.Providers.KubernetesCRD = &crd.Provider{}

	return &s, nil
}

// AddDockerProvider adds Docker Provider to the current configuration.
func (s StaticConfBuilder) AddDockerProvider() (*StaticConfBuilder, error) {
	if s.conf.Providers != nil && s.conf.Providers.Docker != nil {
		return nil, errors.New("the Docker provider already exists")
	}

	s.conf.Providers.Docker = &docker.Provider{}

	return &s, nil
}

// AddFileProvider adds File Provider to the current configuration.
func (s StaticConfBuilder) AddFileProvider(directory string) (*StaticConfBuilder, error) {
	if s.conf.Providers != nil && s.conf.Providers.File != nil {
		return nil, errors.New("the File provider already exists")
	}

	s.conf.Providers.File = &file.Provider{
		Directory: directory,
	}

	return &s, nil
}

// AddEntryPoint adds an entrypoint to the current configuration.
func (s StaticConfBuilder) AddEntryPoint(name string, address string) (*StaticConfBuilder, error) {
	if _, ok := s.conf.EntryPoints[name]; ok {
		return nil, fmt.Errorf("EntryPoint %s already exists", name)
	}

	s.conf.EntryPoints[name] = &static.EntryPoint{Address: address}

	return &s, nil
}

package cmd

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"path"
	"sort"
	"strings"

	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/provider/docker"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd"
)

type traefikConf struct {
	Labels []string
	Ports  []entryPoint
}

type entryPoint struct {
	Name  string
	Value string
}

// GetDefaultConf generate a default configuration for a given provider.
func GetDefaultConf(provider string) static.Configuration {
	defaultConf := static.Configuration{
		EntryPoints: static.EntryPoints{
			"web":       {Address: ":8000"},
			"websecure": {Address: ":8443"},
		},
	}

	if defaultConf.Providers == nil {
		defaultConf.Providers = &static.Providers{}
	}

	if provider == "docker" {
		defaultConf.Providers.Docker = &docker.Provider{}
		defaultConf.Providers.Docker.SetDefaults()
	}

	if provider == "kubernetes" {
		defaultConf.Providers.KubernetesCRD = &crd.Provider{}
	}

	return defaultConf
}

// ExportConf export a configuration applying a specific template to the given output.
func ExportConf(conf static.Configuration, templatePath string, output io.Writer) error {
	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create the template: %w", err)
	}

	labels, err := getLabels(conf)
	if err != nil {
		return fmt.Errorf("failed to extract label from static configuration: %w", err)
	}

	ports, err := getPorts(conf.EntryPoints)
	if err != nil {
		return fmt.Errorf("failed to get ports from static configuration: %w", err)
	}

	err = tmpl.Execute(output, traefikConf{Labels: labels, Ports: ports})
	if err != nil {
		return fmt.Errorf("failed to execute the template: %w", err)
	}

	return nil
}

func getLabels(conf static.Configuration) ([]string, error) {
	var cleanedLabels []string
	labels, err := parser.Encode(conf, "")
	if err != nil {
		return cleanedLabels, fmt.Errorf("failed to parse the configuration: %w", err)
	}

	defaultLabels, err := getDefaultsLabel(conf)
	if err != nil {
		return cleanedLabels, fmt.Errorf("failed to get the default labels: %w", err)
	}

	for key, value := range labels {
		defaultValue, ok := defaultLabels[key]
		if ok && defaultValue == value {
			continue
		}

		if len(key) > 0 && len(value) > 0 {
			cleanedLabels = append(cleanedLabels, strings.ToLower(fmt.Sprintf("%s=%s", key[1:], value)))
		}
	}

	// TODO refactor, very naive approach
	if conf.Providers != nil {
		if conf.Providers.Docker != nil {
			cleanedLabels = append(cleanedLabels, "providers.docker")
		}
		if conf.Providers.KubernetesCRD != nil {
			cleanedLabels = append(cleanedLabels, "providers.kubernetescrd")
		}
	}

	// To keep the result consistent.
	sort.Strings(cleanedLabels)
	return cleanedLabels, nil
}

func getDefaultsLabel(conf static.Configuration) (map[string]string, error) {
	defaultConf := &static.Configuration{}
	if conf.Providers != nil {
		defaultConf.Providers = &static.Providers{}

		if conf.Providers.Docker != nil {
			defaultConf.Providers.Docker = &docker.Provider{}
			defaultConf.Providers.Docker.SetDefaults()
		}

		if conf.Providers.KubernetesCRD != nil {
			defaultConf.Providers.KubernetesCRD = &crd.Provider{}
		}
	}

	return parser.Encode(defaultConf, "")
}

func getPorts(entryPoints static.EntryPoints) ([]entryPoint, error) {
	var ports []entryPoint

	for name, entrypoint := range entryPoints {
		_, port, err := net.SplitHostPort(entrypoint.Address)
		if err != nil {
			return ports, fmt.Errorf("cannot process ports :%w", err)
		}

		ports = append(ports, entryPoint{
			Name:  name,
			Value: port,
		})
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Name < ports[j].Name
	})

	return ports, nil
}

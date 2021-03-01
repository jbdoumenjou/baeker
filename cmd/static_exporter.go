package cmd

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"path"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/traefik/paerser/parser"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"github.com/traefik/traefik/v2/pkg/provider/docker"
	"github.com/traefik/traefik/v2/pkg/provider/file"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd"
	"gopkg.in/yaml.v2"
)

type traefikConf struct {
	Labels []string
	Ports  []entryPoint
}

type entryPoint struct {
	Name  string
	Value string
}

func getDefaultsLabel(conf *static.Configuration) (map[string]string, error) {
	defaultConf := &static.Configuration{}
	if conf.Providers != nil {
		defaultConf.Providers = &static.Providers{}

		if conf.Providers.Docker != nil {
			defaultConf.Providers.Docker = &docker.Provider{}
		}

		if conf.Providers.KubernetesCRD != nil {
			defaultConf.Providers.KubernetesCRD = &crd.Provider{}
		}

		if conf.Providers.File != nil {
			defaultConf.Providers.File = &file.Provider{}
		}
	}

	labels, err := parser.Encode(defaultConf, "")
	if err != nil {
		return labels, fmt.Errorf("cannot encode default configuration: %w", err)
	}

	return labels, nil
}

func getLabels(conf *static.Configuration, prefix string) ([]string, error) {
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
			cleanedLabels = append(cleanedLabels, strings.ToLower(fmt.Sprintf("%s%s=%s", prefix, key[1:], value)))
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

// ExportToml exports static configuration to a toml format.
func ExportToml(config *static.Configuration, output io.Writer) error {
	if err := toml.NewEncoder(output).Encode(config); err != nil {
		// failed to encode
		return fmt.Errorf("cannot encode static configuration in TOML: %w", err)
	}

	return nil
}

// ExportYaml exports static configuration to a yaml format.
func ExportYaml(config *static.Configuration, output io.Writer) error {
	if err := yaml.NewEncoder(output).Encode(config); err != nil {
		// failed to encode
		return fmt.Errorf("cannot encode static configuration in YAML: %w", err)
	}

	return nil
}

// ExportCLI exports static configuration to a CLI format.
func ExportCLI(config *static.Configuration, output io.Writer) error {
	labels, err := getLabels(config, "--")
	if err != nil {
		return err
	}
	sort.Strings(labels)

	str := strings.Join(labels, " ")
	_, err = output.Write([]byte(str + "\n"))
	if err != nil {
		return fmt.Errorf("cannot write to the standard output: %w", err)
	}

	return nil
}

// ExportKubernetes export static configuration to a kubernetes crd format.
func ExportKubernetes(config *static.Configuration, templatePath string, output io.Writer) error {
	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create the template: %w", err)
	}

	var labels []string
	labels, err = getLabels(config, "")
	if err != nil {
		return err
	}

	ports, err := getPorts(config.EntryPoints)
	if err != nil {
		return fmt.Errorf("failed to get ports from static configuration: %w", err)
	}

	err = tmpl.Execute(output, traefikConf{Labels: labels, Ports: ports})
	if err != nil {
		return fmt.Errorf("failed to execute the template: %w", err)
	}

	return nil
}

// ExportDocker export static configuration to docker-compose format.
func ExportDocker(config *static.Configuration, templatePath string, output io.Writer) error {
	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create the template: %w", err)
	}

	var labels []string
	labels, err = getLabels(config, "")
	if err != nil {
		return err
	}

	ports, err := getPorts(config.EntryPoints)
	if err != nil {
		return fmt.Errorf("failed to get ports from static configuration: %w", err)
	}

	err = tmpl.Execute(output, traefikConf{Labels: labels, Ports: ports})
	if err != nil {
		return fmt.Errorf("failed to execute the template: %w", err)
	}

	return nil
}

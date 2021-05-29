package exporter

import (
	"fmt"
	"html/template"
	"io"
	"net"
	"path"
	"sort"
	"strings"

	"github.com/traefik/traefik/v2/pkg/provider/docker"

	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd"

	"github.com/traefik/paerser/parser"

	"github.com/BurntSushi/toml"
	"github.com/traefik/traefik/v2/pkg/config/static"
	"gopkg.in/yaml.v2"
)

func ExportToml(config static.Configuration, output io.Writer) error {
	if err := toml.NewEncoder(output).Encode(config); err != nil {
		// failed to encode
		return fmt.Errorf("Cannot encode static configuration in TOML: %w", err)
	}

	return nil
}

func ExportYaml(config static.Configuration, output io.Writer) error {
	if err := yaml.NewEncoder(output).Encode(config); err != nil {
		// failed to encode
		return fmt.Errorf("Cannot encode static configuration in YAML: %w", err)
	}

	return nil
}

func ExportCLI(config static.Configuration, output io.Writer) error {
	exported, err := parser.Encode(config, "")
	if err != nil {
		return err
	}

	var labels []string
	for k, v := range exported {
		labels = append(labels, fmt.Sprintf("--%s=%s", strings.ToLower(k[1:]), v))
	}
	sort.Strings(labels)

	str := strings.Join(labels, " ")
	output.Write([]byte(str + "\n"))

	return nil
}

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
	}

	return parser.Encode(defaultConf, "")
}

func getLabels(conf *static.Configuration) ([]string, error) {
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

func ExportKubernetes(config *static.Configuration, templatePath string, output io.Writer) error {
	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create the template: %w", err)
	}

	var labels []string
	labels, err = getLabels(config)
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

func ExportDocker(config *static.Configuration, templatePath string, output io.Writer) error {
	tmpl, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to create the template: %w", err)
	}

	var labels []string
	labels, err = getLabels(config)
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

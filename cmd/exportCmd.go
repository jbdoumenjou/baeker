package cmd

import (
	"fmt"
	"io"
	"os"
)

// ExportCmd Exports a yml static configuration file to standard output with a specified format.
func ExportCmd(input io.Reader, to string) error {
	// TODO: add other source type
	conf, err := ImportYaml(input)
	if err != nil {
		return fmt.Errorf("cannot import source file:%w", err)
	}

	switch to {
	case "cli":
		err := ExportCLI(conf, os.Stdout)
		if err != nil {
			return fmt.Errorf("cannot export to cli format:%w", err)
		}
	case "toml":
		err := ExportToml(conf, os.Stdout)
		if err != nil {
			return fmt.Errorf("cannot export to toml format:%w", err)
		}
	case "yml", "yaml":
		err := ExportYaml(conf, os.Stdout)
		if err != nil {
			return fmt.Errorf("cannot export to yaml format:%w", err)
		}
	case "docker":
		err := ExportDocker(conf, "./cmd/docker-compose-tpl.yml", os.Stdout)
		if err != nil {
			return fmt.Errorf("cannot export to docker format:%w", err)
		}
	case "kubernetes", "crd", "k8s":
		err := ExportKubernetes(conf, "./cmd/traefik-lb-svc-tpl.yml", os.Stdout)
		if err != nil {
			return fmt.Errorf("cannot export to kubernetes format:%w", err)
		}
	}

	return nil
}

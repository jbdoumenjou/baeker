package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jbdoumenjou/baeker/cmd"
	"github.com/spf13/cobra"
)

const (
	inDockerCompose          = "In a Docker Compose File"
	asKubernetesLoadBalancer = "As a Kubernetes Load Balancer"
	asTOMLFile               = "As a Toml File"
	asYAMLFile               = "As a Yaml File"
	asCLI                    = "As CLI"
)

func main() {
	rootCmd := createRootCmd()
	rootCmd.AddCommand(createExportCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "baeker",
		Short:   "Baeker - The Traefik Configuration generator.",
		Long:    `Baeker - The Traefik Configuration generator.`,
		Version: "0.0.1",
		RunE: func(_ *cobra.Command, _ []string) error {
			return rootRun()
		},
	}
}

func rootRun() error {
	answers := struct {
		Provider string `survey:"provider"`
	}{}

	// perform the questions
	err := survey.Ask([]*survey.Question{
		{
			Name: "Provider",
			Prompt: &survey.Select{
				Message: "Where do you want to define Traefik?",
				Options: []string{inDockerCompose, asKubernetesLoadBalancer, asTOMLFile, asYAMLFile, asCLI},
				Default: inDockerCompose,
				Help:    "https://doc.traefik.io/traefik/v2.4/providers/overview/#supported-providers",
			},
		},
	}, &answers)
	if err != nil {
		return fmt.Errorf("cannot create survey:%w", err)
	}

	switch answers.Provider {
	case inDockerCompose:
		exportToDockerCompose()
	case asKubernetesLoadBalancer:
		exportToKubernetesCRD()
	case asTOMLFile:
		exportToTomlFile()
	case asYAMLFile:
		exportToYamlFile()
	case asCLI:
		exportToCLI()
	default:
		fmt.Printf("%s not supported", answers.Provider)
	}

	return nil
}

func createExportCmd() *cobra.Command {
	var to string
	cmd := &cobra.Command{
		Use:     "export [file path]",
		Aliases: []string{"e"},
		Short:   "Exports a yml static configuration file to standard output with a specified format.",
		Long:    "Exports a yml static configuration file to standard output with a specified format.",
		RunE: func(_ *cobra.Command, args []string) error {
			confReader, err := os.Open(filepath.FromSlash(args[0]))
			if err != nil {
				return fmt.Errorf("cannot open source file:%w", err)
			}
			err = cmd.ExportCmd(confReader, to)
			if err != nil {
				return fmt.Errorf("cannot export %s to %s: %w", args[0], to, err)
			}

			return nil
		},
		Example: `  $ baeker export traefik.yml
  $ baeker e traefik.yml`,
	}
	cmd.Flags().StringVarP(&to, "to", "t", "cli", "to output format")

	return cmd
}

func exportToDockerCompose() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0o750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Cannot create directory to export conf: %v", err)
			return
		}
	}

	f, err := os.Create("./out/docker-compose.yml")
	if err != nil {
		fmt.Println("cannot create file: ", err)
		return
	}

	builder, err := cmd.NewStaticConfBuilder().AddDockerProvider()
	if err != nil {
		fmt.Printf("cannot add Docker provider to the static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("web", ":8000")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("websecure", ":8443")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	err = cmd.ExportDocker(builder.GetConfiguration(), "./cmd/docker-compose-tpl.yml", f)
	if err != nil {
		fmt.Printf("Failed to export Traefik configuration in %s: %q\n", f.Name(), err.Error())
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s\n", f.Name())
}

func exportToKubernetesCRD() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0o750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Cannot create directory to export conf: %v", err)
			return
		}
	}

	f, err := os.Create("./out/traefik-lb-svc.yml")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	builder, err := cmd.NewStaticConfBuilder().AddKubernetesProvider()
	if err != nil {
		fmt.Printf("cannot create static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("web", ":8000")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("websecure", ":8443")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	err = cmd.ExportDocker(builder.GetConfiguration(), "./cmd/traefik-lb-svc-tpl.yml", f)
	if err != nil {
		fmt.Printf("Failed to export Traefik configuration in %s: %q\n", f.Name(), err.Error())
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s\n", f.Name())
}

func exportToTomlFile() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0o750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Cannot create directory to export conf: %v", err)
			return
		}
	}

	f, err := os.Create("./out/traefik.toml")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	builder, err := cmd.NewStaticConfBuilder().AddFileProvider("conf")
	if err != nil {
		fmt.Printf("cannot AddFileProvider to the static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("web", ":8000")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("websecure", ":8443")
	if err != nil {
		fmt.Printf("cannot AddEntryPoint to the static configuration: %s", err)
		return
	}

	err = cmd.ExportToml(builder.GetConfiguration(), f)
	if err != nil {
		fmt.Printf("Failed to export Traefik configuration in %s: %q\n", f.Name(), err.Error())
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s\n", f.Name())
}

func exportToYamlFile() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0o750)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Cannot create directory to export conf: %v", err)
			return
		}
	}

	f, err := os.Create("./out/traefik.yaml")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	builder, err := cmd.NewStaticConfBuilder().AddFileProvider("conf")
	if err != nil {
		fmt.Printf("cannot create static configuration: %s", err)
	}

	_, err = builder.AddEntryPoint("web", ":8000")
	if err != nil {
		fmt.Printf("cannot create static configuration: %s", err)
	}

	_, err = builder.AddEntryPoint("websecure", ":8443")
	if err != nil {
		fmt.Printf("cannot create static configuration: %s", err)
	}

	err = cmd.ExportYaml(builder.GetConfiguration(), f)
	if err != nil {
		fmt.Printf("Failed to export Traefik configuration in %s: %q\n", f.Name(), err.Error())
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s\n", f.Name())
}

func exportToCLI() {
	builder, err := cmd.NewStaticConfBuilder().AddFileProvider("conf")
	if err != nil {
		fmt.Printf("cannot create static configuration: %s", err)
		return
	}

	_, err = builder.AddEntryPoint("web", ":8000")
	if err != nil {
		fmt.Printf("Failed to add entrypoint: %q\n", err.Error())
		return
	}

	_, err = builder.AddEntryPoint("websecure", ":8443")
	if err != nil {
		fmt.Printf("Failed to add entrypoint: %q\n", err.Error())
		return
	}

	err = cmd.ExportCLI(builder.GetConfiguration(), os.Stdout)
	if err != nil {
		fmt.Printf("Failed to export Traefik configuration: %q\n", err.Error())
		return
	}

	fmt.Println("Successfully exported Traefik configuration")
}

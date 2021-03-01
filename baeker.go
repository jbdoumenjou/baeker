package main

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jbdoumenjou/baeker/cmd"
)

const (
	inDockerCompose          = "In a Docker Compose File"
	asKubernetesLoadBalancer = "As a Kubernetes Load Balancer"
	asTOMLFile               = "As a Toml File"
	asYAMLFile               = "As a Yaml File"
	asCLI                    = "As CLI"
)

func main() {
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
				Help:    "https://doc.traefik.io/traefik/v2.3/providers/overview/#supported-providers",
			},
		},
	}, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
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

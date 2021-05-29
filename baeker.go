package main

import (
	"fmt"
	"os"

	"github.com/jbdoumenjou/baeker/cmd/exporter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jbdoumenjou/baeker/cmd"
)

const (
	inDockerCompose          = "In a Docker Compose File"
	asKubernetesLoadBalancer = "As a Kubernetes Load Balancer"
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
				Options: []string{inDockerCompose, asKubernetesLoadBalancer},
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
	default:
		fmt.Printf("%s not supported", answers.Provider)
	}
}

func exportToDockerCompose() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0750)
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
		fmt.Printf("cannot create static configuration: %s", err)
	}

	builder.AddEntryPoint("web", ":8000")
	builder.AddEntryPoint("websecure", ":8443")
	exporter.ExportDocker(builder.GetConfiguration(), "./cmd/exporter/docker-compose-tpl.yml", f)

	fmt.Printf("Successfully exported Traefik configuration in %s", f.Name())
}

func exportToKubernetesCRD() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0750)
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
	}

	builder.AddEntryPoint("web", ":8000")
	builder.AddEntryPoint("websecure", ":8443")
	exporter.ExportDocker(builder.GetConfiguration(), "./cmd/exporter/traefik-lb-svc-tpl.yml", f)

	fmt.Printf("Successfully exported Traefik configuration in %s", f.Name())
}

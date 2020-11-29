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
)

var qs = []*survey.Question{
	{
		Name: "Provider",
		Prompt: &survey.Select{
			Message: "Where do you want to define Traefik?",
			Options: []string{inDockerCompose, asKubernetesLoadBalancer},
			Default: inDockerCompose,
			Help:    "https://doc.traefik.io/traefik/v2.3/providers/overview/#supported-providers",
		},
	},
}

func main() {
	answers := struct {
		Provider string `survey:"provider"`
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
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
		err := os.Mkdir("out", 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("Cannot create directory to export conf: %v", err)
			return
		}
	}

	f, err := os.Create("./out/docker-compose.yml")
	if err != nil {
		fmt.Println("create file: ", err)
		return
	}

	err = cmd.ExportConf(cmd.GetDefaultConf("docker"), "./cmd/docker-compose-tpl.yml", f)
	if err != nil {
		fmt.Printf("Didn't succeed to export the configuration in docker-compose file: %v", err)
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s", f.Name())
}

func exportToKubernetesCRD() {
	if _, err := os.Stat("output"); os.IsNotExist(err) {
		err := os.Mkdir("out", 0755)
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

	err = cmd.ExportConf(cmd.GetDefaultConf("kubernetes"), "./cmd/traefik-lb-svc-tpl.yml", f)
	if err != nil {
		fmt.Printf("Didn't succeed to export in a kubernetes configuration file: %v", err)
		return
	}

	fmt.Printf("Successfully exported Traefik configuration in %s", f.Name())
}

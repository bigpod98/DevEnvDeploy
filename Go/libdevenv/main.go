package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

func main() {
	fmt.Println("yey")
}

type settings struct {
	commitName  string
	commitEmail string

	templateLocation      string
	localTemplateLocaiton string
}

func BuildEnv(nameofimage string, template string, PAT string, git_repo string, username string, settingsIN string) {
	fmt.Println("empty")

	var envSettings settings

	err := json.Unmarshal([]byte(settingsIN), &envSettings)
	if err != nil {
		log.Fatal(err)
	}
	var fileTemplate []byte

	if envSettings.templateLocation == "local" {
		fmt.Printf("test")

		fileTemplate, err = os.ReadFile(fmt.Sprintf("%s/%s", envSettings.localTemplateLocaiton, template))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(fileTemplate)

	} else {
		fmt.Println("else")
		//TODO: func top access apiserver with templates
	}

	fileTemplateString := string(fileTemplate)

	fileTemplateString = strings.Replace(fileTemplateString, "@PAT", PAT, 1)
	fileTemplateString = strings.Replace(fileTemplateString, "@GITREPO", git_repo, 1)
	fileTemplateString = strings.Replace(fileTemplateString, "@USERNAME", username, 1)
	fileTemplateString = strings.Replace(fileTemplateString, "@GITCONFIGNAME", envSettings.commitName, 1)
	fileTemplateString = strings.Replace(fileTemplateString, "@GITCONFIGEMAIL", envSettings.commitEmail, 1)

	fileTemplate = []byte(fileTemplateString)

	os.WriteFile(fmt.Sprintf("/tmp/%s.dockerfile", nameofimage), fileTemplate, os.ModeTemporary)

	cmd := exec.Command("docker", "build", ".", "-f", fmt.Sprintf("/tmp/%s.dockerfile", nameofimage), "-t", fmt.Sprintf("%s:latest", nameofimage))
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(stdout))
}

func RunEnv(nameofimage string, nameofenv string) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	config := docker.Config{
		Image: nameofimage,
	}

	containeropts := docker.CreateContainerOptions{
		Config: &config,
	}

	container, err := client.CreateContainer(containeropts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.StartContainer(container.ID, container.HostConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func StopEnv(nameofenv string) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	listopts := docker.ListContainersOptions{
		All: true,
	}

	containers, err := client.ListContainers(listopts)
	if err != nil {
		log.Fatal(err)
	}

	var containerID string
	isbroken := true

	for _, element := range containers {
		names := element.Names
		for _, nameElement := range names {
			if nameElement == nameofenv {
				containerID = nameElement
				isbroken = true
				break
			}
			if isbroken {
				break
			}
		}
	}

	client.StopContainer(containerID, 10)
}

func StartEnv(nameofenv string) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	listopts := docker.ListContainersOptions{
		All: true,
	}

	containers, err := client.ListContainers(listopts)
	if err != nil {
		log.Fatal(err)
	}

	var containerID string
	isbroken := true

	for _, element := range containers {
		names := element.Names
		for _, nameElement := range names {
			if nameElement == nameofenv {
				containerID = nameElement
				isbroken = true
				break
			}
			if isbroken {
				break
			}
		}
	}

	inspectopts := docker.InspectContainerOptions{
		ID: containerID,
	}

	container, err := client.InspectContainerWithOptions(inspectopts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.StartContainer(container.ID, container.HostConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveEnv(nameofenv string) {
	fmt.Println("not implemented yet")
	//TODO: do as soon as you can
}
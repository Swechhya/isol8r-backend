package services

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Resource struct {
	Path string `yaml:"path"`
}

type Environment struct {
	Name string   `yaml:"name"`
	Envs []string `yaml:"envs"`
}

type Config struct {
	Namespace          string        `yaml:"namespace"`
	Resources          []string      `yaml:"resources"`
	Patches            []Resource    `yaml:"patches"`
	ConfigMapGenerator []Environment `yaml:"configMapGenerator"`
}

func GenerateManifest(branch string) error {
	dir := fmt.Sprintf("manifest/%s", branch)

	_, err := os.Stat(dir)
	if err == nil {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	err = os.Mkdir(dir, 0777)
	if err != nil {
		return err
	}

	os.Chdir(dir)

	k := Config{
		Namespace: branch,
		Resources: []string{"../../base"},
		Patches: []Resource{
			{Path: "deployment.yml"},
			{Path: "service.yml"},
		},
		ConfigMapGenerator: []Environment{
			{
				Name: branch,
				Envs: []string{".env"},
			},
		},
	}

	data, err := yaml.Marshal(k)
	if err != nil {
		return err
	}

	err = os.WriteFile("kustomization.yml", data, 0666)
	if err != nil {
		return err
	}

	return nil
}

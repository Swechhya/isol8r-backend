package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func GenerateDeployManifest(branch string) error {
	dir := fmt.Sprintf("deploy-manifest/overlay/%s", branch)

	_, err := os.Stat(dir)
	if err == nil {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// os.Chdir(dir)

	kustomize := KustomizationConfig{
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

	deployment := DeploymentConfig{
		APIVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: Metadata{
			Name: "web-deployment",
		},
		Spec: Spec{
			Replicas: 1,
			Selector: Selector{
				MatchLabels: MatchLabels{App: "web"},
			},
			Template: Template{
				Metadata: Metadata{
					Name: "web-deployment",
				},
				Spec: SpecInner{
					Containers: []Container{
						{
							Name:  "feature",
							Image: "nginx:1.14.2",
							Ports: []ContainerPort{
								{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	service := DeploymentConfig{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: Metadata{
			Name: "web-service",
		},
		Spec: Spec{
			Selector: Selector{
				App: "web",
			},
			Ports: []Port{
				{
					Name: "http",
					Port: 80,
				},
			},
		},
	}

	ingress := DeploymentConfig{
		APIVersion: "networking.k8s.io/v1",
		Kind:       "Ingress",
		Metadata: Metadata{
			Name: "app-ingress",
		},
		Spec: Spec{
			Rules: []Rule{
				{
					Host: fmt.Sprintf("%s.demo.prajeshpradhan.com.np", branch),
				},
			},
		},
	}

	ns := DeploymentConfig{
		APIVersion: "v1",
		Kind:       "Namespace",
		Metadata: Metadata{
			Name: branch,
		},
	}

	err = saveYaml(filepath.Join(dir, "kustomization.yml"), kustomize)
	if err != nil {
		return err
	}
	err = saveYaml(filepath.Join(dir, "deployment.yml"), deployment)
	if err != nil {
		return err
	}
	err = saveYaml(filepath.Join(dir, "service.yml"), service)
	if err != nil {
		return err
	}
	err = saveYaml(filepath.Join(dir, "ingress.yml"), ingress)
	if err != nil {
		return err
	}
	err = saveYaml(filepath.Join(dir, "ns.yml"), ns)
	if err != nil {
		return err
	}

	return nil
}

func GenerateBuildManifest(branch string) error {
	dir := fmt.Sprintf("./build-manifest/overlay/%s", branch)

	_, err := os.Stat(dir)
	if err == nil {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// os.Chdir(dir)

	// gitPass := GetConfig("")
	gitPass, _ := GetInstallationToken(context.Background())
	dest := "dest-tag"
	repoUrl := "git://github.com/Swechhya/ci-cd-demo.git"

	kustomize := KustomizationConfig{
		Namespace: branch,
		Resources: []string{
			"../../base",
		},
		Patches: []Resource{
			{Path: "pod.yml"},
		},
		SecretGenerator: []Secret{
			{
				Name: "git-token",
				Literals: []string{
					fmt.Sprintf("GIT_PASSWORD=%s", gitPass),
				},
			},
		},
	}

	pod := DeploymentConfig{
		APIVersion: "v1",
		Kind:       "Pod",
		Metadata: Metadata{
			Name: "kaniko",
		},
		Spec: Spec{
			Containers: []Container{
				{
					Name:  "kaniko",
					Image: "gcr.io/kaniko-project/executor:latest",
					Args: []string{
						"--dockerfile=Dockerfile",
						fmt.Sprintf("--context=%s", repoUrl),
						fmt.Sprintf("--destination=%s", dest),
					},
					Env: []Environment{
						{
							Name:  "GIT_USERNAME",
							Value: "x-access-token",
						},
						{
							Name: "GIT_PASSWORD",
							ValueFrom: ValueFrom{
								SecretKeyRef: map[string]string{
									"name": "git-token",
									"key":  "GIT_PASSWORD",
								},
							},
						},
					},
				},
			},
		},
	}

	ns := DeploymentConfig{
		APIVersion: "v1",
		Kind:       "Namespace",
		Metadata: Metadata{
			Name: branch,
		},
	}

	err = saveYaml(filepath.Join(dir, "kustomization.yml"), kustomize)
	if err != nil {
		return err
	}
	err = saveYaml(filepath.Join(dir, "pod.yml"), pod)
	if err != nil {
		return err
	}

	return nil
}

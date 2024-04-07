package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Swechhya/isol8r-backend/services/s3"
)

func GenerateDeployManifest(namespace string, dest string, envUri string) error {
	dir := fmt.Sprintf("deploy-manifest/overlay/%s", namespace)

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
		Namespace: namespace,
		Resources: []string{"../../base", "ns.yml"},
		Patches: []Resource{
			{Path: "deployment.yml"},
			{Path: "service.yml"},
			{Path: "ingress.yml"},
		},
		ConfigMapGenerator: []Environment{
			{
				Name: namespace,
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
							Image: dest,
							Ports: []ContainerPort{
								{
									ContainerPort: 3000,
									Name:          "http",
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
					Name:       "http",
					Port:       3000,
					TargetPort: "http",
					Protocol:   "TCP",
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
					Host: fmt.Sprintf("%s.demo.prajeshpradhan.com.np", namespace),
					Http: Http{
						Paths: []Path{
							{
								Path:     "/",
								PathType: "Prefix",
								Backend: Backend{
									Service: BackendService{
										ServiceName: "ssl-redirect",
										Port:        Port{Name: "use-annotation"},
									},
								},
							},
							{
								Path:     "/",
								PathType: "Prefix",
								Backend: Backend{
									Service: BackendService{
										ServiceName: "web-service",
										Port:        Port{Name: "http"},
									},
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
			Name: namespace,
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

	bucketName := os.Getenv("APP_ENV_BUCKET")
	fullFilePath := filepath.Join(dir, ".env")
	client := s3.GetClient()
	err = client.DownloadFileToPath(context.Background(), bucketName, envUri, fullFilePath)
	if err != nil {
		return err
	}

	return nil
}

func GenerateBuildManifest(namespace string, repoFullName string, branch string, dest string) error {
	dir := fmt.Sprintf("./build-manifest/overlay/%s", namespace)

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

	gitPass, err := GetInstallationToken(context.Background())
	if err != nil {
		return err
	}

	// ecr := GetConfig("dest")
	// if ecr == "" {
	// 	return fmt.Errorf("empty ecr")
	// }

	// ecr := "654654451390.dkr.ecr.us-east-1.amazonaws.com/test:"

	// repoName := strings.Split(repoFullName, "/")[1]
	// dest := fmt.Sprintf("%s%s-%s", ecr, namespace, repoName)
	repoUrl := fmt.Sprintf("git://github.com/%s.git#refs/heads/%s", repoFullName, branch)

	kustomize := KustomizationConfig{
		Namespace: namespace,
		Resources: []string{
			"../../base",
			"ns.yml",
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
			Name: namespace,
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
	err = saveYaml(filepath.Join(dir, "ns.yml"), ns)
	if err != nil {
		return err
	}

	return nil
}

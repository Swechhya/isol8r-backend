package services

type Resource struct {
	Path string `yaml:"path,omitempty"`
}

type Environment struct {
	Name      string    `yaml:"name,omitempty"`
	Value     string    `yaml:"value,omitempty"`
	ValueFrom ValueFrom `yaml:"valueFrom,omitempty"`
	Envs      []string  `yaml:"envs,omitempty"`
}

type ValueFrom struct {
	SecretKeyRef map[string]string `yaml:"secretKeyRef,omitempty"`
}

type KustomizationConfig struct {
	Namespace          string        `yaml:"namespace,omitempty"`
	Resources          []string      `yaml:"resources,omitempty"`
	Patches            []Resource    `yaml:"patches,omitempty"`
	ConfigMapGenerator []Environment `yaml:"configMapGenerator,omitempty"`
	SecretGenerator    []Secret      `yaml:"secretGenerator,omitempty"`
}

type Metadata struct {
	Name string `yaml:"name,omitempty"`
}

type Selector struct {
	App         string      `yaml:"app,omitempty"`
	MatchLabels MatchLabels `yaml:"matchLabels,omitempty"`
}

type MatchLabels struct {
	App string `yaml:"app,omitempty"`
}

type ContainerPort struct {
	ContainerPort int    `yaml:"containerPort,omitempty"`
	Name          string `yaml:"name,omitempty"`
}

type Container struct {
	Name  string          `yaml:"name,omitempty"`
	Image string          `yaml:"image,omitempty"`
	Ports []ContainerPort `yaml:"ports,omitempty"`

	Args []string      `yaml:"args,omitempty"`
	Env  []Environment `yaml:"env,omitempty"`
}

type Template struct {
	Metadata Metadata  `yaml:"metadata,omitempty"`
	Spec     SpecInner `yaml:"spec,omitempty"`
}

type SpecInner struct {
	Containers []Container `yaml:"containers,omitempty"`
}

type DeploymentConfig struct {
	APIVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata,omitempty"`
	Spec       Spec     `yaml:"spec,omitempty"`
}

//

type Port struct {
	Name       string `yaml:"name,omitempty"`
	Port       int    `yaml:"port,omitempty"`
	TargetPort string `yaml:"targetPort,omitempty"`
	Protocol   string `yaml:"protocol,omitempty"`
}

type Service struct {
	APIVersion string   `yaml:"apiVersion,omitempty"`
	Kind       string   `yaml:"kind,omitempty"`
	Metadata   Metadata `yaml:"metadata,omitempty"`
	Spec       Spec     `yaml:"spec,omitempty"`
}

type Spec struct {
	Replicas   int         `yaml:"replicas,omitempty"`
	Selector   Selector    `yaml:"selector,omitempty"`
	Ports      []Port      `yaml:"ports,omitempty"`
	Template   Template    `yaml:"template,omitempty"`
	Containers []Container `yaml:"containers,omitempty"`
	Rules      []Rule      `yaml:"rules,omitempty"`
}

type Rule struct {
	Host string `yaml:"host,omitempty"`
	Http Http   `yaml:"http,omitempty"`
}

type Http struct {
	Paths []Path `yaml:"paths,omitempty"`
}

type Secret struct {
	Name     string   `yaml:"name,omitempty"`
	Literals []string `yaml:"literals,omitempty"`
}

type Backend struct {
	Service BackendService `ymal:"service,omitempty"`
}

type BackendService struct {
	ServiceName string `yaml:"name,omitempty"`
	Port        Port   `yaml:"port,omitempty"`
}

type Path struct {
	Path     string  `yaml:"path,omitempty"`
	PathType string  `yaml:"pathType,omitempty"`
	Backend  Backend `yaml:"backend,omitempty"`
}

// type HTTPRule struct {
// 	Host   string `yaml:"host,omitempty"`
// 	Paths  []Path `yaml:"paths,omitempty"`
// }

// type Ingress struct {
// 	Rules []HTTPRule `yaml:"rules,omitempty"`
// }

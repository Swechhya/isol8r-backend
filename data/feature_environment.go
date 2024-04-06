package data

import "time"

type FeatureEnvironment struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Identifier  string     `json:"identifier"`
	Description string     `json:"description"`
	DBType      string     `json:"dbType"`
	CreatedBy   string     `json:"createdBy"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	Resources   []Resource `json:"resources"`
}

type Resource struct {
	FeatureEnvID int    `json:"featureEnvironmentId"`
	RepoID       int    `json:"repoId"`
	IsAutoUpdate bool   `json:"isAutoUpdate"`
	Branch       string `json:"branch"`
	Link         string `json:"link"`
}

type Commit struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

// RepositoryData represents data about a GitHub repository
type Branch struct {
	Name      string `json:"name"`
	Commit    Commit `json:"commit"`
	Protected bool   `json:"protected"`
}

type ReadyRepositories struct {
	ID        uint64    `json:"id"`
	RepoID    int       `json:"repo_id"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	UserLogin string    `json:"user_login"`
	URL       string    `json:"url"`
	Setup     bool      `json:"setup"`
	EnvURI    string    `json:"env_uri"`
	Branch    []*Branch `json:"branches"`
}

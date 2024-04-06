package data

import "time"

type FeatureEnvironment struct {
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

type RepoList struct {
	Repositories []*Repo `json:"repositories"`
}

type Repo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type InstallationToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

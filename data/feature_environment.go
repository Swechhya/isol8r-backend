package data

import "time"

type FeatureEnvironment struct {
	Name      string     `json:"name"`
	FeatureID string     `json:"feature_id"`
	DBType    string     `json:"dbType"`
	CreatedBy string     `json:"createdBy"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Resources []Resource `json:"resources"`
}

type Resource struct {
	FeatureEnvID int    `json:"feature_environment_id"`
	AppName      string `json:"appName"`
	IsAutoUpdate bool   `json:"isAutoUpdate"`
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

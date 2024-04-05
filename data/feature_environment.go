package data

type FeatureEnvironment struct {
	Name      string     `json:"name" bson:"name"`
	DBType    string     `json:"dbType" bson:"dbType"`
	CreatedAt string     `json:"createdAt" bson:"createdAt"`
	CreatedBy string     `json:"createdBy" bson:"createdBy"`
	Resources []Resource `json:"resources" bson:"resources"`
}

type Resource struct {
	AppName      string `json:"appName" bson:"appName"`
	IsAutoUpdate bool   `json:"isAutoUpdate" bson:"isAutoUpdate"`
}

type RepoList struct {
	Repositories []*Repo `json:"repositories"`
}

type Repo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

package data

type Config struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedBy string `json:"createdBy"`
}

type GithubClientSetup struct {
	AppID  string `json:"app_id"`
	InstID string `json:"inst_id"`
}

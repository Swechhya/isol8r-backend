package data

type Config struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedBy string `json:"createdBy"`
}

type GithubClientSetup struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	InstallID    string `json:"installID"`
	PrivateKey   string `json:"privateKey"`
}

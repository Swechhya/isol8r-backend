package data

type InstallationToken struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type Repositories struct {
	List []*Repo `json:"repositories"`
}

type Repo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Url      string `json:"url"`
}

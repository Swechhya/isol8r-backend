package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/jwt"
	"github.com/beatlabs/github-auth/key"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Client     *github.Client
	HttpClient *http.Client
}

var Gh *GitHubClient
var User *github.User

func SetupGithubClient(config *data.GithubClientSetup) error {
	key, err := key.FromFile("./key.pem")
	if err != nil {
		return err
	}

	installID := config.InstallID
	privateKey := config.PrivateKey
	if installID == "" || privateKey == "" {
		return errors.New("installID or privateKey missing in configs")
	}

	install, err := inst.NewConfig(installID, privateKey, key)
	if err != nil {
		return err
	}

	err = saveClientConfigToDB(installID, privateKey)
	if err != nil {
		return err
	}

	ctx := context.Background()
	httpClient := install.Client(ctx)
	client := github.NewClient(httpClient)

	Gh = &GitHubClient{
		HttpClient: httpClient,
		Client:     client,
	}

	return nil
}

func GetRepos(ctx context.Context) ([]*data.Repo, error) {
	r, err := Gh.HttpClient.Get("https://api.github.com/installation/repositories")
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	repos := new(data.RepoList)

	err = json.Unmarshal(b, repos)
	if err != nil {
		return nil, err
	}

	//insert into table
	db := db.DB()
	for _, repo := range repos.Repositories {
		_, err := db.Exec(`
        INSERT INTO repositories (name, full_name, user_login, s3_uri, created_by)
        VALUES ($1, $2, $3, $4, $5)
    `, repo.Name, repo.Name, "", "", "")
		if err != nil {
			return nil, err
		}
	}

	return repos.Repositories, nil
}

func GetBranches(ctx context.Context, repo string) ([]*github.Branch, error) {
	branches, _, err := Gh.Client.Repositories.ListBranches(ctx, *User.Login, repo, nil)
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func GetInstallationToken(ctx context.Context) (string, error) {

	key, err := key.FromFile("./key.pem")
	if err != nil {
		return "", err
	}
	jt := jwt.JWT{AppID: "app_id", PrivateKey: key, Expires: time.Minute * 10}
	je, err := jt.Payload()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", "install_id"), nil)
	if err != nil {
		return "", nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", je))
	req.Header.Add("Accept", "application/vnd.github.v3+json")

	r, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if r.StatusCode < 200 && r.StatusCode >= 400 {
		return "", fmt.Errorf("error executing request")
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	token := new(data.InstallationToken)

	err = json.Unmarshal(b, token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

func saveClientConfigToDB(installID, privateKey string) error {
	db := db.DB()

	dq := goqu.Insert("core_config").
		Cols("key", "value").
		Vals(goqu.Vals{"installID", installID},
			goqu.Vals{"privateKey", privateKey},
		)

	insertSql, args, err := dq.ToSQL()
	if err != nil {
		return err
	}
	_, err = db.Exec(insertSql, args...)

	if err != nil {
		return err
	}
	return nil
}

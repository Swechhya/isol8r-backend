package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Swechhya/panik-backend/data"
	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/jwt"
	"github.com/beatlabs/github-auth/key"
	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Client     *github.Client
	HttpClient *http.Client
}

var Gh *GitHubClient
var User *github.User

func SetupGithubClient() error {
	key, err := key.FromFile("./key.pem")
	if err != nil {
		return err
	}

	install, err := inst.NewConfig("app_id", "install_id", key)
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

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		return err
	}
	User = user
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

	return repos.Repositories, nil
}

func GetBranches(ctx context.Context, branch string) ([]*github.Branch, error) {
	branches, _, err := Gh.Client.Repositories.ListBranches(ctx, *User.Login, branch, nil)
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

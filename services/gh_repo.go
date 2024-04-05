package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Swechhya/panik-backend/data"
	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/key"
	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Client     *github.Client
	HttpClient *http.Client
}

var Gh *GitHubClient

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
	return nil
}

func GetRepos(ctx context.Context) ([]*data.Repo, error) {
	err := SetupGithubClient()
	if err != nil {
		return nil, err
	}

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

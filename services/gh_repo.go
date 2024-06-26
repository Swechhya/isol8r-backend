package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/internal/db"
	"github.com/Swechhya/isol8r-backend/services/s3"
	"github.com/beatlabs/github-auth/app/inst"
	"github.com/beatlabs/github-auth/jwt"
	"github.com/beatlabs/github-auth/key"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

type GitHubClient struct {
	Client     *github.Client
	HttpClient *http.Client
}

var Gh *GitHubClient

func SetupGithubClient(ctx context.Context, config *data.GithubClientSetup) error {
	key, err := key.FromFile(config.PrivateKeyPath)
	if err != nil {
		return err
	}

	installID := config.InstallID
	appID := config.AppID
	privateKeyPath := config.PrivateKeyPath
	if installID == "" || privateKeyPath == "" || appID == "" {
		return errors.New("installID or appID missing in configs")
	}

	install, err := inst.NewConfig(installID, appID, key)
	if err != nil {
		return err
	}

	err = AddConfig("installID", installID)
	if err != nil {
		return err
	}
	err = AddConfig("privateKeyPath", privateKeyPath)
	if err != nil {
		return err
	}

	err = AddConfig("appID", appID)
	if err != nil {
		return err
	}

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

	repos := new(data.Repositories)
	err = json.Unmarshal(b, repos)
	if err != nil {
		return nil, err
	}

	db := db.DB()
	// Delete existing records from the repositories table
	deleteExpr := goqu.Delete("repositories")
	sql, args, err := deleteExpr.ToSQL()
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(sql, args...)
	if err != nil {
		return nil, err
	}

	// Insert new records into the repositories table
	for _, repo := range repos.List {
		insertExpr := goqu.Insert("repositories").Rows(
			goqu.Record{
				"repo_id":    repo.Id,
				"name":       repo.Name,
				"full_name":  repo.FullName,
				"user_login": "",
				"url":        repo.Url,
			},
		)
		sql, args, err := insertExpr.ToSQL()
		if err != nil {
			return nil, err
		}
		_, err = db.Exec(sql, args...)
		if err != nil {
			return nil, err
		}
	}

	return repos.List, nil
}

func GetBranches(ctx context.Context, repoId int64) ([]*github.Branch, error) {
	repo, _, err := Gh.Client.Repositories.GetByID(ctx, repoId)
	if err != nil {
		return nil, err
	}

	repoName := *repo.Name
	owner := *repo.Owner.Login

	branches, _, err := Gh.Client.Repositories.ListBranches(ctx, owner, repoName, nil)
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func GetInstallationToken(ctx context.Context) (string, error) {
	// key, err := key.FromFile(GetConfig("privateKeyPath"))

	key, err := key.FromFile("key.pem")
	if err != nil {
		return "", err
	}
	// jt := jwt.JWT{AppID: GetConfig("appID"), PrivateKey: key, Expires: time.Minute * 10}

	jt := jwt.JWT{AppID: "870502", PrivateKey: key, Expires: time.Minute * 10}
	je, err := jt.Payload()
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	// req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", GetConfig("installID")), nil)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", "49284279"), nil)
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

func UploadEnvFile(c *gin.Context, file io.Reader, repoId, fileName string) (string, error) {
	bucketName := os.Getenv("APP_ENV_BUCKET")

	client := s3.GetClient()
	uri, err := client.UploadFile(c, bucketName, fileName, file)
	if err != nil {
		return "", err
	}

	// Update repository in the database
	if err := updateRepository(repoId, fileName); err != nil {
		return "", err
	}

	return *uri, nil
}

func updateRepository(repoId, uri string) error {
	db := db.DB()

	// Update the repositories table
	updateExpr, args, err := goqu.Update("repositories").
		Set(goqu.Record{"env_uri": uri, "setup": true}).
		Where(goqu.Ex{
			"repo_id": goqu.Op{"eq": repoId},
		}).ToSQL()

	if err != nil {
		return err
	}

	// Execute the SQL query
	_, err = db.Exec(updateExpr, args...)
	if err != nil {
		return err
	}

	return nil
}

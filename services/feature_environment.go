package services

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime/debug"
	"strings"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
)

func GetAllFeatureEnvironments() ([]*data.FeatureEnvironment, error) {
	query := goqu.
		From(goqu.T("feature_environments")).
		Select(
			"feature_environments.id",
			"feature_environments.name",
			"feature_environments.identifier",
			"feature_environments.created_at",
			"feature_environments.created_by",
			"resources.feature_environment_id",
			"resources.repo_id",
			"resources.is_auto_update",
			"resources.branch",
			"resources.link",
			"resources.port",
		).
		LeftJoin(goqu.T("resources"), goqu.On(goqu.I("feature_environments.id").Eq(goqu.I("resources.feature_environment_id"))))

	fmt.Print(query.ToSQL())
	selectSQL, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.DB().Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	envMap := make(map[int]*data.FeatureEnvironment)
	for rows.Next() {
		var fe data.FeatureEnvironment
		var res data.Resource
		if err := rows.Scan(
			&fe.ID, &fe.Name, &fe.Identifier, &fe.CreatedAt, &fe.CreatedBy,
			&res.FeatureEnvID, &res.RepoID, &res.IsAutoUpdate, &res.Branch, &res.Link, &res.Port,
		); err != nil {
			return nil, err
		}

		if env, ok := envMap[fe.ID]; ok {
			// If feature environment already exists in the map, add the resource to its resources slice
			env.Resources = append(env.Resources, res)
		} else {
			// If feature environment doesn't exist in the map, create a new entry and add it to the map
			fe.Resources = append(fe.Resources, res)
			envMap[fe.ID] = &fe
		}
	}

	var envLists []*data.FeatureEnvironment
	for _, env := range envMap {
		envLists = append(envLists, env)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return envLists, nil
}

func GetFeatureEnvironmentById(id int) (*data.FeatureEnvironment, error) {
	query := goqu.
		From(goqu.T("feature_environments")).
		Select(
			"feature_environments.id",
			"feature_environments.name",
			"feature_environments.identifier",
			"feature_environments.created_at",
			"feature_environments.created_by",
			"resources.feature_environment_id",
			"resources.repo_id",
			"resources.is_auto_update",
			"resources.branch",
			"resources.link",
			"resources.port",
		).
		LeftJoin(goqu.T("resources"), goqu.On(goqu.I("feature_environments.id").Eq(goqu.I("resources.feature_environment_id")))).
		Where(goqu.Ex{"feature_environments.id": id})

	selectSQL, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.DB().Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fe *data.FeatureEnvironment
	resources := make([]data.Resource, 0)

	for rows.Next() {
		var tempFE data.FeatureEnvironment
		var tempRes data.Resource
		if err := rows.Scan(
			&tempFE.ID, &tempFE.Name, &tempFE.Identifier, &tempFE.CreatedAt, &tempFE.CreatedBy,
			&tempRes.FeatureEnvID, &tempRes.RepoID, &tempRes.IsAutoUpdate, &tempRes.Branch, &tempRes.Link, &tempRes.Port,
		); err != nil {
			return nil, err
		}

		if fe == nil {
			fe = &tempFE
		}

		resources = append(resources, tempRes)
	}

	if fe != nil {
		fe.Resources = resources
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fe, nil
}

func CreateFeatureEnvironment(fe data.FeatureEnvironment, reDeploy bool) (int, error) {
	if reDeploy {
		if err := DeleteFeatureEnvironment(fe.ID); err != nil {
			return 0, err
		}
	}

	if fe.Identifier == "" {
		return 0, fmt.Errorf("empty feature identifier")
	}
	// Insert into table
	db := db.DB()
	dq := goqu.Insert("feature_environments").
		Cols("name", "identifier", "db_type", "created_by").
		Vals(goqu.Vals{fe.Name, fe.Identifier, fe.DBType, fe.CreatedBy})

	insertSql, args, err := dq.ToSQL()
	if err != nil {
		return 0, err
	}

	_, err = db.Exec(insertSql, args...)
	if err != nil {
		// Check if the error is due to duplicate identifier
		if strings.Contains(err.Error(), "feature_environments_identifier_key") {
			return 0, errors.New("duplicate Identifier.")
		}
		return 0, err
	}

	var feID int

	err = db.QueryRow("SELECT lastval()").Scan(&feID)
	if err != nil {
		return 0, err
	}

	// ecr := GetConfig("dest")
	// if ecr == "" {
	// 	return fmt.Errorf("empty ecr")
	// }

	ecr := "654654451390.dkr.ecr.us-east-1.amazonaws.com/test:"

	// Iterate over resources and insert them into the database
	for _, resource := range fe.Resources {
		dq := goqu.From("repositories").Select("full_name", "env_uri")

		sql, args, err := dq.ToSQL()
		for err != nil {
			return 0, err
		}
		rows, err := db.Query(sql, args...)
		for err != nil {
			return 0, err
		}

		defer rows.Close()

		var repoFullName, envUri string
		for rows.Next() {
			if err := rows.Scan(&repoFullName, &envUri); err != nil {
				return 0, err
			}
		}

		repoName := strings.Split(repoFullName, "/")[1]
		dest := fmt.Sprintf("%s%s-%s", ecr, fe.Identifier, repoName)

		err = GenerateBuildManifest(fe.Identifier, repoFullName, resource.Branch, dest)
		if err != nil {
			return 0, err
		}
		err = GenerateDeployManifest(fe.Identifier, dest, envUri)
		if err != nil {
			return 0, err
		}
		err = DeployEnvironment(fe.Identifier)
		if err != nil {
			return 0, err
		}

		resource.FeatureEnvID = feID
		if err := insertResource(resource); err != nil {
			return 0, err
		}
	}

	return feID, err
}

func DeleteFeatureEnvironment(feID int) error {
	db := db.DB()

	q := goqu.From("feature_environments").Select("identifier")
	sql, args, err := q.ToSQL()
	for err != nil {
		return err
	}
	rows, err := db.Query(sql, args...)
	for err != nil {
		return err
	}

	defer rows.Close()
	var identifier string
	for rows.Next() {
		if err := rows.Scan(&identifier); err != nil {
			return err
		}
	}

	err = runDeleteKCommand(identifier)
	if err != nil {
		fmt.Println("Err: "+err.Error(), string(debug.Stack()))
		return err
	}

	deleteResourcesExpr := goqu.Delete("resources").Where(goqu.I("feature_environment_id").Eq(feID))
	deleteResourcesSQL, args, err := deleteResourcesExpr.ToSQL()
	if err != nil {
		return err
	}
	_, err = db.Exec(deleteResourcesSQL, args...)
	if err != nil {
		return err
	}

	deleteExpr := goqu.Delete("feature_environments").Where(goqu.I("id").Eq(feID))
	deleteSQL, args, err := deleteExpr.ToSQL()
	if err != nil {
		return err
	}
	_, err = db.Exec(deleteSQL, args...)
	if err != nil {
		return err
	}

	return nil
}

func insertResource(resource data.Resource) error {
	// Insert Resource into database
	db := db.DB()

	dq := goqu.Insert("resources").
		Cols("feature_environment_id", "repo_id", "branch", "is_auto_update", "link", "port").
		Vals(goqu.Vals{resource.FeatureEnvID, resource.RepoID, resource.Branch, resource.IsAutoUpdate, resource.Link, resource.Port})

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

func DeployEnvironment(namespace string) error {
	buildManifestPath := fmt.Sprintf("build-manifest/overlay/%s", namespace)

	// kubectl kustomize ./overlay/feature | kubectl apply -f -
	err := runStartKCommand(buildManifestPath)
	if err != nil {

		fmt.Println("Err: "+err.Error(), string(debug.Stack()))
		return err
	}

	deployManifestPath := fmt.Sprintf("deploy-manifest/overlay/%s", namespace)
	// kubectl kustomize ./overlay/feature | kubectl apply -f -
	err = runStartKCommand(deployManifestPath)
	if err != nil {

		fmt.Println("Err: "+err.Error(), string(debug.Stack()))
		return err
	}

	// run kustomize and kubectl
	return nil
}

func runStartKCommand(path string) error {
	kustomizeCmd := exec.Command("kubectl", "kustomize", path)
	applyCmd := exec.Command("kubectl", "apply", "-f", "-")

	output, err := kustomizeCmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe for kustomize command: %v\n", err)
		return err
	}
	applyCmd.Stdin = output

	if err := kustomizeCmd.Start(); err != nil {
		fmt.Printf("Error starting kustomize command: %v\n", err)
		return err
	}
	// if err := applyCmd.Start(); err != nil {
	// 	return err
	// }
	o, err := applyCmd.Output()
	if err != nil {
		fmt.Println(o)
		return err
	}
	fmt.Println(o)

	return nil
}

func runDeleteKCommand(namespace string) error {
	kustomizeCmd := exec.Command("kubectl", "delete", namespace)

	if err := kustomizeCmd.Start(); err != nil {
		fmt.Printf("Error deleting kustomize command: %v\n", err)
		return err
	}
	return nil
}

func FetchLaunchReadyRepos(c *gin.Context) ([]*data.ReadyRepositories, error) {
	db := db.DB()

	query := goqu.From("repositories").Select("id", "repo_id", "name", "full_name", "url", "setup", "env_uri").
		Where(goqu.Ex{"setup": true})
	selectSQL, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repoLists []*data.ReadyRepositories

	for rows.Next() {
		var repo data.ReadyRepositories
		if err := rows.Scan(&repo.ID, &repo.RepoID, &repo.Name, &repo.FullName, &repo.URL, &repo.Setup, &repo.EnvURI); err != nil {
			return nil, err
		}

		branches, err := GetBranches(c.Request.Context(), int64(repo.RepoID))
		if err != nil {
			return nil, err
		}

		repo.Branch = ConvertToDataBranches(branches)
		repoLists = append(repoLists, &repo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If successful, return the list of launch ready repos
	return repoLists, err
}

func ConvertToDataBranches(ghBranches []*github.Branch) []*data.Branch {
	var dataBranches []*data.Branch
	for _, b := range ghBranches {
		dataBranches = append(dataBranches, &data.Branch{
			Name:      *b.Name,
			Commit:    data.Commit{SHA: *b.Commit.SHA, URL: *b.Commit.URL},
			Protected: *b.Protected,
		})
	}
	return dataBranches
}

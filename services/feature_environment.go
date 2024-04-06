package services

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
)

func GetAllFeatureEnvironments() ([]*data.FeatureEnvironment, error) {
	query := goqu.From("feature_environments").Select("name", "identifier", "description", "db_type", "created_at", "created_by")
	fmt.Println(query.ToSQL())
	selectSQL, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := db.DB().Query(selectSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var envLists []*data.FeatureEnvironment

	for rows.Next() {
		var fe data.FeatureEnvironment
		if err := rows.Scan(&fe.Name, &fe.Identifier, &fe.Description, &fe.DBType, &fe.CreatedAt, &fe.CreatedBy); err != nil {
			return nil, err
		}

		envLists = append(envLists, &fe)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if envLists == nil {
		envLists = []*data.FeatureEnvironment{
			{
				Name:      "Feature Environment 1",
				DBType:    "MySQL",
				CreatedBy: "John Doe",
				CreatedAt: time.Date(2022, 4, 7, 10, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2022, 4, 7, 10, 30, 0, 0, time.UTC),
				Resources: []data.Resource{
					{
						FeatureEnvID: 1,
						IsAutoUpdate: true,
						Link:         "https://example.com/app1",
					},
					{
						FeatureEnvID: 1,
						IsAutoUpdate: false,
						Link:         "https://example.com/app2",
					},
				},
			},
			{
				Name:      "Feature Environment 2",
				DBType:    "PostgreSQL",
				CreatedBy: "Jane Smith",
				CreatedAt: time.Date(2022, 4, 8, 9, 30, 0, 0, time.UTC),
				UpdatedAt: time.Date(2022, 4, 8, 10, 15, 0, 0, time.UTC),
				Resources: []data.Resource{
					{
						FeatureEnvID: 2,
						IsAutoUpdate: true,
						Link:         "https://example.com/app3",
					},
					{
						FeatureEnvID: 2,
						IsAutoUpdate: true,
						Link:         "https://example.com/app4",
					},
				},
			},
		}
	}

	return envLists, nil
}

func CreateFeatureEnvironment(fe data.FeatureEnvironment) error {
	if fe.Identifier == "" {
		return fmt.Errorf("empty feature identifier")
	}
	// Insert into table
	db := db.DB()
	dq := goqu.Insert("feature_environments").
		Cols("name", "identifier", "db_type", "created_by").
		Vals(goqu.Vals{fe.Name, fe.Identifier, fe.DBType, fe.CreatedBy})

	insertSql, args, err := dq.ToSQL()
	if err != nil {
		return err
	}

	_, err = db.Exec(insertSql, args...)
	if err != nil {
		return err
	}

	var feID int

	err = db.QueryRow("SELECT lastval()").Scan(&feID)
	if err != nil {
		return err
	}

	// ecr := GetConfig("dest")
	// if ecr == "" {
	// 	return fmt.Errorf("empty ecr")
	// }

	ecr := "654654451390.dkr.ecr.us-east-1.amazonaws.com/test:"

	// Iterate over resources and insert them into the database
	for _, resource := range fe.Resources {
		dq := goqu.From("repositories").Select("full_name")

		sql, args, err := dq.ToSQL()
		for err != nil {
			return err
		}
		rows, err := db.Query(sql, args...)
		for err != nil {
			return err
		}
		var repoFullName string
		for rows.Next() {
			if err := rows.Scan(&repoFullName); err != nil {
				return err
			}
		}

		repoName := strings.Split(repoFullName, "/")[1]
		dest := fmt.Sprintf("%s%s-%s", ecr, fe.Identifier, repoName)

		err = GenerateBuildManifest(fe.Identifier, repoName, resource.Branch, dest)
		if err != nil {
			return err
		}
		err = GenerateDeployManifest(fe.Identifier, dest)
		if err != nil {
			return err
		}
		err = DeployEnvironment(fe.Identifier)
		if err != nil {
			return err
		}

		resource.FeatureEnvID = feID
		if err := insertResource(resource); err != nil {
			return err
		}
	}

	return nil
}

func DeleteFeatureEnvironment(identifier string) error {
	err := runDeleteKCommand(identifier)
	if err != nil {
		return err
	}

	db := db.DB()
	deleteExpr := goqu.Delete("feature_environments").Where(goqu.I("identifier").Eq(identifier))
	sql, args, err := deleteExpr.ToSQL()
	if err != nil {
		return err
	}
	_, err = db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func insertResource(resource data.Resource) error {
	// Insert Resource into database
	db := db.DB()

	dq := goqu.Insert("resources").
		Cols("feature_environment_id", "repo_id", "branch", "is_auto_update", "link").
		Vals(goqu.Vals{resource.FeatureEnvID, resource.RepoID, resource.Branch, resource.IsAutoUpdate, resource.Link})

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
		return err
	}

	deployManifestPath := fmt.Sprintf("deploy-manifest/overlay/%s", namespace)
	// kubectl kustomize ./overlay/feature | kubectl apply -f -
	err = runStartKCommand(deployManifestPath)
	if err != nil {
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
	if err := applyCmd.Start(); err != nil {
		return err
	}
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

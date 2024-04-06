package services

import (
	"errors"
	"fmt"
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

	//TODO :: REMOVE LATER
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

func GetFeatureEnvironmentById(id int) (*data.FeatureEnvironment, error) {
	query := goqu.From("feature_environments").
		Select("name", "identifier", "description", "db_type", "created_at", "created_by").
		Where(goqu.Ex{"id": id})

	selectSQL, _, err := query.ToSQL()
	if err != nil {
		return nil, err
	}

	fe := new(data.FeatureEnvironment)
	err = db.DB().QueryRow(selectSQL).Scan(&fe.Name, &fe.Identifier, &fe.Description, &fe.DBType, &fe.CreatedAt, &fe.CreatedBy)
	if err != nil {
		return nil, err
	}

	return fe, nil

}

func CreateFeatureEnvironment(fe data.FeatureEnvironment) error {
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
		// Check if the error is due to duplicate identifier
		if strings.Contains(err.Error(), "feature_environments_identifier_key") {
			return errors.New("duplicate Identifier.")
		}
		return err
	}

	var feID int

	err = db.QueryRow("SELECT lastval()").Scan(&feID)
	if err != nil {
		return err
	}

	// Iterate over resources and insert them into the database
	for _, resource := range fe.Resources {
		resource.FeatureEnvID = feID
		if err := insertResource(resource); err != nil {
			return err
		}
	}

	return nil
}

func DeleteFeatureEnvironment(feID int) error {
	db := db.DB()
	deleteExpr := goqu.Delete("feature_environments").Where(goqu.I("id").Eq(feID))
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

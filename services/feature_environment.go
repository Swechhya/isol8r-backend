package services

import (
	"fmt"

	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
)

func GetAllFeatureEnvironments() ([]*data.FeatureEnvironment, error) {
	query := goqu.From("feature_environments").Select("name", "feature_id", "db_type", "created_at", "created_by")
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
		if err := rows.Scan(&fe.Name, &fe.FeatureID, &fe.DBType, &fe.CreatedAt, &fe.CreatedBy); err != nil {
			return nil, err
		}

		envLists = append(envLists, &fe)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return envLists, nil
}

func CreateFeatureEnvironment(fe data.FeatureEnvironment) error {
	// Insert into table
	db := db.DB()
	dq := goqu.Insert("feature_environments").
		Cols("name", "feature_id", "db_type", "created_at", "created_by").
		Vals(goqu.Vals{fe.Name, fe.FeatureID, fe.DBType, fe.CreatedAt, fe.CreatedBy})

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

	// Iterate over resources and insert them into the database
	for _, resource := range fe.Resources {
		resource.FeatureEnvID = feID
		if err := insertResource(resource); err != nil {
			return err
		}
	}

	return nil
}

func insertResource(resource data.Resource) error {
	// Insert Resource into database
	db := db.DB()

	dq := goqu.Insert("resources").
		Cols("app_name", "feature_environment_id", "is_auto_update", "link").
		Vals(goqu.Vals{resource.AppName, resource.FeatureEnvID, resource.IsAutoUpdate, resource.Link})

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

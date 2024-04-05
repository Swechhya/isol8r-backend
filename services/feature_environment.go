package services

import (
	"time"

	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
)

func GetAllFeatureEnvironments() ([]*data.FeatureEnvironment, error) {
	//Dummy data
	envLists := []*data.FeatureEnvironment{
		{
			Name:      "Environment 1",
			DBType:    "mongodb",
			CreatedAt: time.Now().Format(time.RFC3339),
			CreatedBy: "User 1",
			Resources: []data.Resource{
				{
					AppName:      "App 1",
					IsAutoUpdate: true,
				},
				{
					AppName:      "App 2",
					IsAutoUpdate: false,
				},
			},
		},
		{
			Name:      "Environment 2",
			DBType:    "mysql",
			CreatedAt: time.Now().Format(time.RFC3339),
			CreatedBy: "User 2",
			Resources: []data.Resource{
				{
					AppName:      "App 3",
					IsAutoUpdate: true,
				},
				{
					AppName:      "App 4",
					IsAutoUpdate: true,
				},
			},
		},
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

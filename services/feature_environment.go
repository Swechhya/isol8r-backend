package services

import (
	"time"

	"github.com/Swechhya/panik-backend/data"
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

package services

import (
	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/internal/db"
)

func AddNewConfig(config data.Config) error {
	// Insert into table
	db := db.DB()
	_, err := db.Exec(`
        INSERT INTO core_config (key, value, created_by)
        VALUES ($1, $2, $3)
    `, config.Key, config.Value, config.CreatedBy)

	if err != nil {
		return err
	}

	return nil
}

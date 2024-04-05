package services

import (
	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
)

func AddNewConfig(config data.Config) error {
	// Insert into table
	db := db.DB()

	dq := goqu.Insert("core_config").
		Cols("key", "value").
		Vals(goqu.Vals{config.Key, config.Value})

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

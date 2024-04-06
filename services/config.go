package services

import (
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/doug-martin/goqu/v9"
)

var coreConfig map[string]string

func GetConfig(key string) string {
	return coreConfig[key]
}

func saveConfig(key, value string) {
	coreConfig[key] = value
}
func AddConfig(key, value string) error {
	// Insert into table
	db := db.DB()
	k := GetConfig(key)

	dq := goqu.Insert("core_config").
		Cols("key", "value").
		Vals(goqu.Vals{key, value})
	sql, args, err := dq.ToSQL()

	if k != "" {
		dq := goqu.Update("core_config").
			Set(goqu.Record{"value": value}).
			Where(goqu.Ex{
				"key": key})
		sql, args, err = dq.ToSQL()

	}

	if err != nil {
		return err
	}
	_, err = db.Exec(sql, args...)

	if err != nil {
		return err
	}

	saveConfig(key, value)

	return nil
}

func LoadConfigOnInitialSetup() error {
	coreConfig = make(map[string]string, 0)
	dq := goqu.From("core_config").Select("key", "value")

	sql, args, err := dq.ToSQL()
	for err != nil {
		return err
	}

	db := db.DB()

	rows, err := db.Query(sql, args...)
	for err != nil {
		return err
	}

	defer rows.Close()

	var key, value string
	for rows.Next() {
		if err := rows.Scan(&key, &value); err != nil {
			return err
		}
		saveConfig(key, value)
	}

	return nil

}

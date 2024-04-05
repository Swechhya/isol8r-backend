package db

import (
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool

func SetupDB(url, dbName, sslmode string) error {
	connURL := fmt.Sprintf("%s/%s?sslmode=%s", url, dbName, sslmode)
	connectionConfig, err := pgx.ParseURI(connURL)
	if err != nil {
		return err
	}

	maxConnections := 50
	timeOut := 5 * time.Minute
	poolConfig := pgx.ConnPoolConfig{ConnConfig: connectionConfig, MaxConnections: maxConnections, AfterConnect: nil, AcquireTimeout: timeOut}
	pgxPool, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		return err
	}

	pool = pgxPool
	fmt.Printf("Connected to db using sslmode=%s", sslmode)

	return nil
}

func DB() *pgx.ConnPool {

	if pool == nil {
		url := "postgres://postgres:root@localhost:5432"
		dbName := "panik_fe_db"
		sslmode := "disable"

		err := SetupDB(url, dbName, sslmode)
		if err != nil {
			fmt.Printf("Error connecting to db. Error: %s", err)
			return nil
		}
	}

	return pool
}

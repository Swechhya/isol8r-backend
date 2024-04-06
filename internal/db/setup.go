package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

var pool *pgx.ConnPool

func SetupDB(url, dbName, sslmode, sslrootcert string) error {
	connURL := fmt.Sprintf("%s/%s?sslmode=%s", url, dbName, sslmode)

	if sslmode != "disable" && strings.Trim(sslrootcert, " ") != "" {
		connURL = fmt.Sprintf("%s&sslrootcert=%s", connURL, sslrootcert)
	}

	connectionConfig, err := pgx.ParseURI(connURL)

	if err != nil {
		return err
	}

	maxConnections := 50
	timeOut := 5 * time.Minute

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connectionConfig,
		MaxConnections: maxConnections,
		AfterConnect:   nil,
		AcquireTimeout: timeOut,
	}

	pgxPool, err := pgx.NewConnPool(poolConfig)

	if err != nil {
		return err
	}

	pool = pgxPool
	fmt.Printf("Connected to postgres using sslmode=%s", sslmode)

	return nil
}

func DB() *pgx.ConnPool {

	if pool == nil {
		url := os.Getenv("DB_URL")
		dbName := os.Getenv("DB_NAME")
		sslmode := os.Getenv("SSL_MODE")
		sslrootcert := os.Getenv("SSL_ROOT_CERT")

		err := SetupDB(url, dbName, sslmode, sslrootcert)

		if err != nil && sslmode != "disable" {
			fmt.Printf("Error connecting to postgres using sslmode=%s. Falling back to sslmode=disable", sslmode)
			err = SetupDB(url, dbName, "disable", "")
		}

		if err != nil {
			return nil
		}

	}

	return pool
}

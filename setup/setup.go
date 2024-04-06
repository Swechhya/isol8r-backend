package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/Swechhya/panik-backend/services"
	"github.com/joho/godotenv"
)

func Setup() {
	loadEnv()

	err := setupDatabase()
	if err != nil {
		log.Panic("Unable to load config")
	}

	err = services.LoadConfigOnInitialSetup()
	if err != nil {
		log.Panic("Unable to load config")
	}

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Panicf("Unable to load env, err: %s", err)
	}
}

func setupDatabase() error {
	url := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")
	sslrootcert := os.Getenv("SSL_ROOT_CERT")

	err := db.SetupDB(url, dbName, sslmode, sslrootcert)

	// Fallback to ssl disable
	if err != nil && sslmode != "disable" {
		fmt.Printf("Error connecting to postgres using sslmode=%s. Trying with sslmode=disable", sslmode)
		err = db.SetupDB(url, dbName, "disable", "")
	}

	if err != nil {
		log.Panicf("Error connecting to report db. Error message: %s", err)
		return nil
	}

	return nil
}

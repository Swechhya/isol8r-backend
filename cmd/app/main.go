package main

import (
	"net/http"

	"github.com/Swechhya/panik-backend/handlers"
	"github.com/Swechhya/panik-backend/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {

	// d, _ := services.GetRepos(context.TODO())
	// fmt.Println(d)
	err := Setup()
	if err != nil {
		return
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server Running...")
	})

	router.POST("/config/create", handlers.ConfigHandler)
	router.POST("/repos", handlers.GetRepos)

	//feature-enviroment-handler
	featureEnvironment := router.Group("/fe")
	{
		featureEnvironment.GET("/apps", handlers.AppListHandler)
		featureEnvironment.GET("/list", handlers.FEListHandler)
		featureEnvironment.POST("/create", handlers.FECreateHandler)
		featureEnvironment.POST("/delete", handlers.FEDeleteHandler)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

func Setup() error {
	url := "postgres://postgres:root@localhost:5432"
	dbName := "panik_fe_db"
	sslMode := "disable"

	err := db.SetupDB(url, dbName, sslMode)
	if err != nil {
		return err
	}

	return nil
}

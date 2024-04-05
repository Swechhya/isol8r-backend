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
	router.Use(corsHandler())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server Running...")
	})

	router.POST("/config/create", handlers.ConfigHandler)
	router.POST("/repos", handlers.GetRepos)

	gh := router.Group("/gh")
	{
		gh.POST("/setup", handlers.SetupGithub)
		gh.POST("/repos", handlers.GetRepos)
		gh.GET("/branches/:repo", handlers.GetBranches)
	}

	//feature-enviroment-handler
	featureEnvironment := router.Group("/fe")
	{
		featureEnvironment.GET("/list", handlers.FEListHandler)
		featureEnvironment.POST("/create", handlers.FECreateHandler)
		featureEnvironment.POST("/delete", handlers.FEDeleteHandler)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

func corsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
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

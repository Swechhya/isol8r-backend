package main

import (
	"net/http"

	"github.com/Swechhya/panik-backend/handlers"
	"github.com/Swechhya/panik-backend/middleware"
	"github.com/Swechhya/panik-backend/setup"
	"github.com/gin-gonic/gin"
)

func main() {
	setup.Setup()

	router := gin.Default()
	router.Use(middleware.CorsHandler())
	router.Use(middleware.RcoveryHandler())

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
}

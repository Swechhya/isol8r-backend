package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Swechhya/isol8r-backend/handlers"
	"github.com/Swechhya/isol8r-backend/middleware"
	"github.com/Swechhya/isol8r-backend/services"
	"github.com/Swechhya/isol8r-backend/setup"
	"github.com/gin-gonic/gin"
)

func main() {
	// TestRun()
	setup.Setup()

	router := gin.Default()
	router.Use(middleware.CorsHandler())
	router.Use(middleware.RcoveryHandler())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server Running...")
	})

	gh := router.Group("/gh")
	{
		gh.POST("/setup", handlers.SetupGithub)
		gh.GET("/repos", handlers.GetRepos)
		gh.GET("/branches/:repoId", handlers.GetBranches)
		gh.POST("/save-env/:repoId", handlers.UploadEnvFile)
	}

	//feature-enviroment-handler
	featureEnvironment := router.Group("/fe")
	{
		featureEnvironment.GET("/list", handlers.FEListHandler)
		featureEnvironment.POST("/create", handlers.FECreateHandler)
		featureEnvironment.POST("/delete/:id", handlers.FEDeleteHandler)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
}

func TestRun() error {

	ecr := "654654451390.dkr.ecr.us-east-1.amazonaws.com/test:"

	repoName := strings.Split("Swechhya/Contact-Manager-Frontend", "/")[1]
	dest := fmt.Sprintf("%s%s-%s", ecr, "feature-new-test", repoName)
	err := services.GenerateBuildManifest("feature-new-test", "Swechhya/Contact-Manager-Frontend", "main", dest)
	if err != nil {
		return err
	}
	err = services.GenerateDeployManifest("feature-new-test", dest)
	if err != nil {
		return err
	}
	err = services.DeployEnvironment("feature-new-test")
	if err != nil {
		return err
	}
	return nil
}

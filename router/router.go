package router

import (
	"net/http"

	"github.com/Swechhya/isol8r-backend/handlers"
	"github.com/Swechhya/isol8r-backend/middleware"
	"github.com/Swechhya/isol8r-backend/setup"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
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
		featureEnvironment.GET("/repos", handlers.FEGetRepoHandler)
		featureEnvironment.GET("/list", handlers.FEListHandler)
		featureEnvironment.POST("/create", handlers.FECreateHandler)
		featureEnvironment.POST("/delete/:id", handlers.FEDeleteHandler)
		featureEnvironment.POST("/redeploy/:id", handlers.FERedeployHandler)
		featureEnvironment.POST("/edit/:id", handlers.FEEditHandler)
	}

	return router
}

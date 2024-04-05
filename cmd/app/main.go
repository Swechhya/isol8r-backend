package main

import (
	"github.com/Swechhya/panik-backend/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/gh-auth", handlers.GitHubAuthorizeHandler)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

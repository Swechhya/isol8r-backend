package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Swechhya/panik-backend/services"
)

func GitHubAuthorizeHandler(c *gin.Context) {

	err := services.AuthenticateGitHub()
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

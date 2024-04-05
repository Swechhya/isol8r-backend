package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppListHandler(c *gin.Context) {
	//TODO : LIST REPOS
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func FEListHandler(c *gin.Context) {
	//TODO : LIST FEATURE ENVIRONMENTS
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func FECreateHandler(c *gin.Context) {
	//TODO : CREATE FEATURE ENVIRONMENTS
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func FEDeleteHandler(c *gin.Context) {
	//TODO :: DELETE FEATURE ENVIRONMENTS
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

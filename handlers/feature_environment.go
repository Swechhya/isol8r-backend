package handlers

import (
	"net/http"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/services"
	"github.com/gin-gonic/gin"
)

func AppListHandler(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func FEListHandler(c *gin.Context) {
	//TODO : LIST FEATURE ENVIRONMENTS
	fe, err := services.GetAllFeatureEnvironments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   fe,
	})
}

func FECreateHandler(c *gin.Context) {
	var fe data.FeatureEnvironment
	if err := c.ShouldBindJSON(&fe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateFeatureEnvironment(fe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"resource": fe,
	})
}

func FEDeleteHandler(c *gin.Context) {
	//TODO :: DELETE FEATURE ENVIRONMENTS
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

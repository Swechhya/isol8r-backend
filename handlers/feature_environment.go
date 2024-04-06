package handlers

import (
	"net/http"
	"strconv"

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
	fe, err := services.GetAllFeatureEnvironments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "ERROR",
		})
	}

	SuccessResponse(c, fe)
}

func FECreateHandler(c *gin.Context) {
	var fe data.FeatureEnvironment
	if err := c.ShouldBindJSON(&fe); err != nil {
		ErrorReponse(c, err)
		return
	}

	if err := services.CreateFeatureEnvironment(fe); err != nil {
		ErrorReponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

func FEDeleteHandler(c *gin.Context) {
	feID := c.Param("id")
	if err := services.DeleteFeatureEnvironment(feID); err != nil {
		ErrorReponse(c, err)
		return
	}
}

func FERedeployHandler(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		ErrorReponse(c, err)
		return
	}

	fe, err := services.GetFeatureEnvironmentById(id)
	if err != nil {
		ErrorReponse(c, err)
		return
	}

	if err := services.CreateFeatureEnvironment(*fe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	SuccessResponse(c, fe)
}

func FEEditHandler(c *gin.Context) {
	// TODO: DELETE FEATURE ENV

	// Recreate feature env
	var fe data.FeatureEnvironment
	if err := c.ShouldBindJSON(&fe); err != nil {
		ErrorReponse(c, err)
		return
	}

	if err := services.CreateFeatureEnvironment(fe); err != nil {
		ErrorReponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

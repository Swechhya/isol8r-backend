package handlers

import (
	"net/http"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/services"
	"github.com/gin-gonic/gin"
)

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

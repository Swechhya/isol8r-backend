package handlers

import (
	"strconv"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/services"
	"github.com/gin-gonic/gin"
)

func FEGetRepoHandler(c *gin.Context) {
	repos, err := services.FetchLaunchReadyRepos(c)
	if err != nil {
		ErrorResponse(c, err)
		return
	}
	SuccessResponse(c, repos)
}

func FEListHandler(c *gin.Context) {
	fe, err := services.GetAllFeatureEnvironments()
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

func FEDetailsHandler(c *gin.Context) {
	ids := c.Param("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	fe, err := services.GetFeatureEnvironmentById(id)
	if err != nil {
		ErrorResponse(c, err)
		return
	}
	SuccessResponse(c, fe)
}

func FECreateHandler(c *gin.Context) {
	var fe data.FeatureEnvironment
	if err := c.ShouldBindJSON(&fe); err != nil {
		ErrorResponse(c, err)
		return
	}

	isReDeploy := false
	_, err := services.CreateFeatureEnvironment(fe, isReDeploy)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

func FEDeleteHandler(c *gin.Context) {
	feID := c.Param("id")
	id, err := strconv.Atoi(feID)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if err := services.DeleteFeatureEnvironment(id); err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, "OK")
}

func FERedeployHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	fe, err := services.GetFeatureEnvironmentById(id)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	isReDeploy := true
	fe, err = createOrUpdateFeatureEnvironment(fe, isReDeploy)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

func FEEditHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	if err := services.DeleteFeatureEnvironment(id); err != nil {
		ErrorResponse(c, err)
		return
	}

	var fe *data.FeatureEnvironment
	if err := c.ShouldBindJSON(&fe); err != nil {
		ErrorResponse(c, err)
		return
	}

	isReDeploy := false
	fe, err = createOrUpdateFeatureEnvironment(fe, isReDeploy)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, fe)
}

func createOrUpdateFeatureEnvironment(fe *data.FeatureEnvironment, isReDeploy bool) (*data.FeatureEnvironment, error) {
	feID, err := services.CreateFeatureEnvironment(*fe, isReDeploy)
	if err != nil {
		return nil, err
	}

	fe.ID = feID
	return fe, nil
}

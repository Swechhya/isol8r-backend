package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Swechhya/panik-backend/services"
)

func ErrorReponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func SetupGithub(c *gin.Context) {
	err := services.SetupGithubClient()
	if err != nil {
		ErrorReponse(c, err)
		return
	}
	SuccessResponse(c, "success")
	return
}

func GetRepos(c *gin.Context) {
	repos, err := services.GetRepos(c)
	if err != nil {
		ErrorReponse(c, err)
		return
	}
	SuccessResponse(c, repos)
	return
}

func GetBranches(c *gin.Context) {
	var body map[string]string
	if err := c.ShouldBindJSON(&body); err != nil {
		ErrorReponse(c, err)
		return
	}
	b := body["branch"]
	branches, err := services.GetBranches(c, b)
	if err != nil {
		ErrorReponse(c, err)
	}
	SuccessResponse(c, branches)
}

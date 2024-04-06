package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Swechhya/panik-backend/data"
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
	var config *data.GithubClientSetup
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.SetupGithubClient(c.Request.Context(), config)
	if err != nil {
		ErrorReponse(c, err)
		return
	}

	SuccessResponse(c, "success")
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
	repo := c.Param("repo")
	branches, err := services.GetBranches(c, repo)
	if err != nil {
		ErrorReponse(c, err)
	}
	SuccessResponse(c, branches)
}

func UploadEnvFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer fileReader.Close()

	// Read the file into a buffer
	fileBytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := services.UploadEnvFile(c, bytes.NewReader(fileBytes)); err != nil {
		return
	}
}

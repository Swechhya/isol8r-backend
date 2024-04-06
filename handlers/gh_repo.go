package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/services"
	"github.com/Swechhya/panik-backend/services/s3"
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

	bucketName := "panik-private-key"
	bucketKey := "key.pem"
	filePath, _ := os.Getwd()
	fullFilePath := filepath.Join(filePath, bucketKey)

	client := s3.GetClient()
	err := client.DownloadFileToPath(c.Request.Context(), bucketName, bucketKey, fullFilePath)
	if err != nil {
		ErrorReponse(c, err)
		return
	}

	config.PrivateKeyPath = fullFilePath
	err = services.SetupGithubClient(c.Request.Context(), config)
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
	repoId := c.Param("repoId")
	branches, err := services.GetBranches(c, repoId)
	if err != nil {
		ErrorReponse(c, err)
	}
	SuccessResponse(c, branches)
}

func UploadEnvFile(c *gin.Context) {
	repoId := c.Param("repoId")
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

	uri, err := services.UploadEnvFile(c, bytes.NewReader(fileBytes), repoId)
	if err != nil {
		return
	}

	SuccessResponse(c, uri)
}

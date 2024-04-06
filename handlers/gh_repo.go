package handlers

import (
	"bytes"
	"fmt"
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
	accessKey := "ASIAZQ3DSF27AOODHVWH"
	secretKey := "M6yMnGPB7hLf9YdVxEbatoSfS16eA0Q915ZayfMh"
	region := "us-east-1"
	filePath, _ := os.Getwd()
	fullFilePath := filepath.Join(filePath, bucketKey)
	client := s3.GetClient(accessKey, secretKey, region)
	err := client.DownloadFileToPath(c.Request.Context(), bucketName, bucketKey, fullFilePath)
	if err != nil {
		ErrorReponse(c, err)
		return
	}

	config.PrivateKey = fullFilePath
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

	uri, err := services.UploadEnvFile(c, bytes.NewReader(fileBytes))
	if err != nil {
		return
	}

	fmt.Println(uri)
}

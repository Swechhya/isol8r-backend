package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Swechhya/isol8r-backend/data"
	"github.com/Swechhya/isol8r-backend/services"
	"github.com/Swechhya/isol8r-backend/services/s3"
)

func ErrorResponse(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status": "ERROR",
		"error":  err.Error(),
	})
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   data,
	})
}

func SetupGithub(c *gin.Context) {
	var config *data.GithubClientSetup
	if err := c.ShouldBindJSON(&config); err != nil {
		ErrorResponse(c, err)
	}

	bucketName := os.Getenv("PRIVATE_KEY_BUCKET")
	fileName := os.Getenv("PRIVATE_KEY_FILENAME")
	filePath, _ := os.Getwd()
	fullFilePath := filepath.Join(filePath, fileName)

	client := s3.GetClient()
	err := client.DownloadFileToPath(c.Request.Context(), bucketName, fileName, fullFilePath)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	config.PrivateKeyPath = fullFilePath
	err = services.SetupGithubClient(c.Request.Context(), config)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, "success")
}

func GetRepos(c *gin.Context) {
	repos, err := services.GetRepos(c)
	if err != nil {
		ErrorResponse(c, err)
		return
	}
	SuccessResponse(c, repos)
}

func GetBranches(c *gin.Context) {
	repoParam := c.Param("repoId")
	repoId, err := strconv.ParseInt(repoParam, 10, 64)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	branches, err := services.GetBranches(c, repoId)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	SuccessResponse(c, branches)
}

func UploadEnvFile(c *gin.Context) {
	repoId := c.Param("repoId")
	file, err := c.FormFile("file")
	fileName := file.Filename
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		ErrorResponse(c, err)
		return
	}
	defer fileReader.Close()

	// Read the file into a buffer
	fileBytes, err := ioutil.ReadAll(fileReader)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	uri, err := services.UploadEnvFile(c, bytes.NewReader(fileBytes), repoId, fileName)
	if err != nil {
		return
	}

	SuccessResponse(c, uri)
}

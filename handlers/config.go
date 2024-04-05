package handlers

import (
	"net/http"

	"github.com/Swechhya/panik-backend/data"
	"github.com/Swechhya/panik-backend/services"
	"github.com/gin-gonic/gin"
)

func ConfigHandler(c *gin.Context) {
	var config data.Config
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.AddNewConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"data":   config,
	})
}

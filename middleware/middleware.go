package middleware

import (
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func RcoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rc := recover(); rc != nil {
				var err error
				switch x := rc.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("unknown panic")
				}
				if err != nil {

					log.Fatal("PANIC: %s", err.Error())

					if etrace := debug.Stack(); etrace != nil {
						log.Fatal("STACKTRACE: %s", etrace)
					}
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
		}()

		c.Next()
	}
}

package loggers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Started %s %s", c.Request.Method, c.Request.URL.Path)

		c.Next()

		log.Printf("Completed %s %s %d", c.Request.Method, c.Request.URL.Path, c.Writer.Status())
	}
}

func RecoveryWithLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic occurred: %v", err)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

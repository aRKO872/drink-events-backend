package routers

import (
	loadenv "github.com/drink-events-backend/cmd/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Middlewares
	r.Use(loadenv.LoadEnv)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
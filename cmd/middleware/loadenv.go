package loadenv

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func LoadEnv (c *gin.Context) {
	godotenv.Load(".env")
	c.Next()
}
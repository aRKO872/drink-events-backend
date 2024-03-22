package middlewares

import (
	"fmt"

	internal_database "github.com/drink-events-backend/internal"
	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func LoadDatabase (c *gin.Context) {
	connString := pkg_config.GetProjectConfig().DATABASE_URL
	_, dbErr := internal_database.GetDB(connString)

	// Initializing Database
	if dbErr != nil {
		c.JSON(200, gin.H{
			"status": false,
			"errMessage": dbErr.Error(),
		})
		c.Abort()
		return
	}

	// Initializing Redis
	if _, err := internal_database.GetRDB(); err != nil {
		c.JSON(200, gin.H{
			"status": false,
			"errMessage": err.Error(),
		})
		c.Abort()
		return
	}

	fmt.Println("DB and RDB is set up!")
	c.Next()
}
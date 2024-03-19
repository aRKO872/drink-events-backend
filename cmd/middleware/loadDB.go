package middlewares

import (
	"fmt"
	"os"

	internal_database "github.com/drink-events-backend/internal"
	"github.com/gin-gonic/gin"
)

func LoadDatabase (c *gin.Context) {
	connString := os.Getenv("DATABASE_URL")
	dbErr := internal_database.InitDB(connString)

	if dbErr != nil {
		c.JSON(200, gin.H{
			"status": false,
			"errMessage": dbErr.Error(),
		})
		c.Abort()
	} else {
		fmt.Println("DB is set up!")
		c.Next()
	}
}
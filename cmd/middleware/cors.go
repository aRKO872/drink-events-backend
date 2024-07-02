package middlewares

import (
	"net/http"

	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/gin-gonic/gin"
)

func EnableCors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", pkg_config.GetProjectConfig().FE_SOURCE)
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Request.Header.Del("Origin")

	c.Next()
}
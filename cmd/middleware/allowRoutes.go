package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AllowRoutesMiddleware(
	middleware gin.HandlerFunc,
	routesToAvoid... string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowRoutes := false
		for _, route := range routesToAvoid {
			allowRoutes = allowRoutes || strings.HasPrefix(c.Request.URL.Path, route)
		}

		if allowRoutes {
			middleware(c)
		} else {
			c.Next()
		}
	}
}
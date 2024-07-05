package middlewares

import (
	"fmt"
	"net/http"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func (c *gin.Context) {
		user, fetchUserErr := service.FetchUserFromToken(c.Request, true)
		if fetchUserErr != nil {
			c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
				Status: false,
				ErrorMsg: fetchUserErr.Error(),
			})
			c.Abort()
			return
		}

		c.Request.Header.Set(literals.HEADER_USER_ID, user.Id)
		c.Request.Header.Set(literals.HEADER_USER_SEARCH_DISTANCE, fmt.Sprintf("%d", user.SearchRadius))
		c.Request.Header.Set(literals.HEADER_USER_TYPE, user.UserType)

		c.Next()
	}
}
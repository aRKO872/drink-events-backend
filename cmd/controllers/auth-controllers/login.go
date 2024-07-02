package auth_controller

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gin-gonic/gin"
)

func LogIn(c *gin.Context) {
	var input *models.LogInInput
	bindingErr := c.Bind(&input)
	ctx := c.Request.Context()

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	user := &models.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	logInStatus, logInObj := service.LogIn(ctx, input, user)

	if !logInStatus {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: logInObj.Status,
			ErrorMsg: logInObj.ErrorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, &logInObj)
}
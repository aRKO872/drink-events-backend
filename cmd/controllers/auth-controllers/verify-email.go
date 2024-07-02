package auth_controller

import (
	"fmt"
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	var input *models.VerifyUserEmailInput

	ctx := c.Request.Context()
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindDataErr.Error(),
		})
		return
	}

	user := &models.Users{
		Email: input.Email,
	}

	isSuccessful, emailSendErr := service.VerifyEmail(ctx, user);

	if !isSuccessful && emailSendErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: emailSendErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : true,
		"msg": fmt.Sprintf("email successfully sent to %s", input.Email),
	})
}
package auth_controller

import (
	"fmt"
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gin-gonic/gin"
)

func VerifyOTP(c *gin.Context) {
	var input *models.VerifyOTP
	ctx := c.Request.Context()
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindDataErr.Error(),
		})
		return
	}

	// Business Logic
	user := &models.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	isSuccessful, otpCheckErr := service.VerifyOTP(ctx, input.Otp, input.Event, user);

	if !isSuccessful && otpCheckErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: otpCheckErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status" : true,
		"msg": fmt.Sprintf("%s is verified!", input.Email),
	})
}
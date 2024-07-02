package auth_controller

import (
	"fmt"
	"net/http"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/drink-events-backend/models"
	"github.com/gin-gonic/gin"
)

func ResendEmailOTP(c *gin.Context) {
	var input *models.ResendOTP
	ctx := c.Request.Context()

	if bindingErr := c.Bind(&input); bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	var user = &models.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	resendConfirmation, resendOTPError := service.ResendOTPForVerification(ctx, user, input)

	if !resendConfirmation && resendOTPError != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resendOTPError.Error(),
		})
		return
	}

	if input.Phone != "" {
		c.JSON(http.StatusOK, gin.H{
			"status" : true,
			"msg": fmt.Sprintf("otp successfully resent to %s", input.Phone),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status" : true,
			"msg": fmt.Sprintf("otp successfully resent to %s", input.Email),
		})
	}
}
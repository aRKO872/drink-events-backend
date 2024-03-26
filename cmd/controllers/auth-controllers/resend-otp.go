package auth_controller

import (
	"fmt"
	"net/http"

	pkg_component_user "github.com/drink-events-backend/pkg/components/user_component"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	"github.com/gin-gonic/gin"
)

func ResendEmailOTP(c *gin.Context) {
	var input *endpoint_inputs.ResendOTP

	if bindingErr := c.Bind(&input); bindingErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	var user = &pkg_component_user.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	resendConfirmation, resendOTPError := user.ResendOTPForVerification(input)

	if !resendConfirmation && resendOTPError != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
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
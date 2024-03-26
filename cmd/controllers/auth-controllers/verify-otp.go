package auth_controller

import (
	"fmt"
	"net/http"

	pkg_component "github.com/drink-events-backend/pkg/components/user_component"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	"github.com/gin-gonic/gin"
)

func VerifyOTP(c *gin.Context) {
	var input *endpoint_inputs.VerifyOTP
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindDataErr.Error(),
		})
		return
	}

	// Business Logic
	user := &pkg_component.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	isSuccessful, otpCheckErr := user.VerifyOTP(input.Otp, input.Event);

	if !isSuccessful && otpCheckErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
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
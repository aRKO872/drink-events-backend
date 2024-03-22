package auth_controller

import (
	"fmt"
	"net/http"

	pkg_component "github.com/drink-events-backend/pkg/components/user_component"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	"github.com/gin-gonic/gin"
)

func VerifyEmail(c *gin.Context) {
	var input *endpoint_inputs.VerifyUserEmailInput
	bindDataErr := c.Bind(&input);

	if bindDataErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindDataErr.Error(),
		})
		return
	}

	user := &pkg_component.Users{
		Email: input.Email,
	}

	isSuccessful, emailSendErr := user.VerifyEmail();

	if !isSuccessful && emailSendErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
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
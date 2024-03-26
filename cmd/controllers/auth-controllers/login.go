package auth_controller

import (
	"net/http"

	pkg_component_user "github.com/drink-events-backend/pkg/components/user_component"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	"github.com/gin-gonic/gin"
)

func LogIn(c *gin.Context) {
	var input *endpoint_inputs.LogInInput
	bindingErr := c.Bind(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	user := &pkg_component_user.Users{
		Email: input.Email,
		Phone: input.Phone,
	}

	logInStatus, logInObj := user.LogIn(input)

	if !logInStatus {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: logInObj.Status,
			ErrorMsg: logInObj.ErrorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, &logInObj)
}
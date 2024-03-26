package auth_controller

import (
	"net/http"

	pkg_component_user "github.com/drink-events-backend/pkg/components/user_component"
	endpoint_inputs "github.com/drink-events-backend/pkg/endpoint-inputs"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var input *endpoint_inputs.SignUpInput

	if bindingErr := c.Bind(&input); bindingErr != nil {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	var user = &pkg_component_user.Users{
		Email: input.Email,
		Name: input.Name,
		Phone: input.Phone,
		UserType: "user",
		Bio: input.Bio,
	}

	signUpStatus, signUpObj := user.SignUp()

	if !signUpStatus {
		c.JSON(http.StatusBadRequest, &endpoint_inputs.CommonErrorOutput{
			Status: signUpObj.Status,
			ErrorMsg: signUpObj.ErrorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, &signUpObj)
}
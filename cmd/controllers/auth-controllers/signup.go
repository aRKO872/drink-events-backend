package auth_controller

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var input *models.SignUpInput
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
		Name: input.Name,
		Phone: input.Phone,
		UserType: "user",
		Bio: input.Bio,
	}

	signUpStatus, signUpObj := service.SignUp(ctx, user)

	if !signUpStatus {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: signUpObj.Status,
			ErrorMsg: signUpObj.ErrorMsg,
		})
		return
	}

	c.JSON(http.StatusOK, &signUpObj)
}
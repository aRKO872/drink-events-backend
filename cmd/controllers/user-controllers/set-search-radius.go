package user_controllers

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
)

func SetSearchRadius(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	var input *models.SearchRadiusSetInput
	bindingErr := c.ShouldBindJSON(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	searchRadiusOutput := service.SetSearchRadius(
		c.Request.Context(),
		tokenUser.UserId,
		input.SearchRadius,
	)

	if !searchRadiusOutput.Status {
		c.JSON(http.StatusBadRequest, &searchRadiusOutput)
		return
	}

	c.JSON(http.StatusOK, &searchRadiusOutput)
}
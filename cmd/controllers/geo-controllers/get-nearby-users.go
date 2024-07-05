package geo_controllers

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
)

func GetNearbyUsers(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	var input *models.LocationSetEvent
	bindingErr := c.ShouldBindJSON(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	isSuccessful, nearbyUsersOutput := service.GetPeopleNearby(c.Request, input, tokenUser)

	if !isSuccessful {
		c.JSON(http.StatusBadRequest, &nearbyUsersOutput)
		return
	}

	c.JSON(http.StatusOK, &nearbyUsersOutput)
}
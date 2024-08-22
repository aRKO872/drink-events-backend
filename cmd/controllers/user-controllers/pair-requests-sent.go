package user_controllers

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
)

func PairRequestsSentByMe(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	var input *models.PairRequestsSentInput
	bindingErr := c.ShouldBindJSON(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	isSuccessful, prSentByMe := service.GetPRSentByMe(*input, tokenUser.UserId)

	if !isSuccessful {
		c.JSON(http.StatusBadRequest, &prSentByMe)
		return
	}

	c.JSON(http.StatusOK, &prSentByMe)
}

func PairRequestsSentToMe(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	var input *models.PairRequestsSentInput
	bindingErr := c.ShouldBindJSON(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	isSuccessful, prSentByMe := service.GetPRSentToMe(*input, tokenUser.UserId)

	if !isSuccessful {
		c.JSON(http.StatusBadRequest, &prSentByMe)
		return
	}

	c.JSON(http.StatusOK, &prSentByMe)
}
package user_controllers

import (
	"net/http"

	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
)


func ChangeProfilePicture(c *gin.Context) {
	tokenizedUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	imgFile, imgDetails, err := pkg_helpers.ValidateImageRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: err.Error(),
		})
		return
	}

	savedProfilePicData := service.SaveProfilePicture(c.Request.Context(), tokenizedUser, imgFile, imgDetails)
	
	c.JSON(http.StatusOK, &savedProfilePicData)
}
package websockets

import (
	"encoding/json"
	"net/http"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	"github.com/drink-events-backend/pkg/business_logic/service"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
)

func (m *Manager) AcceptPairRequest(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	var input *models.AcceptPairRequestInput
	bindingErr := c.ShouldBindJSON(&input)

	if bindingErr != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: bindingErr.Error(),
		})
		return
	}

	// Update Pair Request isActive = false in DB
	// Delete Pair Request Object from Redis
	// Create new Friends Object, persist it in DB
	// Send notification to sender of pair request that they're friends

	friendRecord, wsOutput, err := service.CreateNewFriend(c.Request.Context(), tokenUser.UserId, *input)
	if err != nil {
		c.JSON(http.StatusBadRequest, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: err.Error(),
		})
		return
	}

	if senderClient, ok := m.ClientMap[input.SenderId]; ok {
		// Form positive pair request accepted message event evt
		wsOutputBytes, _ := json.Marshal(wsOutput)

		evt := Event{
			Type: literals.PAIR_REQUEST_ACCEPTED,
			Payload: wsOutputBytes,
		}
		senderClient.Feed <- evt
	}

	c.JSON(http.StatusOK, &models.AcceptPairRequestResponse{
		Status: true,
		ErrorMsg: "friend mapping successful",
		FriendRecord: models.Friends{
			Id: friendRecord.Id,
			IsActive: friendRecord.IsActive,
			CreatedAt: friendRecord.CreatedAt,
		},
	})
}
package websockets

import (
	"encoding/json"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
)

func FormNotifyMessageEvent(msg string) (Event, error) {
	msgPayload := models.NotifyMsgEvent{
		Msg: msg,
	}
	payload, err := json.Marshal(msgPayload)
	if err != nil {
		return Event{}, err
	}

	return Event{
		Type: literals.NOTIFY_MESSAGE,
		Payload: payload,
	}, nil
}
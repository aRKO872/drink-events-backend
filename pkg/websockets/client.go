package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/drink-events-backend/models"
	"github.com/gorilla/websocket"
)

var (
	pongWait = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type Client struct {
	models.TokenizedUserDetails
	Conn *websocket.Conn
	Manager *Manager
	Feed chan Event
}

func NewClient(
	userDtls models.TokenizedUserDetails,
	m *Manager,
	conn *websocket.Conn,
) *Client {
	return &Client{
		TokenizedUserDetails: userDtls,
		Manager: m,
		Conn: conn,
		Feed: make(chan Event),
	}
}

func (c *Client) ReadMessages() {
	defer c.Manager.RemoveClient(c.UserId)
	c.Conn.SetReadLimit(10 << 10) // 10kb

	if err := c.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	c.Conn.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				evt, _ := FormNotifyMessageEvent(fmt.Sprintf("error reading message: %s", err.Error()))
				c.Feed <- evt
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			evt, _ := FormNotifyMessageEvent(fmt.Sprintf("error marshalling message: %s", err.Error()))
			c.Feed <- evt
			continue
		}

		if err := c.Manager.RouteEvent(request, c); err != nil {
			evt, _ := FormNotifyMessageEvent(err.Error())
			c.Feed <- evt
		}
	}
}

func (c *Client) pongHandler(pongMsg string) error {
	return c.Conn.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) WriteMessages() {
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		c.Manager.RemoveClient(c.UserId)
	}()

	for {
		select {
		case message, ok := <-c.Feed:
			// Ok will be false Incase the Feed channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend
				if err := c.Conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					evt, _ := FormNotifyMessageEvent(fmt.Sprintf("connection closed: %s", err.Error()))
					c.Feed <- evt
				}
				// Return to close the goroutine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				evt, _ := FormNotifyMessageEvent(fmt.Sprint("unable to marshall message"))
				c.Feed <- evt
				continue
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				evt, _ := FormNotifyMessageEvent(fmt.Sprintf("error sending message: %s", err.Error()))
				c.Feed <- evt
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}

	}
}
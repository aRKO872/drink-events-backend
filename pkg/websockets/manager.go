package websockets

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	pkg_helpers "github.com/drink-events-backend/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	ClientMap map[string]*Client
	Handlers map[string]EventHandler
	*sync.RWMutex
}

func (m *Manager) setupHandlers() {
	m.Handlers[literals.SEND_PAIR_REQUEST] = SendPairRequest
}

func NewWSManager() Manager {
	m := Manager {
		RWMutex: new(sync.RWMutex),
		ClientMap: make(map[string]*Client),
		Handlers: make(map[string]EventHandler),
	}
	m.setupHandlers()
	return m
}

func (m *Manager) AddClient(userId string, conn Client) {
	m.Lock()
	defer m.Unlock()

	m.ClientMap[userId] = &conn
}

func (m *Manager) RemoveClient(userId string) {
	m.Lock()
	defer m.Unlock()

	if client, conn := m.ClientMap[userId]; conn {
		client.Conn.Close()
		delete(m.ClientMap, userId)
	}
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *Manager) RouteEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.Handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("event not supported")
	}
}

func (m *Manager) ServeWS(c *gin.Context) {
	tokenUser, resolveFromHeaderErr := pkg_helpers.GetUserDetailsFromReqHeader(c.Request)

	if resolveFromHeaderErr != nil {
		c.JSON(http.StatusUnauthorized, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: resolveFromHeaderErr.Error(),
		})
		return
	}

	conn, wsUpgradeErr := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if wsUpgradeErr != nil {
		log.Println("error : ", wsUpgradeErr)
		c.JSON(http.StatusBadGateway, &models.CommonErrorOutput{
			Status: false,
			ErrorMsg: wsUpgradeErr.Error(),
		})
		return
	}

	clientObj := NewClient(*tokenUser, m, conn)
	m.AddClient(tokenUser.UserId, *clientObj)

	go clientObj.ReadMessages()
	go clientObj.WriteMessages()
}
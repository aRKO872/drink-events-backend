package pkg_component_websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/drink-events-backend/literals"
	"github.com/drink-events-backend/models"
	pkg_config "github.com/drink-events-backend/pkg/config"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserSocketClient struct {
	*websocket.Conn
	userId string
	send chan []byte
	hub *UserSocketHub
}

func (c *UserSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.Conn.Close()
	}()
	
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { 
		c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		log.Println("pong")
		return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.hub.broadcast <- message
	}
}

func (c *UserSocketClient) writePump() {
	pingTicker := time.NewTicker(pingPeriod)
	fetchNearbyPeriod := pkg_config.GetProjectConfig().FETCH_NEARBY_ACTIVE_PERIOD
	fetchNearbyTicker := time.NewTicker(time.Duration(fetchNearbyPeriod) * time.Second)
	
	defer func() {
		pingTicker.Stop()
		fetchNearbyTicker.Stop()
		c.Conn.Close()
	}()

	loop:
	for {
		select {
		case <-pingTicker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		
		case sentBytes := <-c.send :
			var socketEvtObj models.SocketEventMain
			if err := json.Unmarshal(sentBytes, &socketEvtObj); err != nil {
				log.Println("error while unmarshalling")
				continue loop
			}

			rso, fetchRSOErr := GetRedisSocketOperator()
			if fetchRSOErr != nil {
				log.Println("error while fetching Redis Socket Operator")
				continue loop
			} 
			if socketEvtObj.EventType == literals.SOCKET_LOCATION_SET_TYPE {
				rso.SetGEOLocation(socketEvtObj.LocationEvt, c.userId)
			}


		}
	}
}
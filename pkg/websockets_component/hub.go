package pkg_component_websockets

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/drink-events-backend/models"
)

type UserSocketHub struct {
	userClientMap map[string]*UserSocketClient
	register	chan *UserSocketClient
	unregister chan *UserSocketClient
	broadcast chan []byte
	sync.RWMutex
}

func NewHub() *UserSocketHub {
	return &UserSocketHub{
		userClientMap: make(map[string]*UserSocketClient),
		register: make(chan *UserSocketClient),
		unregister: make(chan *UserSocketClient),
		broadcast: make(chan []byte),
	}
}

func (h *UserSocketHub) run() {
	loop:
	for {
		select {
		case client := <-h.register :
			h.Lock()
			defer h.Unlock()
			h.userClientMap[client.userId] = client
		case client := <-h.unregister:
			h.Lock()
			defer h.Unlock()
			delete(h.userClientMap, client.userId)
			close(client.send)
		case sentbytes := <-h.broadcast: 
			var socketEvtObj models.SocketEventMain
			if err := json.Unmarshal(sentbytes, &socketEvtObj); err != nil {
				log.Println("error while unmarshalling")
				continue loop
			}

			
		}
	}
}
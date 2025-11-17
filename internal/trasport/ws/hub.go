package ws

import "go.uber.org/zap"

type Hub struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Clients    map[*Client]bool
	log        *zap.SugaredLogger
}

func NewHub(log *zap.SugaredLogger) *Hub {
	return &Hub{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
		Clients:    make(map[*Client]bool),
		log:        log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.log.Info("Registered client")
			h.Clients[c] = true
		case c := <-h.Unregister:
			h.log.Info("Unregistered client")
			if _, ok := h.Clients[c]; ok {
				delete(h.Clients, c)
				close(c.Send)
			}
		case msg := <-h.Broadcast:
			for c := range h.Clients {
				select {
				case c.Send <- msg:
				default:
					delete(h.Clients, c)
					close(c.Send)
				}
			}
		}
	}
}

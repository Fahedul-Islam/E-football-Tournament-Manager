package ws

type Hub struct {
	Clients    map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Notification
}

type Notification struct {
	UserID  int
	Message []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Notification),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				close(client.Send)
				_ = client.Conn.Close()
				delete(h.Clients, client.UserID)
			}
		case notification := <-h.Broadcast:
			if client, ok := h.Clients[notification.UserID]; ok {
				select {
				case client.Send <- notification.Message:
				default:
					close(client.Send)
					_ = client.Conn.Close()
					delete(h.Clients, notification.UserID)
				}
			}

		}
	}
}

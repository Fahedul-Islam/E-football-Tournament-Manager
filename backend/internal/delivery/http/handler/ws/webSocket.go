package ws

import (
	"net/http"
	"strconv"
	"tournament-manager/infra/ws"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	Hub *ws.Hub
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userIDStr, ok := r.Context().Value("user_id").(string)
	if !ok {
		conn.Close()
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		conn.Close()
		return
	}
	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.Hub.Register <- client
	defer func() {
		h.Hub.Unregister <- client
	}()
	go client.WritePump()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

}

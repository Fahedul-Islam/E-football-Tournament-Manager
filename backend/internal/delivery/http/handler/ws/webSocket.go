package ws

import (
	"net/http"
	"strconv"
	"tournament-manager/infra/ws"
	"tournament-manager/internal/delivery/http/middleware"

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

	userIDStr, ok := r.Context().Value(middleware.ContextKeyUserID).(string)
	if !ok {
		_ = conn.Close()
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		_ = conn.Close()
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

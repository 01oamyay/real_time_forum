package controller

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocket struct {
	sync.Mutex
	connections map[id]*websocket.Conn
}

type id uint

func newWS() *WebSocket {
	return &WebSocket{
		connections: make(map[id]*websocket.Conn),
		Mutex:       sync.Mutex{},
	}
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	userId := r.Context().Value(h.service.Keys.IDKey).(id)

	h.webSocket.Lock()
	h.webSocket.connections[userId] = conn
	h.webSocket.Unlock()

	go func() {
		defer func() {
			conn.Close()
			h.webSocket.Lock()
			delete(h.webSocket.connections, userId)
			h.webSocket.Unlock()
		}()

		// first sent to all other connections that the user is online
		for id, conn := range h.webSocket.connections {
			if id != userId {
				err := conn.WriteJSON(map[string]interface{}{
					"event": "user-online",
					"data":  userId,
				})
				if err != nil {
					h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
					return
				}
			}
		}

		// second: wait for messages
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
				return
			}
			// TODO save to db

			// broadcast to all other connections

			for id, conn := range h.webSocket.connections {
				if id != userId {
					err := conn.WriteJSON(map[string]interface{}{
						"event": "message",
						"data":  message,
					})
					if err != nil {
						h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
						return
					}
				}
			}
		}
	}()
}

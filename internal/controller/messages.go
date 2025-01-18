package controller

import (
	"net/http"
	"sync"

	"rlf/internal/entity"

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
					conn.WriteJSON(map[string]interface{}{
						"event": "error",
						"error": err.Error(),
					})
					return
				}
				return
			}
		}

		// second: wait for messages
		for {
			msg := entity.Message{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"error": err.Error(),
				})
				return
			}

			chat, err := h.service.Message.GetChatById(r.Context(), msg.ChatId)
			if err != nil {
				conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"error": err.Error(),
				})
				return
			}

			var receiver_id uint
			if chat.UserID == uint(userId) {
				receiver_id = chat.UserId1
			} else {
				receiver_id = chat.UserID
			}

			receiverConn, ok := h.webSocket.connections[id(receiver_id)]
			if !ok {
				conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"error": "you cannot send a message to offline users",
				})
				return
			}

			msg, _, err = h.service.CreateMessage(r.Context(), msg)
			if err != nil {
				conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"error": err.Error(),
				})
				return
			}

			err = receiverConn.WriteJSON(map[string]interface{}{
				"event": "msg",
				"data":  msg,
			})
			if err != nil {
				conn.WriteJSON(map[string]interface{}{
					"event": "error",
					"error": err.Error(),
				})
				return
			}
		}
	}()
}

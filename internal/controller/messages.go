package controller

import (
	"encoding/json"
	"fmt"
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

	userId := id(uint(r.Context().Value(h.service.Keys.IDKey).(int)))

	h.webSocket.Lock()
	h.webSocket.connections[userId] = conn
	h.webSocket.Unlock()

	conn.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("Connection closing with code: %d, text: %s\n", code, text)

		h.webSocket.Mutex.Lock()
		delete(h.webSocket.connections, userId)
		h.webSocket.Mutex.Unlock()

		// Notify all other connected users
		for id, conn := range h.webSocket.connections {
			if id != userId {
				err := conn.WriteJSON(map[string]interface{}{
					"event": "user-offline",
					"data":  userId,
				})
				if err != nil {
					fmt.Printf("Error notifying user %d: %v\n", id, err)
					continue
				}
			}
		}

		return nil
	})

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
}

func (h *Handler) GetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	contacts, status, err := h.service.GetContacts(r.Context())
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}

	for k := range h.webSocket.connections {
		fmt.Println(k)
	}

	for i := 0; i < len(contacts); i++ {
		user_id := id(contacts[i].UserID)
		_, ok := h.webSocket.connections[user_id]

		if ok {
			contacts[i].IsOnline = true
		} else {
			contacts[i].IsOnline = false
		}
	}

	if err := json.NewEncoder(w).Encode(contacts); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		h.webSocket.Lock()
		delete(h.webSocket.connections, userId)
		// Notify others that user is offline
		for id, conn := range h.webSocket.connections {
			if id != userId {
				conn.WriteJSON(map[string]interface{}{
					"event": "user-offline",
					"data":  userId,
				})
			}
		}
		h.webSocket.Unlock()
		return nil
	})

	// Notify others that user is online
	h.webSocket.Lock()
	for id, conn := range h.webSocket.connections {
		if id != userId {
			err := conn.WriteJSON(map[string]interface{}{
				"event": "user-online",
				"data":  userId,
			})
			if err != nil {
				log.Printf("Error notifying user %v: %v", id, err)
			}
		}
	}
	h.webSocket.Unlock()

	defer conn.Close()

	// Message handling loop
	for {
		var msg entity.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		chat, err := h.service.Message.GetChatById(r.Context(), msg.ChatId)
		if err != nil {
			conn.WriteJSON(map[string]interface{}{
				"event": "error",
				"error": err.Error(),
			})
			continue // Don't close connection, just skip this message
		}

		var receiver_id uint
		if chat.UserID == uint(userId) {
			receiver_id = chat.UserId1
		} else {
			receiver_id = chat.UserID
		}

		h.webSocket.Lock()
		receiverConn, ok := h.webSocket.connections[id(receiver_id)]
		h.webSocket.Unlock()

		if !ok {
			conn.WriteJSON(map[string]interface{}{
				"event": "error",
				"error": "recipient is offline",
			})
			continue
		}

		msg, _, err = h.service.CreateMessage(r.Context(), msg)
		if err != nil {
			conn.WriteJSON(map[string]interface{}{
				"event": "error",
				"error": err.Error(),
			})
			continue
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
			continue
		}

		err = conn.WriteJSON(map[string]interface{}{
			"event": "msg",
			"data":  msg,
		})
		if err != nil {
			conn.WriteJSON(map[string]interface{}{
				"event": "error",
				"error": err.Error(),
			})
			continue
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

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	strChatID := r.URL.Path[len("/api/chat/"):]
	id, err := strconv.Atoi(strChatID)
	if err != nil || id < 0 {
		h.errorHandler(w, r, http.StatusBadRequest, "invalid chat id")
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, "Invalid limit")
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.errorHandler(w, r, http.StatusBadRequest, "Invalid offset")
		return
	}

	chat, messages, status, err := h.service.Message.GetMessagesByChat(r.Context(), uint(id), limit, offset)
	if err != nil {
		h.errorHandler(w, r, status, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(entity.MsgEvent{
		Chat:     chat,
		Messages: messages,
	}); err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.errorHandler(w, r, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}

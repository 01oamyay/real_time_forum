package controller

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Printf("Connection closed for user %v with code %d: %s", userId, code, text)
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

	// defer func() {
	// 	conn.Close()
	// 	h.webSocket.Lock()
	// 	delete(h.webSocket.connections, userId)
	// 	// Notify others that user is offline
	// 	for id, conn := range h.webSocket.connections {
	// 		if id != userId {
	// 			conn.WriteJSON(map[string]interface{}{
	// 				"event": "user-offline",
	// 				"data":  userId,
	// 			})
	// 		}
	// 	}
	// 	h.webSocket.Unlock()
	// }()

	// const (
	// 	writeWait  = 10 * time.Second
	// 	pongWait   = 60 * time.Second
	// 	pingPeriod = (pongWait * 9) / 10
	// )

	// conn.SetReadDeadline(time.Now().Add(pongWait))
	// conn.SetPongHandler(func(appData string) error {
	// 	conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })

	// ticker := time.NewTicker(pingPeriod)
	// defer ticker.Stop()

	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

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

	// Message handling loop
	for {
		msg := entity.Message{}
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

package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

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
	connections map[*websocket.Conn]int
}

type id uint

func newWS() *WebSocket {
	return &WebSocket{
		connections: make(map[*websocket.Conn]int),
		Mutex:       sync.Mutex{},
	}
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer conn.Close()

	userId := r.Context().Value(h.service.Keys.IDKey).(int)

	h.webSocket.Lock()
	h.webSocket.connections[conn] = userId
	h.webSocket.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)

	conn.SetCloseHandler(func(code int, text string) error {
		h.webSocket.Lock()
		delete(h.webSocket.connections, conn)

		hasOtherConnections := false
		for _, id := range h.webSocket.connections {
			if id == userId {
				hasOtherConnections = true
				break
			}
		}

		// Notify others that user is offline
		if !hasOtherConnections {
			for conn := range h.webSocket.connections {
				conn.WriteJSON(map[string]interface{}{
					"event": "user-offline",
					"data":  userId,
				})
			}
		}
		h.webSocket.Unlock()
		wg.Done()
		return nil
	})

	// Notify others that user is online
	h.webSocket.Lock()
	for conn, id := range h.webSocket.connections {
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
	go func() {
		var typingTimer *time.Timer
		// defer wg.Done()
		for {

			var event struct {
				Type    string          `json:"event"`
				Payload json.RawMessage `json:"payload"`
			}

			if err := conn.ReadJSON(&event); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("error: %v", err)
				}
				break
			}

			switch event.Type {
			case "msg":
				{
					var msg entity.Message
					if err := json.Unmarshal(event.Payload, &msg); err != nil {
						log.Printf("error unmarshaling message: %v", err)
						continue
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

					var receiverConn *websocket.Conn
					h.webSocket.Lock()
					for conn, id := range h.webSocket.connections {
						if id == int(receiver_id) {
							receiverConn = conn
							break
						}
					}
					h.webSocket.Unlock()

					if receiverConn == nil {
						conn.WriteJSON(map[string]interface{}{
							"event": "msg-error",
							"error": "recipient is offline",
						})
						continue
					}

					msg, _, err = h.service.CreateMessage(r.Context(), msg)
					if err != nil {
						conn.WriteJSON(map[string]interface{}{
							"event": "msg-error",
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
							"event": "msg-error",
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
							"event": "msg-error",
							"error": err.Error(),
						})
						continue
					}
				}
			case "typing":
				{
					var typing entity.Typing
					if err := json.Unmarshal(event.Payload, &typing); err != nil {
						log.Printf("error unmarshaling message: %v", err)
						continue
					}
					chat, err := h.service.Message.GetChatById(r.Context(), uint(typing.ChatID))
					if err != nil {
						conn.WriteJSON(map[string]interface{}{
							"event": "error",
							"error": err.Error(),
						})
						continue
					}

					var receiver_id uint
					if chat.UserID == uint(userId) {
						receiver_id = chat.UserId1
					} else {
						receiver_id = chat.UserID
					}

					typing.UserID = int(userId)

					user, err := h.service.User.GetUserById(r.Context(), userId)
					if err != nil {
						conn.WriteJSON(map[string]interface{}{
							"event": "error",
							"error": err.Error(),
						})
					}

					typing.Nickname = user.NickName

					var receiverConn *websocket.Conn
					h.webSocket.Lock()
					for conn, id := range h.webSocket.connections {
						if id == int(receiver_id) {
							receiverConn = conn
							break
						}
					}
					h.webSocket.Unlock()

					if receiverConn == nil {
						continue
					}

					err = receiverConn.WriteJSON(map[string]interface{}{
						"event":  "typing",
						"typing": typing,
					})
					if err != nil {
						conn.WriteJSON(map[string]interface{}{
							"event": "error",
							"error": err.Error(),
						})
					}

					if typingTimer != nil {
						typingTimer.Stop()
					}

					typingTimer = time.AfterFunc(2*time.Second, func() {
						typing.IsTyping = false
						err = receiverConn.WriteJSON(map[string]interface{}{
							"event":  "typing",
							"typing": typing,
						})
						if err != nil {
							conn.WriteJSON(map[string]interface{}{
								"event": "error",
								"error": err.Error(),
							})
						}
					})
				}
			}
		}
	}()

	wg.Wait()
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

	// Get the user IDs of online users
	onlineUserIds := make([]uint, 0)
	h.webSocket.Lock()
	for _, id := range h.webSocket.connections {
		onlineUserIds = append(onlineUserIds, uint(id))
	}
	h.webSocket.Unlock()

	// Update the contacts with online status
	for i := range contacts {
		for _, onlineId := range onlineUserIds {
			if contacts[i].UserID == onlineId {
				contacts[i].IsOnline = true
				break
			}
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

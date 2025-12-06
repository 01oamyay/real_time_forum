package entity

import (
	"database/sql"
	"time"
)

// Message is persisted chat history returned over HTTP and WS.
type Message struct {
	ID        uint      `json:"id"`
	ChatId    uint      `json:"chat_id"`
	SenderId  uint      `json:"sender_id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Chat ties two users together and keeps metadata such as last message time.
type Chat struct {
	ID      uint         `json:"id"`
	UserID  uint         `json:"user_id"`
	UserId1 uint         `json:"user_id_1"`
	LastMsg sql.NullTime `json:"last_msg"`
}

// Contact is used by the contact list API.
type Contact struct {
	UserID     uint         `json:"user_id"`
	FirstName  string       `json:"firstName"`
	LastName   string       `json:"lastName"`
	Nickname   string       `json:"nickname"`
	IsOnline   bool         `json:"isOnline"`
	LastMsg    sql.NullTime `json:"last_msg"`
	HasLastMsg bool         `json:"has_last_msg"`
}

// MsgEvent is streamed to the client when it loads a chat.
type MsgEvent struct {
	Chat     Chat      `json:"chat"`
	Messages []Message `json:"messages"`
}

// Typing indicates whether a chat partner is typing.
type Typing struct {
	ChatID   int    `json:"chat_id"`
	IsTyping bool   `json:"is_typing"`
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

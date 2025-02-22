package entity

import (
	"database/sql"
	"time"
)

type Message struct {
	ID        uint      `json:"id"`
	ChatId    uint      `json:"chat_id"`
	SenderId  uint      `json:"sender_id"`
	Nickname  string    `json:"nickname"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Chat struct {
	ID      uint         `json:"id"`
	UserID  uint         `json:"user_id"`
	UserId1 uint         `json:"user_id_1"`
	LastMsg sql.NullTime `json:"last_msg"`
}

type Contact struct {
	UserID     uint         `json:"user_id"`
	FirstName  string       `json:"firstName"`
	LastName   string       `json:"lastName"`
	Nickname   string       `json:"nickname"`
	IsOnline   bool         `json:"isOnline"`
	LastMsg    sql.NullTime `json:"last_msg"`
	HasLastMsg bool         `json:"has_last_msg"`
}

type MsgEvent struct {
	Chat     Chat      `json:"chat"`
	Messages []Message `json:"messages"`
}

type Typing struct {
	ChatID   int    `json:"chat_id"`
	IsTyping bool   `json:"is_typing"`
	UserID   int    `json:"user_id"`
	Nickname string `json:"nickname"`
}

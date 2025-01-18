package entity

import "time"

type Message struct {
	ID        uint      `json:"id"`
	ChatId    uint      `json:"chat_id"`
	SenderId  uint      `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Chat struct {
	ID      uint `json:"id"`
	UserID  uint `json:"user_id"`
	UserId1 uint `json:"user_id_1"`
}

type Contact struct {
	UserID    uint   `json:"user_id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nickname  string `json:"nickname"`
	IsOnline  bool   `json:"isOnline"`
}
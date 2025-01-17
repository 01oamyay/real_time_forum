package repository

import (
	"context"
	"database/sql"
	"net/http"

	"rlf/internal/entity"
)

type MessagesRepository struct {
	db   *sql.DB
	Keys Keys
}

func newMessagesRepo(db *sql.DB, keys Keys) *MessagesRepository {
	return &MessagesRepository{
		db:   db,
		Keys: keys,
	}
}

func (r *MessagesRepository) GetMessagesByChat(ctx context.Context, chatId uint, limit, offset int) ([]entity.Message, int, error) {
	query := `
		SELECT (chat_id, sender_id, content, created_at) FROM message
		WHERE chat_id = ?
		ORDER BY created_at
		LIMIT ? OFFSET ?
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer prep.Close()

	rows, err := prep.QueryContext(ctx, chatId, limit, offset)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer rows.Close()

	messages := []entity.Message{}

	for rows.Next() {
		msg := entity.Message{}
		if err := rows.Scan(&msg.ChatId, &msg.SenderId, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return messages, http.StatusOK, nil
}

func (r *MessagesRepository) GetAllUserChats(ctx context.Context) ([]entity.Chat, int, error) {
	userId := ctx.Value(r.Keys.IDKey).(int)

	query := `
		SELECT (id, user_id, user_id_1) FROM chat
		WHERE user_id = ? OR user_id_1 = ?
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()
	rows, err := prep.QueryContext(ctx, userId, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()
	chats := []entity.Chat{}
	for rows.Next() {
		chat := entity.Chat{}
		if err = rows.Scan(&chat.ID, &chat.UserID, &chat.UserId1); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		chats = append(chats, chat)
	}
	if err = rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return chats, http.StatusOK, nil
}

func (r *MessagesRepository) CreateChat(ctx context.Context, second_user uint) (entity.Chat, int, error) {
	userId := ctx.Value(r.Keys.IDKey).(int)
	query := `
		INSERT INTO chat (user_id, user_id_1)
		VALUES (?, ?)
		RETURNING id, user_id, user_id_1
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.Chat{}, http.StatusInternalServerError, err
	}
	defer prep.Close()

	chat := entity.Chat{}
	if err = prep.QueryRowContext(ctx, userId, second_user).Scan(&chat.ID, &chat.UserID, &chat.UserId1); err != nil {
		return chat, http.StatusInternalServerError, err
	}

	return chat, http.StatusOK, nil
}

func (r *MessagesRepository) CreateMessage(ctx context.Context, chatId uint, text string) (entity.Message, int, error) {
	senderId := ctx.Value(r.Keys.IDKey).(int)

	query := `
		INSERT INTO message (chat_id, sender_id, content)
		VALUES (?, ?, ?)
		RETURNING id, chat_id, sender_id, content, created_at
	`

	msg := entity.Message{}

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, chatId, senderId, text).Scan(&msg.ID, &msg.ChatId, &msg.SenderId, &msg.Content, &msg.CreatedAt); err != nil {
		return msg, http.StatusInternalServerError, err
	}

	return msg, http.StatusOK, nil
}

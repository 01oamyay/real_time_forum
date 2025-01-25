package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		SELECT id, chat_id, sender_id, content, created_at FROM message
		WHERE chat_id = ?
		ORDER BY created_at DESC
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
		if err := rows.Scan(&msg.ID, &msg.ChatId, &msg.SenderId, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var nickname string
		err = r.db.QueryRowContext(ctx, "SELECT nickname FROM users WHERE id = ?", msg.SenderId).Scan(&nickname)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		msg.Nickname = nickname
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
		SELECT id, user_id, user_id_1, last_msg FROM chat
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
		if err = rows.Scan(&chat.ID, &chat.UserID, &chat.UserId1, &chat.LastMsg); err != nil {
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

func (r *MessagesRepository) ChatExist(ctx context.Context, second_user uint) (entity.Chat, int, error) {
	userId := ctx.Value(r.Keys.IDKey)
	query := `
		SELECT id, user_id, user_id_1 FROM chat
		WHERE (user_id = ? AND user_id_1 = ?) OR (user_id = ? AND user_id_1 = ?)
	`

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return entity.Chat{}, http.StatusInternalServerError, err
	}
	defer prep.Close()

	chat := entity.Chat{}
	row := prep.QueryRowContext(ctx, userId, second_user, second_user, userId)
	if err = row.Scan(&chat.ID, &chat.UserID, &chat.UserId1); err != nil {
		if err == sql.ErrNoRows {
			return chat, http.StatusNotFound, err
		}
		return chat, http.StatusInternalServerError, err
	}

	return chat, http.StatusOK, nil
}

func (r *MessagesRepository) ChatExistsById(ctx context.Context, chat_id uint) (bool, int, error) {
	exists := false
	query := `
	SELECT EXISTS (
    	SELECT 1 FROM chat 
    	WHERE id = ?
	)
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return exists, http.StatusInternalServerError, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, chat_id).Scan(&exists); err != nil {
		return exists, http.StatusInternalServerError, err
	}
	return exists, http.StatusOK, nil
}

func (r *MessagesRepository) CreateMessage(ctx context.Context, chatId uint, text string) (entity.Message, int, error) {
	senderId := ctx.Value(r.Keys.IDKey).(int)

	query := `
		INSERT INTO message(chat_id, sender_id, content)
		VALUES(?, ?, ?)
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

	// Get the nickname of the sender
	var nickname string
	err = r.db.QueryRowContext(ctx, "SELECT nickname FROM users WHERE id = ?", senderId).Scan(&nickname)
	if err != nil {
		return msg, http.StatusInternalServerError, err
	}

	// Update the message struct with the nickname
	msg.Nickname = nickname

	return msg, http.StatusOK, nil
}

func (r *MessagesRepository) GetChatById(ctx context.Context, chat_id uint) (entity.Chat, error) {
	query := `SELECT id, user_id, user_id_1 FROM chat WHERE id = ?`
	chat := entity.Chat{}
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return chat, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, chat_id).Scan(&chat.ID, &chat.UserID, &chat.UserId1); err != nil {
		fmt.Println("here", err)
		return chat, err
	}

	if chat.ID == 0 {
		return chat, errors.New("chat not found")
	}
	return chat, nil
}

func (r *MessagesRepository) GetContacts(ctx context.Context) ([]entity.Contact, int, error) {
	userId := ctx.Value(r.Keys.IDKey).(int)
	query := `SELECT id, nickname, firstName, lastName FROM users WHERE id != ?`

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	rows, err := prep.QueryContext(ctx, userId)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	contacts := []entity.Contact{}
	for rows.Next() {
		contact := entity.Contact{}
		if err := rows.Scan(&contact.UserID, &contact.Nickname, &contact.FirstName, &contact.LastName); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(contacts) == 0 {
		return nil, http.StatusNoContent, errors.New("no contact")
	}
	return contacts, http.StatusOK, nil
}

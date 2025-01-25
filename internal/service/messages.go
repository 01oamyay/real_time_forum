package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"rlf/internal/entity"
	"rlf/internal/repository"
)

type MessagesService struct {
	messagesRepo repository.Message
	usersRepo    repository.User
}

func newMessagesService(msgRepo repository.Message, userRepo repository.User) *MessagesService {
	return &MessagesService{
		messagesRepo: msgRepo,
		usersRepo:    userRepo,
	}
}

func (s *MessagesService) GetMessagesByChat(ctx context.Context, second_user uint, limit, offset int) (entity.Chat, []entity.Message, int, error) {
	chat, status, _ := s.messagesRepo.ChatExist(ctx, second_user)
	if status == http.StatusNotFound {
		chat, status, err := s.messagesRepo.CreateChat(ctx, second_user)
		if err != nil {
			return chat, nil, status, err
		}
		return chat, []entity.Message{}, http.StatusOK, nil
	}

	messages, status, err := s.messagesRepo.GetMessagesByChat(ctx, uint(chat.ID), limit, offset)
	return chat, messages, status, err
}

func (s *MessagesService) GetAllUserChats(ctx context.Context) ([]entity.Chat, int, error) {
	return s.messagesRepo.GetAllUserChats(ctx)
}

func (s *MessagesService) CreateChat(ctx context.Context, second_user uint) (entity.Chat, int, error) {
	chat := entity.Chat{}
	userExists, status, err := s.usersRepo.Exists(ctx, second_user)
	if err != nil {
		return chat, status, err
	}

	if !userExists {
		return chat, http.StatusBadRequest, errors.New("second user not found")
	}

	chat, status, err = s.messagesRepo.CreateChat(ctx, second_user)
	return chat, status, err
}

func (s *MessagesService) CreateMessage(ctx context.Context, msg entity.Message) (entity.Message, int, error) {
	exists, status, err := s.messagesRepo.ChatExistsById(ctx, msg.ChatId)
	if err != nil {
		return msg, status, err
	}

	if !exists {
		return msg, http.StatusBadRequest, errors.New("chat not found")
	}

	return s.messagesRepo.CreateMessage(ctx, msg.ChatId, msg.Content)
}

func (s *MessagesService) GetChatById(ctx context.Context, chat_id uint) (entity.Chat, error) {
	return s.messagesRepo.GetChatById(ctx, chat_id)
}

func (s *MessagesService) GetContacts(ctx context.Context) ([]entity.Contact, int, error) {
	allChats, status, err := s.messagesRepo.GetAllUserChats(ctx)
	if err != nil {
		return nil, status, err
	}

	contacts, status, err := s.messagesRepo.GetContacts(ctx)
	if err != nil {
		return nil, status, err
	}

	chatMap := make(map[uint]sql.NullTime)
	for _, chat := range allChats {
		// Check both potential user ID fields
		if chat.UserID != 0 {
			chatMap[chat.UserID] = chat.LastMsg
		}
		if chat.UserId1 != 0 {
			chatMap[chat.UserId1] = chat.LastMsg
		}

	}

	// Update contacts with last message time if found
	for i := range contacts {
		if lastMsg, exists := chatMap[contacts[i].UserID]; exists {
			contacts[i].LastMsg = lastMsg
			contacts[i].HasLastMsg = true
		}
	}

	return contacts, status, err
}

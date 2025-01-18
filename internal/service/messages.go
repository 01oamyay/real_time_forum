package service

import (
	"context"
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

func (s *MessagesService) GetMessagesByChat(ctx context.Context, second_user uint, limit, offset int) ([]entity.Message, int, error) {
	chatId, status, _ := s.messagesRepo.ChatExist(ctx, second_user)
	if status == http.StatusNotFound {
		_, status, err := s.messagesRepo.CreateChat(ctx, second_user)
		if err != nil {
			return nil, status, err
		}
		return []entity.Message{}, http.StatusNoContent, nil
	}

	return s.messagesRepo.GetMessagesByChat(ctx, uint(chatId), limit, offset)
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

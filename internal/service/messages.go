package service

import (
	"context"

	"rlf/internal/entity"
	"rlf/internal/repository"
)

type MessagesService struct {
	messagesRepo repository.Message
}

func newMessagesService(msgRepo repository.Message) *MessagesService {
	return &MessagesService{
		messagesRepo: msgRepo,
	}
}

func (s *MessagesService) GetMessagesByChat(ctx context.Context, chatId uint, limit, offset int) ([]entity.Message, int, error)

func (s *MessagesService) GetAllUserChats(ctx context.Context) ([]entity.Chat, int, error)

func (s *MessagesService) CreateChat(ctx context.Context, second_user uint) (entity.Chat, int, error)

func (s *MessagesService) CreateMessage(ctx context.Context, chatId uint, text string) (entity.Message, int, error)

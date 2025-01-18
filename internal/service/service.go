package service

import (
	"context"

	"rlf/internal/entity"
	"rlf/internal/repository"
)

type User interface {
	Create(ctx context.Context, user entity.User) (int, error)
	SignIn(ctx context.Context, user entity.User) (string, int, error)
}

type Session interface {
	IsTokenExist(ctx context.Context, token string) (bool, error)
	DeleteSessionByToken(ctx context.Context, token string) error
	DeleteSessionByUserID(ctx context.Context, userID uint) error
}

type Post interface {
	CreatePost(ctx context.Context, input entity.Post) (uint, int, error)
	UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error)
	GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error)
	GetAllByCategory(ctx context.Context, tagName string, limit, offset int) ([]entity.Post, int, error)
	GetAllByUserID(ctx context.Context, userID uint, limit, offset int) ([]entity.Post, int, error)
	GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool, limit, offset int) ([]entity.Post, int, error)
}

type Category interface {
	GetAllCategorys(ctx context.Context) ([]entity.Category, int, error)
}

type Comment interface {
	CreateComment(ctx context.Context, input entity.Comment) (int, error)
	UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error)
}

type Message interface {
	GetMessagesByChat(ctx context.Context, second_user uint, limit, offset int) ([]entity.Message, int, error)
	GetAllUserChats(ctx context.Context) ([]entity.Chat, int, error)
	CreateChat(ctx context.Context, second_user uint) (entity.Chat, int, error)
	CreateMessage(ctx context.Context, msg entity.Message) (entity.Message, int, error)
	GetChatById(ctx context.Context, chat_id uint) (entity.Chat, error)
}

type Service struct {
	User
	Session
	Post
	Comment
	Category
	Message
	repository.Keys
}

func NewService(repo *repository.Repository, secret string) *Service {
	return &Service{
		User:     newUserService(repo.User, repo.Session, secret),
		Session:  newSessionService(repo.Session),
		Post:     newPostService(repo.Post, repo.Category),
		Comment:  newCommentService(repo.Comment),
		Category: newCategoryService(repo.Category),
		Message:  newMessagesService(repo.Message, repo.User),
		Keys:     repo.Keys,
	}
}

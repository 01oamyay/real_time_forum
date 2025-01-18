package repository

import (
	"context"
	"database/sql"

	"rlf/internal/entity"
)

type User interface {
	Create(ctx context.Context, user entity.User) (int, error)
	GetUserIDByEmail(ctx context.Context, email string) (entity.User, int, error)
	Exists(ctx context.Context, userId uint) (bool, int, error)
}

type Session interface {
	IsTokenExist(ctx context.Context, token string) (bool, error)
	DeleteSessionByToken(ctx context.Context, token string) error
	PostSession(ctx context.Context, session entity.Session) (int, error)
	DeleteSessionByUserID(ctx context.Context, userID uint) error
}

type Post interface {
	CreatePost(ctx context.Context, input entity.Post) (uint, int, error)
	DeletePostByID(ctx context.Context, PostID uint, userID uint) (int, error)
	UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error)
	GetAllByCategory(ctx context.Context, categoryName string, limit, offset int) ([]entity.Post, int, error)
	GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error)
	GetAllByUserID(ctx context.Context, userID uint, limit, offset int) ([]entity.Post, int, error)
	GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool, limit, offset int) ([]entity.Post, int, error)
}

type Category interface {
	CreateCategorys(ctx context.Context, categorysName []string) (int, error)
	GetCategorysIDByName(ctx context.Context, categorysName []string) ([]uint, int, error)
	CreateCategorysAndPostCon(ctx context.Context, categorysID []uint, postID uint) (int, error)
	CategoryExist(ctx context.Context, categoryName string) (bool, int, error)
	GetAllCategories(ctx context.Context) ([]entity.Category, int, error)
}

type Comment interface {
	CreateComment(ctx context.Context, input entity.Comment) (int, error)
	UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error)
}

type Message interface {
	GetMessagesByChat(ctx context.Context, chatId uint, limit, offset int) ([]entity.Message, int, error)
	GetAllUserChats(ctx context.Context) ([]entity.Chat, int, error)
	CreateChat(ctx context.Context, second_user uint) (entity.Chat, int, error)
	ChatExist(ctx context.Context, second_user uint) (int, int, error)
	ChatExistsById(ctx context.Context, chat_id uint) (bool, int, error)
	CreateMessage(ctx context.Context, chatId uint, text string) (entity.Message, int, error)
	GetChatById(ctx context.Context, chat_id uint) (entity.Chat, error)

}

type key string

type Keys struct {
	IDKey    key
	TokenKey key
}

type Repository struct {
	Post
	User
	Session
	Category
	Comment
	Message
	Keys
}

func NewRepository(db *sql.DB) *Repository {
	Keys := Keys{
		IDKey:    "id",
		TokenKey: "token",
	}
	return &Repository{
		User:     newUserRepository(db),
		Session:  newSessionRepository(db),
		Post:     newPostRepository(db, Keys),
		Category: newCategoryRepository(db),
		Comment:  newCommentRepository(db),
		Message:  newMessagesRepo(db, Keys),
		Keys:     Keys,
	}
}

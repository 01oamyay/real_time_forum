package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"rlf/internal/entity"
	"rlf/internal/repository"
)

type CommentService struct {
	commentRepo repository.Comment
}

func newCommentService(commentRepo repository.Comment) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(ctx context.Context, input entity.Comment) (int, error) {
	if strings.TrimSpace(input.Data) == "" {
		return http.StatusBadRequest, errors.New("invalid data")
	} else if input.PostID == 0 {
		return http.StatusBadRequest, errors.New("invalid postID")
	}
	return s.commentRepo.CreateComment(ctx, input)
}

func (s *CommentService) UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error) {
	if input.Vote != 0 && input.Vote != 1 {
		return http.StatusBadRequest, errors.New("invalid vote")
	}
	return s.commentRepo.UpsertCommentVote(ctx, input)
}

package service

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"rlf/internal/entity"
	"rlf/internal/repository"
)

type PostService struct {
	postRepo     repository.Post
	categoryRepo repository.Category
}

func newPostService(postRepo repository.Post, categoryRepo repository.Category) *PostService {
	return &PostService{
		postRepo:     postRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *PostService) CreatePost(ctx context.Context, input entity.Post) (uint, int, error) {
	if input.Data == "" || len(input.Data) > 10000 {
		return 0, http.StatusBadRequest, errors.New("size of text must be beyween 1 and 10000")
	} else if input.Title == "" || len(input.Title) > 58 {
		return 0, http.StatusBadRequest, errors.New("size of text must be between 1 and  58")
	} else if len(input.Categorys) == 0 {
		return 0, http.StatusBadRequest, errors.New("categorys is empty")
	} else if len(input.Categorys) > 5 {
		return 0, http.StatusBadRequest, errors.New("max categories are 5")
	}

	for _, category := range input.Categorys {
		exist, status, err := s.categoryRepo.CategoryExist(ctx, category)
		if err != nil {
			return 0, status, err
		}

		if !exist {
			return 0, status, errors.New("invalid category")
		}
	}

	postID, status, err := s.postRepo.CreatePost(ctx, input)
	if err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}

	input.Categorys = append(input.Categorys, "General")

	CategoryIDS, status, err := s.categoryRepo.GetCategorysIDByName(ctx, input.Categorys)
	if err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	if status, err := s.categoryRepo.CreateCategorysAndPostCon(ctx, CategoryIDS, postID); err != nil {
		if _, Posterr := s.postRepo.DeletePostByID(ctx, postID, input.UserID); Posterr != nil {
			log.Println(Posterr)
		}
		return 0, status, err
	}
	return postID, http.StatusOK, nil
}

func (s *PostService) GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error) {
	return s.postRepo.GetPostByID(ctx, postID)
}

func (s *PostService) UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error) {
	if input.Vote != 0 && input.Vote != 1 {
		return http.StatusBadRequest, errors.New("invalid vote")
	}
	return s.postRepo.UpsertPostVote(ctx, input)
}

func (s *PostService) GetAllByCategory(ctx context.Context, categoryName string, limit, offset int) ([]entity.Post, int, error) {
	if strings.TrimSpace(categoryName) == "" {
		return nil, http.StatusBadRequest, errors.New("invalid Category")
	}
	return s.postRepo.GetAllByCategory(ctx, categoryName, limit, offset)
}

func (s *PostService) GetAllByUserID(ctx context.Context, userID uint, limit, offset int) ([]entity.Post, int, error) {
	return s.postRepo.GetAllByUserID(ctx, userID, limit, offset)
}

func (s *PostService) GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool, limit, offset int) ([]entity.Post, int, error) {
	return s.postRepo.GetAllLikedPostsByUserID(ctx, userID, islike, limit, offset)
}

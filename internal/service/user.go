package service

import (
	"context"
	"errors"
	"net/http"

	"rlf/internal/entity"
	"rlf/internal/repository"
	smpljwt "rlf/pkg/smplJwt"
	"rlf/pkg/utils"
)

type UserService struct {
	userRepo    repository.User
	sessionRepo repository.Session
	secret      string
}

var UserS UserService

func newUserService(userRepo repository.User, sessionRepo repository.Session, secret string) *UserService {
	return &UserService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		secret:      secret,
	}
}

func (s *UserService) Create(ctx context.Context, user entity.User) (int, error) {
	if err := utils.IsValidRegister(&user); err != nil {
		return http.StatusBadRequest, err
	}
	status, err := s.userRepo.Create(ctx, user)
	if status == http.StatusBadRequest {
		switch err.Error() {
		case "UNIQUE constraint failed: users.email":
			return status, errors.New("already email is using")
		case "UNIQUE constraint failed: users.nickname":
			return status, errors.New("already nickname is using")
		}
	}
	return status, err
}

func (s *UserService) GetUserById(ctx context.Context, id int) (entity.User, error) {
	return s.userRepo.GetUserById(ctx, id)
}

func (s *UserService) SignIn(ctx context.Context, user entity.UserInput) (string, int, error) {
	if user.Login == "" {
		return "", http.StatusBadRequest, errors.New("invalid credentials")
	} else if user.Password == "" {
		return "", http.StatusBadRequest, errors.New("invalid password")
	}

	repoUserStruct, status, err := s.userRepo.GetUserIDByLogin(ctx, user.Login)
	if err != nil {
		if status == http.StatusBadRequest {
			return "", status, errors.New("invalid email or password")
		}
		return "", status, err
	}
	if err := utils.CompareHashAndPassword(repoUserStruct.Password, user.Password); err != nil {
		return "", http.StatusBadRequest, errors.New("invalid password")
	}
	token, err := smpljwt.NewJWT(repoUserStruct.ID, s.secret)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if status, err = s.sessionRepo.PostSession(ctx, entity.Session{
		UserID: repoUserStruct.ID,
		Token:  token,
	}); err != nil {
		return "", status, err
	}
	return token, http.StatusOK, nil
}

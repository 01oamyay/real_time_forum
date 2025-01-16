package service

import (
	"context"

	"rlf/internal/repository"
)

type SessionService struct {
	sessionRepo repository.Session
}

func newSessionService(sessionRepo repository.Session) *SessionService {
	return &SessionService{sessionRepo: sessionRepo}
}

func (s *SessionService) IsTokenExist(ctx context.Context, token string) (bool, error) {
	return s.sessionRepo.IsTokenExist(ctx, token)
}

func (s *SessionService) DeleteSessionByToken(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteSessionByToken(ctx, token)
}

func (s *SessionService) DeleteSessionByUserID(ctx context.Context, userID uint) error {
	return s.sessionRepo.DeleteSessionByUserID(ctx, userID)
}

package repository

import (
	"context"
	"database/sql"
	"net/http"

	"rlf/internal/entity"
)

type SessionRepository struct {
	db *sql.DB
}

func newSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) IsTokenExist(ctx context.Context, token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sessions WHERE token = $1);`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer prep.Close()
	if err := prep.QueryRowContext(ctx, token).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SessionRepository) DeleteSessionByToken(ctx context.Context, token string) error {
	query := "DELETE FROM sessions WHERE token = $1;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, token); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) DeleteSessionByUserID(ctx context.Context, userID uint) error {
	query := "DELETE FROM sessions WHERE user_id = $1;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, userID); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) PostSession(ctx context.Context, session entity.Session) (int, error) {
	if err := r.DeleteSessionByUserID(ctx, session.UserID); err != nil {
		return http.StatusInternalServerError, err
	}
	query := `INSERT INTO sessions(token, user_id) VALUES($1, $2)`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, session.Token, session.UserID); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

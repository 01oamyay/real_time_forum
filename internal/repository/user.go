package repository

import (
	"context"
	"database/sql"
	"net/http"

	"rlf/internal/entity"
)

type UserRepository struct {
	db *sql.DB
}

func newUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user entity.User) (int, error) {
	query := `INSERT INTO users(username, email, hashPass)
	VALUES($1, $2, $3) RETURNING id;`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err = prep.ExecContext(ctx, user.Username, user.Email, user.Password); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusCreated, nil
}

func (r *UserRepository) GetUserIDByEmail(ctx context.Context, email string) (entity.User, int, error) {
	query := `SELECT id, hashPass FROM users WHERE email = $1 LIMIT 1;`
	prep, err := r.db.PrepareContext(ctx, query)
	user := entity.User{}
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, email).Scan(&user.ID, &user.Password); err != nil {
		return user, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

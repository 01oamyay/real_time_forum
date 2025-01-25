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
	query := `INSERT INTO users(nickname, email, age, gender, firstName, lastName, password)
	VALUES(?, ?, ?, ?, ?, ?, ?) RETURNING id;`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err = prep.ExecContext(ctx, user.NickName, user.Email, user.Age, user.Gender, user.FirstName, user.LastName, user.Password); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusCreated, nil
}

func (r *UserRepository) GetUserIDByLogin(ctx context.Context, login string) (entity.User, int, error) {
	query := `SELECT id, password FROM users WHERE email = ? OR nickname = ? LIMIT 1;`
	prep, err := r.db.PrepareContext(ctx, query)
	user := entity.User{}
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, login, login).Scan(&user.ID, &user.Password); err != nil {
		return user, http.StatusBadRequest, err
	}
	return user, http.StatusOK, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (entity.User, error) {
	query := `SELECT id, nickname FROM users WHERE id = ? LIMIT 1;`
	prep, err := r.db.PrepareContext(ctx, query)
	user := entity.User{}
	if err != nil {
		return user, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, id).Scan(&user.ID, &user.NickName); err != nil {
		return user, err
	}
	return user, nil
}

func (r *UserRepository) Exists(ctx context.Context, userId uint) (bool, int, error) {
	exists := false
	query := `
	SELECT EXISTS (
    	SELECT 1 FROM users 
    	WHERE id = ?
	)
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return exists, http.StatusInternalServerError, err
	}
	defer prep.Close()

	if err = prep.QueryRowContext(ctx, userId).Scan(&exists); err != nil {
		return exists, http.StatusInternalServerError, err
	}
	return exists, http.StatusOK, nil
}

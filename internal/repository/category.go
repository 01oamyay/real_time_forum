package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"rlf/internal/entity"
)

type CategoryRepository struct {
	db *sql.DB
}

func newCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]entity.Category, int, error) {
	query := "SELECT id, name, description FROM category;"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to query categories: %v", err)
	}
	defer rows.Close()

	var categories []entity.Category
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.CategoryID, &category.Name, &category.Description); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("failed to scan category row: %v", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error reading category rows: %v", err)
	}

	if len(categories) == 0 {
		return nil, http.StatusNotFound, errors.New("no categories")
	}

	return categories, http.StatusOK, nil
}

func (r *CategoryRepository) CreateCategorys(ctx context.Context, categorysName []string) (int, error) {
	query := "INSERT OR IGNORE INTO category (name) VALUES ($1);"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	for _, category := range categorysName {
		if _, err := prep.ExecContext(ctx, category); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

func (r *CategoryRepository) GetCategorysIDByName(ctx context.Context, categorysName []string) ([]uint, int, error) {
	ids := []uint{}
	query := "SELECT id FROM category WHERE name = $1;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	for _, category := range categorysName {
		var id uint
		if err = prep.QueryRowContext(ctx, category).Scan(&id); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		ids = append(ids, id)
	}

	return ids, http.StatusOK, nil
}

func (r *CategoryRepository) CategoryExist(ctx context.Context, categoryName string) (bool, int, error) {
	// Validate input
	if strings.TrimSpace(categoryName) == "" {
		return false, http.StatusBadRequest, errors.New("invalid category")
	}

	query := "SELECT EXISTS(SELECT 1 FROM category WHERE name = $1);"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, categoryName).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, http.StatusOK, nil
		}
		return false, http.StatusInternalServerError, fmt.Errorf("error checking category existence: %v", err)
	}

	return exists, http.StatusOK, nil
}

func (r *CategoryRepository) CreateCategorysAndPostCon(ctx context.Context, categorysID []uint, postID uint) (int, error) {
	query := "INSERT INTO Category_and_post(Category_id, post_id) VALUES($1, $2);"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	for _, categoryID := range categorysID {
		if _, err := prep.ExecContext(ctx, categoryID, postID); err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}

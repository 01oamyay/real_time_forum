package repository

import (
	"context"
	"database/sql"
	"net/http"

	"rlf/internal/entity"
)

type CommentRepository struct {
	db *sql.DB
}

func newCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(ctx context.Context, input entity.Comment) (int, error) {
	query := `INSERT INTO comment(user_id, post_id, data) VALUES($1, $2, $3)`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, input.UserID, input.PostID, input.Data); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (r *CommentRepository) UpsertCommentVote(ctx context.Context, input entity.CommentVote) (int, error) {
	query := "SELECT vote FROM comment_vote WHERE user_id = $1 and comment_id = $2;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	var vote int
	if err := prep.QueryRowContext(ctx, input.UserID, input.CommentID).Scan(&vote); err != nil {
		if err == sql.ErrNoRows {
			query = "INSERT INTO comment_vote(user_id, comment_id, vote) VALUES($1, $2, $3);"
			if _, err = r.db.ExecContext(ctx, query, input.UserID, input.CommentID, input.Vote); err != nil {
				return http.StatusBadRequest, err
			}
		} else {
			return http.StatusInternalServerError, err
		}
	} else {
		if vote == input.Vote {
			query = "DELETE FROM comment_vote WHERE user_id = $1 and comment_id = $2;"
			if _, err := r.db.ExecContext(ctx, query, input.UserID, input.CommentID); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			query = "UPDATE comment_vote SET vote = $1 WHERE user_id = $2 and comment_id = $3;"
			if _, err := r.db.ExecContext(ctx, query, input.Vote, input.UserID, input.CommentID); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}
	return http.StatusOK, nil
}

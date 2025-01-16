package repository

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"rlf/internal/entity"
)

type PostRepository struct {
	db   *sql.DB
	keys Keys
}

func newPostRepository(db *sql.DB, keys Keys) *PostRepository {
	return &PostRepository{db: db, keys: keys}
}

func (r *PostRepository) GetAllByCategory(ctx context.Context, categoryName string, limit, offset int) ([]entity.Post, int, error) {
	userId := ctx.Value(r.keys.IDKey).(int)

	query := `
SELECT
    p.id,
    p.user_id,
    p.title,
    p.data,
    u.username,	
    (SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 1) as likes,
    (SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 0) as dislikes,
    (SELECT COUNT(*) FROM comment WHERE post_id = p.id) as comments_count,
    CASE 
	    WHEN $1 = 0 THEN 0
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 1) THEN 1
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 0) THEN 2
        ELSE 0
    END as vote_status
FROM
    post p
    INNER JOIN category_and_post tp ON p.id = tp.post_id
    INNER JOIN users u ON u.id = p.user_id
    INNER JOIN category t ON tp.category_id = t.id
WHERE t.name LIKE($2) 
ORDER BY p.id DESC
LIMIT $3 OFFSET $4;
`

	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	posts := []entity.Post{}
	rows, err := prep.QueryContext(ctx, userId, categoryName, limit, offset)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	defer rows.Close()

	for rows.Next() {

		post := entity.Post{}
		if err := rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.Title,
			&post.Data,
			&post.UserName,
			&post.Likes,
			&post.Dislikes,
			&post.CommentsCount,
			&post.VoteStatus,
		); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		// Fetch categories
		categorys, status, err := r.getCategorysByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Categorys = categorys

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if len(posts) == 0 {
		return nil, http.StatusNoContent, fmt.Errorf("no posts found for category: %s", categoryName)
	}

	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetAllByUserID(ctx context.Context, userID uint, limit, offset int) ([]entity.Post, int, error) {
	posts := []entity.Post{}
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username,
	    (SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 1) as likes,
    	(SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 0) as dislikes,
    	(SELECT COUNT(*) FROM comment WHERE post_id = p.id) as comments_count,
	CASE 
	    WHEN $1 = 0 THEN 0
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 1) THEN 1
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 0) THEN 2
        ELSE 0
    END as vote_status
	FROM
		post p
		INNER JOIN users u ON u.id = p.user_id
	WHERE p.user_id = $1
	LIMIT $2 OFFSET $3;
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()
	rows, err := prep.QueryContext(ctx, userID, limit, offset)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName, &post.Likes, &post.Dislikes, &post.CommentsCount, &post.VoteStatus); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		categorys, status, err := r.getCategorysByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Categorys = categorys
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetAllLikedPostsByUserID(ctx context.Context, userID uint, islike bool, limit, offset int) ([]entity.Post, int, error) {
	posts := []entity.Post{}
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username,
		(SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 1) as likes,
    	(SELECT COUNT(*) FROM post_vote WHERE post_id = p.id AND vote = 0) as dislikes,
    	(SELECT COUNT(*) FROM comment WHERE post_id = p.id) as comments_count
	FROM
		post p
	INNER JOIN users u on p.user_id = u.id
	INNER JOIN post_vote pv ON p.id = pv.post_id
	WHERE
		pv.user_id = $1 AND pv.vote = %d
		LIMIT $2 OFFSET $3;
	`
	if islike {
		query = fmt.Sprintf(query, 1)
	} else {
		query = fmt.Sprintf(query, 0)
	}
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rows, err := prep.QueryContext(ctx, userID, limit, offset)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}
	for rows.Next() {
		post := entity.Post{}
		if err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName, &post.Likes, &post.Dislikes, &post.CommentsCount); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		categorys, status, err := r.getCategorysByPostID(ctx, post.PostID)
		if err != nil {
			return nil, status, err
		}
		post.Categorys = categorys
		if islike {
			post.VoteStatus = 1
		} else {
			post.VoteStatus = 2
		}
		posts = append(posts, post)
	}
	return posts, http.StatusOK, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, postID uint) (entity.Post, int, error) {
	var post entity.Post
	query := `
	SELECT
		p.id,
		p.user_id,
		p.title,
		p.data,
		u.username,
		COALESCE(COUNT(CASE WHEN pv.vote = 1 THEN 1 END), 0) AS voting,
		COALESCE(COUNT(CASE WHEN pv.vote = 0 THEN 1 END), 0) AS voting1,
	CASE 
	    WHEN $1 = 0 THEN 0
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 1) THEN 1
        WHEN EXISTS (SELECT 1 FROM post_vote WHERE post_id = p.id AND user_id = $1 AND vote = 0) THEN 2
        ELSE 0
    END AS vote_status
	FROM
		post p
		INNER JOIN users u ON p.user_id = u.id
		LEFT JOIN post_vote pv ON p.id = pv.post_id
	WHERE
		p.id = $2
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return post, http.StatusInternalServerError, err
	}

	userId := ctx.Value(r.keys.IDKey).(int)

	if err := prep.QueryRowContext(ctx, userId, postID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Data, &post.UserName, &post.Likes, &post.Dislikes, &post.VoteStatus); err != nil {
		return post, http.StatusNotFound, err
	}

	Categorys, status, err := r.getCategorysByPostID(ctx, postID)
	if err != nil {
		return post, status, err
	}
	post.Categorys = Categorys

	comments, status, err := r.getCommentsByPostID(ctx, postID)
	if err != nil {
		return post, status, err
	}
	post.Comments = comments
	return post, http.StatusOK, nil
}

func (r *PostRepository) getCommentsByPostID(ctx context.Context, postID uint) ([]entity.Comment, int, error) {
	query := `
	SELECT 
		c.id,
		c.user_id,
		c.data,
		u.username,
		COALESCE(COUNT(CASE WHEN cv.vote = 1 THEN 1 END), 0) AS voting,
		COALESCE(COUNT(CASE WHEN cv.vote = 0 THEN 1 END), 0) AS voting1,
	CASE 
        WHEN $1 = 0 THEN 0
        WHEN EXISTS (SELECT 1 FROM comment_vote WHERE comment_id = c.id AND user_id = $1 AND vote = 1) THEN 1
        WHEN EXISTS (SELECT 1 FROM comment_vote WHERE comment_id = c.id AND user_id = $1 AND vote = 0) THEN 2
        ELSE 0
    END as vote_status
	FROM 
		comment c
		INNER JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_vote cv ON c.id = cv.comment_id
	WHERE 
		c.post_id = $2
	GROUP BY
		c.id, c.user_id, c.data, u.username;
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer prep.Close()

	userId := ctx.Value(r.keys.IDKey)
	if userId == nil {
		userId = 0
	}
	rows, err := prep.QueryContext(ctx, userId, postID)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	comments := []entity.Comment{}
	for rows.Next() {
		comment := entity.Comment{}
		if err := rows.Scan(&comment.CommentID, &comment.UserID, &comment.Data, &comment.UserName, &comment.Likes, &comment.Dislikes, &comment.VoteStatus); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		comment.PostID = postID
		comments = append(comments, comment)
	}
	return comments, http.StatusOK, nil
}

func (r *PostRepository) getCategorysByPostID(ctx context.Context, postID uint) ([]string, int, error) {
	query := `
	SELECT
		t.name
	FROM 
		category t
		INNER JOIN category_and_post tp ON tp.post_id = $1 and t.id = tp.category_id
	`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	rows, err := prep.QueryContext(ctx, postID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	categorys := []string{}
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, http.StatusBadRequest, err
		}
		categorys = append(categorys, category)
	}

	return categorys, http.StatusOK, nil
}

func (r *PostRepository) CreatePost(ctx context.Context, input entity.Post) (uint, int, error) {
	query := `INSERT INTO post(user_id, title, data) VALUES($1, $2, $3) RETURNING id;`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, http.StatusInternalServerError, err
	}
	defer prep.Close()
	var id uint
	if err = prep.QueryRowContext(ctx, input.UserID, input.Title, input.Data).Scan(&id); err != nil {
		return 0, http.StatusBadRequest, err
	}
	return id, http.StatusOK, nil
}

func (r *PostRepository) DeletePostByID(ctx context.Context, PostID uint, userID uint) (int, error) {
	query := `DELETE FROM post WHERE id = $1 AND user_id = $2`
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	if _, err := prep.ExecContext(ctx, PostID, userID); err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func (r *PostRepository) UpsertPostVote(ctx context.Context, input entity.PostVote) (int, error) {
	query := "SELECT vote FROM post_vote WHERE user_id = $1 and post_id = $2;"
	prep, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer prep.Close()
	var vote int
	if err := prep.QueryRowContext(ctx, input.UserID, input.PostID).Scan(&vote); err != nil {
		if err == sql.ErrNoRows {
			query = "INSERT INTO post_vote(user_id, post_id, vote) VALUES($1, $2, $3);"
			if _, err = r.db.ExecContext(ctx, query, input.UserID, input.PostID, input.Vote); err != nil {
				return http.StatusBadRequest, err
			}
		} else {
			return http.StatusInternalServerError, err
		}
	} else {
		if vote == input.Vote {
			query = "DELETE FROM post_vote WHERE user_id = $1 and post_id = $2;"
			if _, err := r.db.ExecContext(ctx, query, input.UserID, input.PostID); err != nil {
				return http.StatusInternalServerError, err
			}
		} else {
			query = "UPDATE post_vote SET vote = $1 WHERE user_id = $2 and post_id = $3;"
			if _, err := r.db.ExecContext(ctx, query, input.Vote, input.UserID, input.PostID); err != nil {
				return http.StatusInternalServerError, err
			}
		}
	}
	return http.StatusOK, nil
}

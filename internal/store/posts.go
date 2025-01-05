package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	return err
}

func (s *PostStore) GetById(ctx context.Context, postId int) (*Post, error) {
	var post *Post = new(Post)

	query := `
		SELECT * FROM posts WHERE id = $1;
	`

	var created_at time.Time
	var updated_at time.Time

	err := s.db.QueryRowContext(
		ctx,
		query,
		postId,
	).Scan(
		&post.ID,
		&post.Title,
		&post.UserID,
		&post.Content,
		&created_at,
		pq.Array(&post.Tags),
		&updated_at,
	)

	post.CreatedAt = created_at.Format(time.RFC3339)
	post.UpdatedAt = updated_at.Format(time.RFC3339)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) Delete(c context.Context, postID int) error {
	query := `
		DELETE FROM posts WHERE id = $1 RETURNING id, title, content, user_id;
	`

	res, err := s.db.ExecContext(
		c,
		query,
		postID,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostStore) Patch(c context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE id = $3;
	`

	_, err := s.db.ExecContext(c, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	return nil
}

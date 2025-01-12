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
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments"`
	User      User      `json:"user"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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
		&post.Version,
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

	c, cancel := context.WithTimeout(c, QueryTimeoutDuration)
	defer cancel()

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
		SET title = $1, content = $2, version = version+1
		WHERE id = $3 AND version = $4
		RETURNING version;
	`

	c, cancel := context.WithTimeout(c, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(c, query, post.Title, post.Content, post.ID, post.Version).Scan(&post.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetaData, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $1 OR p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetaData

	for rows.Next() {
		var post PostWithMetaData

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCount,
		)

		feed = append(feed, post)

		if err != nil {
			return nil, err
		}
	}

	return feed, nil
}

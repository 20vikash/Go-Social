package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int) (*Post, error)
		Delete(context.Context, int) error
		Patch(context.Context, *Post) error
		GetUserFeed(context.Context, int64) ([]PostWithMetaData, error)
	}

	Users interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, int) (*User, error)
	}

	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(ctx context.Context, postID int64) ([]Comment, error)
	}

	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}

package store

import (
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

// UserStore enables DB operations for user related data
type UserStore interface {

	// GetByEmail returns a user registered with the given email
	GetByEmail(ctx *krogo.Context, email string) (*model.User, error)

	// Create creates a new user in DB
	Create(ctx *krogo.Context, user *model.User) error
}

// BookStore enables DB operations for book related data
type BookStore interface {
	Get(ctx *krogo.Context, page *model.Page, filter string, value string) ([]model.BookRes, error)
	GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error)
	Create(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error)
	Update(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error)
	Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error
}

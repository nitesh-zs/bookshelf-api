package service

import (
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

// UserSvc provides functions for user operations
type UserSvc interface {
	// Exists checks if a user with given email is registered
	Exists(ctx *krogo.Context, email string) (bool, error)

	// Create registers a new user
	Create(ctx *krogo.Context, user *model.User) error
}

type BookSvc interface {
	Get(ctx *krogo.Context, page *model.Page, filter string, value string) ([]model.BookRes, error)
	GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error)
	Create(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error)
	Update(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error)
	Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error
}

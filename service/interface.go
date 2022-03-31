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

	// IsAdmin checks if the user registered with given email is admin or not
	IsAdmin(ctx *krogo.Context, email string) (bool, error)

	// Create registers a new user
	Create(ctx *krogo.Context, user *model.User) error
}

type BookSvc interface {
	Get(ctx *krogo.Context, page *model.Page, filters *model.Filters) ([]model.BookRes, error)
	GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error)

	// Create creates a new book if the user has access to do so
	Create(ctx *krogo.Context, book *model.Book, user *model.User) (*model.BookRes, error)

	// Update modifies a book if the user has access to do so
	Update(ctx *krogo.Context, book *model.Book, user *model.User) (*model.BookRes, error)

	// Delete deletes a book if the user has access to do so
	Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error
}

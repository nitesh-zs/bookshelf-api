package store

import (
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

// UserStore enables DB operations for user related data
type UserStore interface {

	// Exists checks if a user with given ID exists in DB
	Exists(ctx *krogo.Context, email string) (bool, error)

	// Create creates a new user in DB
	Create(ctx *krogo.Context, user *model.User) error
}

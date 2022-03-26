package service

import (
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

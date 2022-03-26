package user

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"

	"github.com/nitesh-zs/bookshelf-api/model"
)

type store struct{}

//nolint:revive //store should not be exported
func New() store {
	return store{}
}

func (s store) Exists(ctx *krogo.Context, email string) (bool, error) {
	var id string

	row := ctx.DB().QueryRow(getUserID, email)
	err := row.Scan(&id)

	if err == sql.ErrNoRows {
		return false, errors.EntityNotFound{Entity: "user", ID: email}
	}

	if err != nil {
		return false, errors.DB{Err: err}
	}

	return true, nil
}

func (s store) Create(ctx *krogo.Context, user *model.User) error {
	_, err := ctx.DB().Exec(createUser, uuid.NewString(), user.Email, user.Name, user.Type)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

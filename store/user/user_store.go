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

func (s store) GetByEmail(ctx *krogo.Context, email string) (*model.User, error) {
	var user model.User
	var id string

	row := ctx.DB().QueryRow(getUserByEmail, email)
	err := row.Scan(&id, &user.Email, &user.Name, &user.Type)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "user", ID: email}
	}

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	user.ID = uid

	return &user, nil
}

func (s store) Create(ctx *krogo.Context, user *model.User) error {
	_, err := ctx.DB().Exec(createUser, uuid.NewString(), user.Email, user.Name, user.Type)
	if err != nil {
		return errors.DB{Err: err}
	}

	return nil
}

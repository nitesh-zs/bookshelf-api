package user

import (
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
	"github.com/nitesh-zs/bookshelf-api/store"
)

type svc struct {
	store store.UserStore
}

//nolint:revive //svc should not be exported
func New(s store.UserStore) svc {
	return svc{s}
}

func (s svc) Exists(ctx *krogo.Context, email string) (bool, error) {
	_, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s svc) Create(ctx *krogo.Context, user *model.User) error {
	if user.Email == ctx.Config.Get("ADMIN") {
		user.Type = "admin"
	} else {
		user.Type = "general"
	}

	return s.store.Create(ctx, user)
}

func (s svc) IsAdmin(ctx *krogo.Context, email string) (bool, error) {
	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if user.Type == "admin" {
		return true, nil
	}

	return false, nil
}

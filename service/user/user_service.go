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
	return s.store.Exists(ctx, email)
}

func (s svc) Create(ctx *krogo.Context, user *model.User) error {
	if user.Email == "nitesh.saxena@zopsmart.com" {
		user.Type = "admin"
	} else {
		user.Type = "general"
	}

	return s.store.Create(ctx, user)
}

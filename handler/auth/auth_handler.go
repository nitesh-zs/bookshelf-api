package auth

import (
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/service"
	"github.com/nitesh-zs/bookshelf-api/util"
)

type handler struct {
	svc service.UserSvc
}

//nolint:revive //handler should not be exported
func New(s service.UserSvc) handler {
	return handler{s}
}

func (h handler) Login(ctx *krogo.Context) (interface{}, error) {
	// get user data
	user, err := util.GetTokenData(ctx)
	if err != nil {
		return nil, errors.Unauthenticated{}
	}

	ok, _ := h.svc.Exists(ctx, user.Email)
	// if user not exists, create one
	if !ok {
		err := h.svc.Create(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return "success", nil
}

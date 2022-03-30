package book

import (
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/service"
)

type handler struct {
	svc service.BookSvc
}

//nolint:revive //handler should not be exported
func New(s service.BookSvc) handler {
	return handler{s}
}

func (h handler) Get(ctx *krogo.Context) (interface{}, error) {
	return nil, nil
}

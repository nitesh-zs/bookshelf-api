package book

import (
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/service"
	"github.com/nitesh-zs/bookshelf-api/util"
)

type handler struct {
	svc service.BookSvc
}

//nolint:revive //handler should not be exported
func New(s service.BookSvc) handler {
	return handler{s}
}

func (h handler) Get(ctx *krogo.Context) (interface{}, error) {
	filter := ctx.Param("filter")
	value := ctx.Param("keyword")

	page, err := util.Pagination(ctx)
	if err != nil {
		return nil, err
	}

	books, err := h.svc.Get(ctx, &page, filter, value)

	return books, err
}

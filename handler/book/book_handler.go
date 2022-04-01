package book

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"

	"github.com/nitesh-zs/bookshelf-api/model"
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

func (h handler) Delete(ctx *krogo.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	err = h.svc.Delete(ctx, uid, &model.User{})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h handler) Create(ctx *krogo.Context) (interface{}, error) {
	var book model.Book
	if err := ctx.Bind(&book); err != nil {
		fmt.Println("We could not bind the data")
	}

	book.ID = uuid.New()
	resp, err := h.svc.Create(ctx, &book, &model.User{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h handler) Update(ctx *krogo.Context) (interface{}, error) {
	id := ctx.PathParam("id")
	uid, err := uuid.Parse(id)

	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	var book model.Book
	if err2 := ctx.Bind(&book); err2 != nil {
		fmt.Println("We could not bind the data")
	}

	book.ID = uid
	resp, err := h.svc.Update(ctx, &book, &model.User{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

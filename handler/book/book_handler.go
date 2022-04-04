package book

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"

	"github.com/nitesh-zs/bookshelf-api/model"
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
	page, err := util.Pagination(ctx)
	if err != nil {
		return nil, err
	}

	err = page.Check()
	if err != nil {
		return nil, err
	}

	filters := &model.Filters{}

	if ctx.Param("genre") != "" {
		filters.Genre = ctx.Param("genre")
	}

	if ctx.Param("author") != "" {
		filters.Author = ctx.Param("author")
	}

	if ctx.Param("year") != "" {
		year, e := strconv.Atoi(ctx.Param("year"))
		if e != nil {
			return nil, errors.InvalidParam{Param: []string{"year"}}
		}

		filters.Year = year
	}

	if ctx.Param("language") != "" {
		filters.Language = ctx.Param("language")
	}

	books, err := h.svc.Get(ctx, page, filters)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (h handler) GetByID(ctx *krogo.Context) (interface{}, error) {
	id := ctx.PathParam("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.InvalidParam{Param: []string{"id"}}
	}

	return h.svc.GetByID(ctx, uid)
}

func (h handler) GetFilters(ctx *krogo.Context) (interface{}, error) {
	filter := ctx.PathParam("param")
	return h.svc.GetFilters(ctx, filter)
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
		return nil, errors.InvalidParam{Param: []string{"body"}}
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
		ctx.Logger.Error("We could not bind the data")
		return nil, errors.InvalidParam{Param: []string{"body"}}
	}

	book.ID = uid
	resp, err := h.svc.Update(ctx, &book, &model.User{})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

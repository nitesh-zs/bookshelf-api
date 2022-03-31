package book

import (
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
	"github.com/nitesh-zs/bookshelf-api/store"
)

type svc struct {
	store store.BookStore
}

//nolint:revive //svc should not be exported
func New(s store.BookStore) svc {
	return svc{s}
}

func (s svc) Get(ctx *krogo.Context, page *model.Page, filters *model.Filters) ([]model.BookRes, error) {
	books, err := s.store.Get(ctx, page, filters)
	return books, err
}

func (s svc) GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error) {
	return nil, nil
}

func (s svc) Create(ctx *krogo.Context, book *model.Book, user *model.User) (*model.BookRes, error) {
	//if user.Type == "admin" {
	if book == nil {
		return nil, errors.Error("No object to create")
	}
	resp, err := s.store.Create(ctx, book)
	if err != nil {
		return nil, errors.DB{}
	}
	return resp, nil
	//}
	//return nil, errors.Unauthorized{}
}

func (s svc) Update(ctx *krogo.Context, book *model.Book, user *model.User) (*model.BookRes, error) {
	if book == nil {
		return nil, errors.Error("No object to update")
	}
	resp, err := s.store.Update(ctx, book)
	if err != nil {
		return nil, errors.DB{}
	}
	return resp, nil
}

func (s svc) Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error {
	//if user.Type == "admin" {
	if id == uuid.Nil {
		return errors.Error("invalid uuid")
	}
	err := s.store.Delete(ctx, id)
	if err != nil {
		return errors.DB{}
	}
	return nil
	//}
	//return errors.Unauthorized{}
}

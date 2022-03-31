package book

import (
	"github.com/google/uuid"
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
	return s.store.Get(ctx, page, filters)
}

func (s svc) GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error) {
	return s.store.GetByID(ctx, id)
}

func (s svc) GetFilters(ctx *krogo.Context, filter string) ([]string, error) {
	return s.store.GetFilters(ctx, filter)
}

func (s svc) Create(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s svc) Update(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s svc) Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error {
	return nil
}

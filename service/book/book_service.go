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

func (s svc) Get(ctx *krogo.Context, page *model.Page, filter string, value string) ([]model.BookRes, error) {
	books, err := s.store.Get(ctx, page, filter, value)
	return books, err
}

func (s svc) GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error) {
	return nil, nil
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

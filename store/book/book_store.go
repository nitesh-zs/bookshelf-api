package book

import (
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/model"
)

type store struct{}

//nolint:revive //store should not be exported
func New() store {
	return store{}
}

func (s store) Get(ctx *krogo.Context, page *model.Page, filter string, value string) ([]model.BookRes, error) {
	books := make([]model.BookRes, 0)

	query, err := s.getQueryBuilder(filter, value)
	if err != nil {
		return nil, err
	}

	rows, err := ctx.DB().Query(query, page.Offset, page.Size)
	if err != nil {
		return nil, errors.DB{Err: err}
	}

	defer func() {
		rows.Close()

		err = rows.Err()
		if err != nil {
			ctx.Logger.Error(err)
		}
	}()

	for rows.Next() {
		var id, genre, author, title, imageURI string
		var book model.BookRes

		err := rows.Scan(&id, &genre, &author, &title, &imageURI)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		book.ID = uid
		book.Genre = genre
		book.Author = author
		book.Title = title
		book.ImageURI = imageURI

		books = append(books, book)
	}

	return books, nil
}

func (s store) GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error) {
	return nil, nil
}

func (s store) Create(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s store) Update(ctx *krogo.Context, book *model.Book, user *model.User) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s store) Delete(ctx *krogo.Context, id uuid.UUID, user *model.User) error {
	return nil
}

func (s store) getQueryBuilder(filter string, value string) (string, error) {
	query := `select id, genre, author, title, image_uri from book`

	if filter != "" && value != "" {
		query += `where ` + filter + `= '` + value + `' offset $1 limit $2;`
		return query, nil
	}

	if filter == "" && value == "" {
		query += ` offset $1 limit $2;`
		return query, nil
	}

	if filter == "" {
		return "", errors.MissingParam{Param: []string{"filter"}}
	}

	return "", errors.MissingParam{Param: []string{"keyword"}}
}

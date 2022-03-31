package book

import (
	"database/sql"
	"strconv"

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

func (s store) Get(ctx *krogo.Context, page *model.Page, filters *model.Filters) ([]model.BookRes, error) {
	books := make([]model.BookRes, 0)

	query := s.getQueryBuilder(filters)

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
		var (
			id   string
			book model.BookRes
		)

		err := rows.Scan(&id, &book.Title, &book.Author, &book.Genre, &book.ImageURI)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, errors.DB{Err: err}
		}

		book.ID = uid

		books = append(books, book)
	}

	return books, nil
}

func (s store) GetByID(ctx *krogo.Context, id uuid.UUID) (*model.BookRes, error) {
	book := &model.BookRes{}

	row := ctx.DB().QueryRow(getByID, id)
	err := row.Scan(&book.Title, &book.Author, &book.Summary, &book.Genre, &book.Year, &book.Publisher, &book.ImageURI)

	if err == sql.ErrNoRows {
		return nil, errors.EntityNotFound{Entity: "book", ID: id.String()}
	}

	if err != nil {
		return nil, errors.DB{Err: err}
	}

	book.ID = id

	return book, nil
}

func (s store) Create(ctx *krogo.Context, book *model.Book) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s store) Update(ctx *krogo.Context, book *model.Book) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func (s store) Delete(ctx *krogo.Context, id uuid.UUID) error {
	return nil
}

func (s store) getQueryBuilder(f *model.Filters) string {
	query := `select id, title, author, genre, image_uri from book`
	whereClause := ""

	if f.Author != "" {
		whereClause += ` author = '` + f.Author + `' AND`
	}

	if f.Language != "" {
		whereClause += ` language = '` + f.Language + `' AND`
	}

	if f.Genre != "" {
		whereClause += ` genre = '` + f.Genre + `' AND`
	}

	if f.Year != 0 {
		year := strconv.Itoa(f.Year)
		whereClause += ` year = ` + year + ` AND`
	}

	if len(whereClause) > 0 {
		whereClause = whereClause[:len(whereClause)-4]
	}

	if whereClause != "" {
		query += " where" + whereClause
	}

	query += ` offset $1 limit $2;`

	return query
}

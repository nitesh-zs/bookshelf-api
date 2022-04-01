package book

import (
	"database/sql"
	"fmt"
	"log"
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

func (s store) Create(ctx *krogo.Context, book *model.Book) (*model.BookRes, error) {
	if book == nil {
		return nil, errors.Error("No object to create")
	}
	_, err := ctx.DB().Exec(createBook, book.ID.String(), book.Title, book.Author, book.Summary, book.Genre, book.Year, book.RegNum, book.Publisher, book.Language, book.ImageURI)
	if err != nil {
		return nil, errors.DB{Err: errors.Error("cannot create object")}
	}
	var bookRes1 = &model.BookRes{
		ID:        book.ID,
		Title:     book.Title,
		Author:    book.Author,
		Summary:   book.Summary,
		Genre:     book.Genre,
		Year:      book.Year,
		Publisher: book.Publisher,
		ImageURI:  book.ImageURI,
	}
	return bookRes1, nil
}

func (s store) Update(ctx *krogo.Context, book *model.Book) (*model.BookRes, error) {
	if book == nil {
		return nil, errors.Error("No object to update")
	}
	query := getUpdateQuery(book)
	_, err := ctx.DB().Exec(query)
	if err != nil {
		log.Println("store err 1: ", err)
		return nil, errors.DB{Err: err}
	}
	response, err := s.GetByID(ctx, book.ID)
	if err != nil {
		log.Println("store err 2: ", err)
		return nil, errors.DB{Err: errors.Error("updated, but can't fetch")}
	}
	return response, nil
}

func (s store) Delete(ctx *krogo.Context, id uuid.UUID) error {
	_, err := ctx.DB().Exec(`DELETE FROM book WHERE id=$1 `, id.String())
	if err != nil {
		return errors.DB{Err: err}
	}
	return nil
}

func getUpdateQuery(book *model.Book) string {
	/*
		update book set
		id='18f8fb8a-689e-45f4-bcd8-4224f07998e2',
		title='Rashmirathi',
		author='Dinkar',
		summary='About something i dont know',
		genre='History',
		year=1985,
		reg_num='Mohitbabazindabad',
		publisher='Rajkamal Prakashan',
		language='Hindi',
		image_uri=''
		where id='18f8fb8a-689e-45f4-bcd8-4224f07998e2'
	*/
	query := "update book set"
	query += " title=" + "'" + book.Title + "'"
	query += ", author=" + "'" + book.Author + "'"
	query += ", summary=" + "'" + book.Summary + "'"
	query += ", genre=" + "'" + book.Genre + "'"
	query += ", year=" + "'" + fmt.Sprint(book.Year) + "'"
	query += ", reg_num=" + "'" + book.RegNum + "'"
	query += ", publisher=" + "'" + book.Publisher + "'"
	query += ", language=" + "'" + book.Language + "'"
	query += ", image_uri=" + "'" + book.ImageURI + "'"
	query += " where id=" + "'" + book.ID.String() + "'"
	return query
}

func (s store) getQueryBuilder(f *model.Filters) string {
	query := `select id, genre, author, title, image_uri from book`
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
		whereClause += ` year = '` + year + `' AND`
	}

	if len(whereClause) > 0 {
		whereClause = whereClause[:len(whereClause)-4]
	}

	if whereClause != "" {
		query += " WHERE" + whereClause
	}

	query += ` OFFSET $1 LIMIT $2`

	return query
}

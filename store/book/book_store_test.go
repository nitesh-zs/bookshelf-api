package book

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmizerany/assert"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/krogertechnology/krogo/pkg/datastore"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/krogertechnology/krogo/pkg/krogo/config"
	"github.com/krogertechnology/krogo/pkg/log"

	"github.com/nitesh-zs/bookshelf-api/model"
)

func getNewBook(id uuid.UUID) *model.Book {
	return &model.Book{
		ID:        id,
		Title:     "Madhushala",
		Author:    "Harivansh Rai Bachchan",
		Summary:   "This is short summary",
		Genre:     "Poetry",
		Year:      1997,
		RegNum:    "ISB8726W821",
		Publisher: "Rajpal Publishing",
		Language:  "Hindi",
		ImageURI:  "https://images-na.ssl-images-amazon.com/images/I/71Hc0nX3UHL.jpg",
	}
}

func getNewBookRes(id uuid.UUID) *model.BookRes {
	return &model.BookRes{
		ID:        id,
		Title:     "Madhushala",
		Author:    "Harivansh Rai Bachchan",
		Summary:   "This is short summary",
		Genre:     "Poetry",
		Year:      1997,
		Publisher: "Rajpal Publishing",
		ImageURI:  "https://images-na.ssl-images-amazon.com/images/I/71Hc0nX3UHL.jpg",
	}
}

func initializeTest(t *testing.T) (sqlmock.Sqlmock, *krogo.Context, store) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	if err != nil {
		t.Fatalf("error in creating mockDB: %v", err)
	}

	c := config.NewGoDotEnvProvider(log.NewLogger(), "../../../configs")
	k := krogo.NewWithConfig(c)
	k.ORM = nil
	gormDB, _ := gorm.Open("postgres", db)

	k.SetORM(datastore.GORMClient{DB: gormDB})

	ctx := krogo.NewContext(nil, nil, k)
	s := New()

	return mock, ctx, s
}

func filter() *model.Filters {
	return &model.Filters{
		Genre: "Self Help",
		Year:  2000,
	}
}

func bookRes() []model.BookRes {
	return []model.BookRes{
		{
			ID:       uuid.Nil,
			Title:    "Meditations",
			Author:   "Marcus Aurelius",
			Genre:    "Self Help",
			ImageURI: "image.com/woo-hoo",
		},
	}
}

func bookRes2() *model.BookRes {
	return &model.BookRes{
		ID:        uuid.Nil,
		Title:     "Meditations",
		Author:    "Marcus Aurelius",
		Summary:   "Lorem Ipsum",
		Genre:     "Self Help",
		Year:      2000,
		Publisher: "Random",
		ImageURI:  "image.com/woo-hoo",
	}
}

func TestStore_Delete(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	id1 := uuid.New()
	tests := []struct {
		desc string
		id   uuid.UUID
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			id1,
			nil,
			mock.ExpectExec(`DELETE FROM book WHERE id=$1`).WithArgs(id1).WillReturnResult(sqlmock.NewResult(0, 1)),
		},
		{
			"DB error",
			id1,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(`DELETE FROM book WHERE id=$1`).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		err := s.Delete(ctx, tc.id)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestStore_Get(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	row1 := sqlmock.NewRows([]string{"id", "title", "author", "genre", "image_uri"}).
		AddRow(uuid.Nil, "Meditations", "Marcus Aurelius", "Self Help", "image.com/woo-hoo")

	row2 := sqlmock.NewRows([]string{"id", "title", "author", "genre", "image_uri"}).
		AddRow(uuid.Nil, "Meditations", "Marcus Aurelius", "Self Help", "image.com/woo-hoo")

	query1 := `select id, title, author, genre, image_uri from book offset $1 limit $2;`
	query2 := `select id, title, author, genre, image_uri from book where genre = 'Self Help' AND year = 2000 offset $1 limit $2;`

	books := bookRes()

	page := &model.Page{
		Offset: 10,
		Size:   20,
	}

	tests := []struct {
		desc    string
		filters *model.Filters
		res     []model.BookRes
		err     error
		query   *sqlmock.ExpectedQuery
	}{
		{
			"Without filters",
			&model.Filters{},
			books,
			nil,
			mock.ExpectQuery(query1).WithArgs(page.Offset, page.Size).WillReturnRows(row1),
		},
		{
			"With filters",
			filter(),
			books,
			nil,
			mock.ExpectQuery(query2).WithArgs(page.Offset, page.Size).WillReturnRows(row2),
		},
		{
			"DB Error",
			&model.Filters{},
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectQuery(query1).WithArgs(page.Offset, page.Size).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		books, err := s.Get(ctx, page, tc.filters)
		assert.Equal(t, tc.res, books, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestStore_Update(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := uuid.New()
	book1 := getNewBook(id)
	bookRes1 := getNewBookRes(id)

	row := sqlmock.NewRows([]string{"title", "author", "summary", "genre", "year", "publisher", "image_uri"}).
		AddRow(book1.Title, book1.Author, book1.Summary, book1.Genre, book1.Year, book1.Publisher, book1.ImageURI)
	tests := []struct {
		desc string
		book *model.Book
		resp *model.BookRes
		err  error
		exec *sqlmock.ExpectedExec
		quer *sqlmock.ExpectedQuery
	}{
		{
			"Success",
			book1,
			bookRes1,
			nil,
			mock.ExpectExec(updateBook).
				WithArgs(book1.Title, book1.Author, book1.Summary,
					book1.Genre, book1.Year, book1.RegNum, book1.Publisher,
					book1.Language, book1.ImageURI, book1.ID).
				WillReturnResult(sqlmock.NewResult(0, 1)),
			mock.ExpectQuery(getByID).WithArgs(book1.ID).WillReturnRows(row),
		},
		{
			"DB error",
			book1,
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(updateBook).WithArgs(
				book1.Title, book1.Author, book1.Summary,
				book1.Genre, book1.Year, book1.RegNum, book1.Publisher,
				book1.Language, book1.ImageURI, book1.ID).WillReturnError(
				errors.Error("DB Error")),
			mock.ExpectQuery(getByID).WithArgs(book1.ID).WillReturnRows(row),
		},
	}

	for _, tc := range tests {
		bookRes, err := s.Update(ctx, tc.book)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, bookRes, tc.desc)
	}
}

func TestStore_GetByID(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	row := sqlmock.NewRows([]string{"title", "author", "summary", "genre", "year", "publisher", "image_uri"}).
		AddRow("Meditations", "Marcus Aurelius", "Lorem Ipsum", "Self Help", 2000, "Random", "image.com/woo-hoo")

	book := bookRes2()

	id1 := uuid.New()
	id2 := uuid.New()

	tests := []struct {
		desc  string
		id    uuid.UUID
		res   *model.BookRes
		err   error
		query *sqlmock.ExpectedQuery
	}{
		{
			"Success",
			uuid.Nil,
			book,
			nil,
			mock.ExpectQuery(getByID).WithArgs(uuid.Nil).WillReturnRows(row),
		},
		{
			"Not Exists",
			id1,
			nil,
			errors.EntityNotFound{Entity: "book", ID: id1.String()},
			mock.ExpectQuery(getByID).WithArgs(id1).WillReturnError(sql.ErrNoRows),
		},
		{
			"DB Error",
			id2,
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectQuery(getByID).WithArgs(id2).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		book, err := s.GetByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.res, book, tc.desc)
	}
}

func TestStore_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := uuid.New()
	book1 := getNewBook(id)
	bookRes1 := getNewBookRes(id)
	tests := []struct {
		desc string
		book *model.Book
		resp *model.BookRes
		err  error
		exec *sqlmock.ExpectedExec
	}{
		{
			"Success",
			book1,
			bookRes1,
			nil,
			mock.ExpectExec(createBook).WithArgs(sqlmock.AnyArg(), book1.Title,
				book1.Author, book1.Summary, book1.Genre, book1.Year, book1.RegNum,
				book1.Publisher, book1.Language, book1.ImageURI).WillReturnResult(
				sqlmock.NewResult(0, 1),
			),
		},
		{
			"DB error",
			book1,
			nil,
			errors.DB{Err: errors.DB{}},
			mock.ExpectExec(createBook).WithArgs(sqlmock.AnyArg(), book1.Title,
				book1.Author, book1.Summary, book1.Genre, book1.Year, book1.RegNum,
				book1.Publisher, book1.Language, book1.ImageURI).WillReturnError(
				errors.DB{Err: nil}),
		},
	}

	for _, tc := range tests {
		book, err := s.Create(ctx, tc.book)
		if err == nil {
			tc.resp.ID = book.ID
		}

		assert.Equal(t, tc.resp, book, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestStore_GetFilters(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	rows := sqlmock.NewRows([]string{"year"}).AddRow(2000).AddRow(2005).AddRow(2011).AddRow(2021)

	tests := []struct {
		desc   string
		filter string
		res    []string
		err    error
		query  *sqlmock.ExpectedQuery
	}{
		{
			"Success",
			"year",
			[]string{"2000", "2005", "2011", "2021"},
			nil,
			mock.ExpectQuery(`select distinct year from book;`).WillReturnRows(rows),
		},
		{
			"DB Error",
			"gibberish",
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectQuery(`select distinct gibberish from book;`).WillReturnError(errors.Error("DB Error")),
		},
	}

	for _, tc := range tests {
		filters, err := s.GetFilters(ctx, tc.filter)
		assert.Equal(t, tc.res, filters, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

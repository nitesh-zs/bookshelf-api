package book

import (
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

func getID() uuid.UUID {
	return uuid.New()
}

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

func TestStore_Delete(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	id1 := getID()
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

func TestStore_Update(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := getID()
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
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnResult(sqlmock.NewResult(0, 1)),
			mock.ExpectQuery(getByID).WithArgs(book1.ID).WillReturnRows(row),
		},
		{
			"DB error",
			book1,
			nil,
			errors.DB{Err: errors.Error("DB Error")},
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnError(errors.Error("DB Error")),
			mock.ExpectQuery(getByID).WithArgs(book1.ID).WillReturnRows(row),
		},
		{
			desc: "error",
			book: nil,
			resp: nil,
			err:  errors.Error("No object to update"),
		},
	}

	for _, tc := range tests {
		bookRes, err := s.Update(ctx, tc.book)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, bookRes, tc.desc)
	}
}

func TestStore_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := getID()
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
			mock.ExpectExec(createBook).WithArgs(book1.ID.String(), book1.Title,
				book1.Author, book1.Summary, book1.Genre, book1.Year, book1.RegNum,
				book1.Publisher, book1.Language, book1.ImageURI).WillReturnResult(
				sqlmock.NewResult(0, 1),
			),
		},
		{
			"DB error",
			book1,
			nil,
			errors.DB{Err: errors.Error("cannot create object")},
			mock.ExpectExec(getUpdateQuery(book1)).WillReturnError(errors.DB{Err: errors.Error("cannot create object")}),
		},
		{
			desc: "error",
			book: nil,
			resp: nil,
			err:  errors.Error("No object to create"),
		},
	}

	for _, tc := range tests {
		bookRes1, err := s.Create(ctx, tc.book)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, bookRes1, tc.desc)
	}
}

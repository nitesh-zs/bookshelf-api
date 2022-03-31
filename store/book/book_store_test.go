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
		assert.Equal(t, tc.res, book, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

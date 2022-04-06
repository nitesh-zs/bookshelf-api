package book

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"

	"github.com/nitesh-zs/bookshelf-api/mocks"
	"github.com/nitesh-zs/bookshelf-api/model"
)

func getUser() *model.User {
	return &model.User{
		Type: "admin",
	}
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

func initializeTest(t *testing.T) (*mocks.MockBookStore, *krogo.Context, svc) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mocks.NewMockBookStore(mockCtrl)
	k := krogo.New()
	ctx := krogo.NewContext(nil, nil, k)
	s := New(mock)

	return mock, ctx, s
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

func filter() *model.Filters {
	return &model.Filters{
		Genre: "Self Help",
		Year:  2000,
	}
}

func TestSvc_Delete(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	id1 := uuid.New()

	tests := []struct {
		desc             string
		id               uuid.UUID
		user             *model.User
		err              error
		mockStoreDelete  *gomock.Call
		mockStoreGetByID *gomock.Call
	}{
		{
			"success",
			id1,
			&model.User{Type: "admin"},
			nil,
			mock.EXPECT().Delete(ctx, id1).Return(nil),
			mock.EXPECT().GetByID(ctx, id1).Return(getNewBookRes(id1), nil),
		},
		{
			"DB Error",
			id1,
			&model.User{Type: "admin"},
			errors.DB{Err: errors.DB{}},
			mock.EXPECT().Delete(ctx, id1).Return(errors.DB{}),
			mock.EXPECT().GetByID(ctx, id1).Return(getNewBookRes(id1), nil),
		},
		{
			desc:            "DB Error",
			id:              uuid.Nil,
			user:            &model.User{Type: "admin"},
			err:             errors.InvalidParam{Param: []string{"id"}},
			mockStoreDelete: mock.EXPECT().Delete(ctx, id1).Return(errors.InvalidParam{Param: []string{"id"}}),
		},
		{
			"Entity not exist",
			id1,
			&model.User{Type: "admin"},
			errors.EntityNotFound{},
			mock.EXPECT().GetByID(ctx, id1).Return(nil, errors.EntityNotFound{}),
			mock.EXPECT().Delete(ctx, id1).Return(errors.EntityNotFound{}),
		},
	}

	for _, tc := range tests {
		err := s.Delete(ctx, tc.id, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Get(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	books := bookRes()
	page := &model.Page{}

	tests := []struct {
		desc      string
		filters   *model.Filters
		res       []model.BookRes
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Without filters",
			&model.Filters{},
			books,
			nil,
			[]*gomock.Call{
				mock.EXPECT().Get(ctx, page, &model.Filters{}).Return(books, nil),
			},
		},
		{
			"With filters",
			filter(),
			books,
			nil,
			[]*gomock.Call{
				mock.EXPECT().Get(ctx, page, filter()).Return(books, nil),
			},
		},
		{
			"DB Error",
			&model.Filters{},
			nil,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().Get(ctx, page, &model.Filters{}).Return(nil, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		books, err := s.Get(ctx, page, tc.filters)
		assert.Equal(t, tc.res, books, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := uuid.New()
	tests := []struct {
		desc      string
		book      *model.Book
		user      *model.User
		resp      *model.BookRes
		err       error
		mockStore *gomock.Call
	}{
		{
			"success",
			getNewBook(id),
			getUser(),
			getNewBookRes(id),
			nil,
			mock.EXPECT().Create(ctx, getNewBook(id)).Return(getNewBookRes(id), nil),
		},
		{
			desc:      "DB Error",
			book:      getNewBook(id),
			user:      getUser(),
			resp:      nil,
			err:       errors.DB{},
			mockStore: mock.EXPECT().Create(ctx, getNewBook(id)).Return(nil, errors.DB{}),
		},
		{
			desc: "no object",
			book: nil,
			user: getUser(),
			resp: nil,
			err:  errors.InvalidParam{Param: []string{"invalid body request"}},
		},
	}

	for _, tc := range tests {
		resp, err := s.Create(ctx, tc.book, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, resp, tc.desc)
	}
}

func TestSvc_GetByID(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	book := bookRes2()
	id1 := uuid.New()
	id2 := uuid.New()

	tests := []struct {
		desc      string
		id        uuid.UUID
		res       *model.BookRes
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Success",
			uuid.Nil,
			book,
			nil,
			[]*gomock.Call{
				mock.EXPECT().GetByID(ctx, uuid.Nil).Return(book, nil),
			},
		},
		{
			"Not Exists",
			id1,
			nil,
			errors.EntityNotFound{Entity: "book", ID: id1.String()},
			[]*gomock.Call{
				mock.EXPECT().GetByID(ctx, id1).Return(nil, errors.EntityNotFound{Entity: "book", ID: id1.String()}),
			},
		},
		{
			"DB Error",
			id2,
			nil,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().GetByID(ctx, id2).Return(nil, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		book, err := s.GetByID(ctx, tc.id)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.res, book, tc.desc)
	}
}

func TestSvc_Update(t *testing.T) {
	id := uuid.New()
	mock, ctx, s := initializeTest(t)
	tests := []struct {
		desc      string
		book      *model.Book
		user      *model.User
		resp      *model.BookRes
		err       error
		mockStore *gomock.Call
	}{
		{
			"error",
			getNewBook(id),
			getUser(),
			nil,
			errors.DB{Err: errors.DB{}},
			mock.EXPECT().Update(ctx, getNewBook(id)).Return(nil, errors.DB{}),
		},
		{
			"success",
			getNewBook(id),
			getUser(),
			getNewBookRes(id),
			nil,
			mock.EXPECT().Update(ctx, getNewBook(id)).Return(getNewBookRes(id), nil),
		},
		{
			desc: "no object",
			book: nil,
			user: getUser(),
			resp: nil,
			err:  errors.InvalidParam{Param: []string{"invalid body request"}},
		},
	}

	for _, tc := range tests {
		resp, err := s.Update(ctx, tc.book, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, resp, tc.desc)
	}
}

func TestSvc_GetFilters(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	filters := []string{"2000", "2005", "2011", "2021"}

	tests := []struct {
		desc      string
		filter    string
		res       []string
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Success",
			"year",
			filters,
			nil,
			[]*gomock.Call{
				mock.EXPECT().GetFilters(ctx, "year").Return(filters, nil),
			},
		},
		{
			"DB Error",
			"random",
			nil,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().GetFilters(ctx, "random").Return(nil, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		f, err := s.GetFilters(ctx, tc.filter)
		assert.Equal(t, tc.res, f, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

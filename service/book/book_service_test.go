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
		assert.Equal(t, tc.res, book, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
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

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

func initializeTest(t *testing.T) (*mocks.MockBookStore, *krogo.Context, svc) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mocks.NewMockBookStore(mockCtrl)
	k := krogo.New()
	ctx := krogo.NewContext(nil, nil, k)
	s := New(mock)

	return mock, ctx, s
}

func TestSvc_Delete(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	id1 := getID()
	tests := []struct {
		desc      string
		id        uuid.UUID
		user      *model.User
		err       error
		mockStore *gomock.Call
	}{
		{
			"success",
			id1,
			&model.User{Type: "admin"},
			nil,
			mock.EXPECT().Delete(ctx, id1).Return(nil),
		},
		{
			"DB Error",
			id1,
			&model.User{Type: "admin"},
			errors.DB{},
			mock.EXPECT().Delete(ctx, id1).Return(errors.DB{}),
		},
		{
			"DB Error",
			uuid.Nil,
			&model.User{Type: "admin"},
			errors.Error("invalid uuid"),
			mock.EXPECT().Delete(ctx, id1).Return(errors.Error("invalid uuid")),
		},
	}

	for _, tc := range tests {
		err := s.Delete(ctx, tc.id, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	id := getID()
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
			"DB Error",
			getNewBook(id),
			getUser(),
			nil,
			errors.DB{},
			mock.EXPECT().Create(ctx, getNewBook(id)).Return(nil, errors.DB{}),
		},
		{
			desc: "no object",
			book: nil,
			user: getUser(),
			resp: nil,
			err:  errors.Error("No object to create"),
		},
	}

	for _, tc := range tests {
		resp, err := s.Create(ctx, tc.book, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, resp, tc.desc)
	}
}

// nolint
func TestSvc_Update(t *testing.T) {
	id := getID()
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
			errors.DB{},
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
			err:  errors.Error("No object to update"),
		},
	}

	for _, tc := range tests {
		resp, err := s.Update(ctx, tc.book, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
		assert.Equal(t, tc.resp, resp, tc.desc)
	}
}

package book

import (
	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/krogertechnology/krogo/pkg/errors"
	"github.com/krogertechnology/krogo/pkg/krogo"
	"github.com/nitesh-zs/bookshelf-api/mocks"
	"github.com/nitesh-zs/bookshelf-api/model"
	"testing"
)

var user1 = &model.User{
	Type: "admin",
}
var user2 = &model.User{
	Type: "non admin",
}
var id = uuid.New()
var book1 = &model.Book{
	ID:        id,
	Title:     "Abc",
	Author:    "X",
	Summary:   "Lorem Ipsum",
	Genre:     "Action",
	Year:      2019,
	RegNum:    "ISB8726W821",
	Publisher: "saiudhiau",
	Language:  "Hebrew",
	ImageURI:  "jncj.ajcbiauadnc.com",
}

var bookRes1 = &model.BookRes{
	ID:        id,
	Title:     "Abc",
	Author:    "X",
	Summary:   "Lorem Ipsum",
	Genre:     "Action",
	Year:      2019,
	Publisher: "saiudhiau",
	ImageURI:  "jncj.ajcbiauadnc.com",
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

	id1 := uuid.New()
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
		//{
		//	"Auth Error",
		//	id1,
		//	&model.User{Type: "non-admin"},
		//	errors.Unauthorized{},
		//	mock.EXPECT().Delete(ctx, id1).Return(errors.Unauthorized{}),
		//},
	}
	for _, tc := range tests {
		err := s.Delete(ctx, tc.id, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Create(t *testing.T) {
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
			"success",
			book1,
			user1,
			bookRes1,
			nil,
			mock.EXPECT().Create(ctx, book1).Return(bookRes1, nil),
		},
		//{
		//	desc: "unauthorised",
		//	book: book1,
		//	user: user2,
		//	resp: nil,
		//	err:  errors.Unauthorized{},
		//},
		{
			"DB Error",
			book1,
			user1,
			nil,
			errors.DB{},
			mock.EXPECT().Create(ctx, book1).Return(nil, errors.DB{}),
		},
		{
			desc: "no object",
			book: nil,
			user: user1,
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

func TestSvc_Update(t *testing.T) {
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
			"success",
			book1,
			user1,
			bookRes1,
			nil,
			mock.EXPECT().Update(ctx, book1).Return(bookRes1, nil),
		},
		//{
		//	desc: "unauthorised",
		//	book: book1,
		//	user2,
		//	nil,
		//	errors.Unauthorized{},
		//},
		{
			"error",
			book1,
			user1,
			nil,
			errors.DB{},
			mock.EXPECT().Update(ctx, book1).Return(nil, errors.DB{}),
		},
		{
			desc: "no object",
			book: nil,
			user: user1,
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

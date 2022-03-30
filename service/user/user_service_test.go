package user

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

func user1() *model.User {
	return &model.User{
		ID:    uuid.New(),
		Email: "nitesh.saxena@zopsmart.com",
		Name:  "Nitesh",
		Type:  "admin",
	}
}

func user2() *model.User {
	return &model.User{
		ID:    uuid.New(),
		Email: "abc@abc.com",
		Name:  "John",
		Type:  "general",
	}
}

func initializeTest(t *testing.T) (*mocks.MockUserStore, *krogo.Context, svc) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := mocks.NewMockUserStore(mockCtrl)

	k := krogo.New()
	ctx := krogo.NewContext(nil, nil, k)
	s := New(mock)

	return mock, ctx, s
}

func TestSvc_Exists(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	tests := []struct {
		desc      string
		email     string
		res       bool
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Exists",
			"nitesh.saxena@zopsmart.com",
			true,
			nil,
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "nitesh.saxena@zopsmart.com").Return(user1(), nil),
			},
		},
		{
			"Not Exists",
			"abc@abc.com",
			false,
			errors.EntityNotFound{Entity: "user", ID: "abc@abc.com"},
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "abc@abc.com").Return(nil, errors.EntityNotFound{Entity: "user", ID: "abc@abc.com"}),
			},
		},
		{
			"Server error",
			"xyz@xyz.com",
			false,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "xyz@xyz.com").Return(nil, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		exists, err := s.Exists(ctx, tc.email)

		assert.Equal(t, tc.res, exists, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_Create(t *testing.T) {
	mock, ctx, s := initializeTest(t)
	user1 := user1()
	user2 := &model.User{}

	tests := []struct {
		desc      string
		user      *model.User
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Success",
			user1,
			nil,
			[]*gomock.Call{
				mock.EXPECT().Create(ctx, user1).Return(nil),
			},
		},
		{
			"Server Error",
			user2,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().Create(ctx, user2).Return(errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		err := s.Create(ctx, tc.user)
		assert.Equal(t, tc.err, err, tc.desc)
	}
}

func TestSvc_IsAdmin(t *testing.T) {
	mock, ctx, s := initializeTest(t)

	tests := []struct {
		desc      string
		email     string
		res       bool
		err       error
		mockStore []*gomock.Call
	}{
		{
			"Admin",
			"nitesh.saxena@zopsmart.com",
			true,
			nil,
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "nitesh.saxena@zopsmart.com").Return(user1(), nil),
			},
		},
		{
			"Non Admin",
			"abc@abc.com",
			false,
			nil,
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "abc@abc.com").Return(user2(), nil),
			},
		},
		{
			"Error",
			"xyz@xyz.com",
			false,
			errors.DB{},
			[]*gomock.Call{
				mock.EXPECT().GetByEmail(ctx, "xyz@xyz.com").Return(nil, errors.DB{}),
			},
		},
	}

	for _, tc := range tests {
		res, err := s.IsAdmin(ctx, tc.email)

		assert.Equal(t, tc.res, res, tc.desc)
		assert.Equal(t, tc.err, err, tc.desc)
	}

}

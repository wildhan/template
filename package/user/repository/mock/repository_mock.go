package mock

import (
	"template/package/user/model"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (_m *UserRepoMock) GetUsers() ([]model.User, error) {
	args := _m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (_m *UserRepoMock) AddUser(user model.User) error {
	args := _m.Called(user)
	return args.Error(0)
}

func (_m *UserRepoMock) EditUser(user model.User) error {
	args := _m.Called(user)
	return args.Error(0)
}

package usecase_test

import (
	"errors"
	"template/package/user/model"
	"template/package/user/repository/mock"
	"template/package/user/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mockUser model.User = model.User{
		Id:        "938866d2-b095-4671-9502-cfd21405a57a",
		Username:  "JohnDoe",
		FirstName: "John",
		LastName:  "Doe",
	}
)

func TestGetUsers(t *testing.T) {
	repo := &mock.UserRepoMock{}
	usecase := usecase.NewUserUsecase(repo)

	users := []model.User{mockUser}
	repo.On("GetUsers").Return(users, nil)

	result, err := usecase.GetUsers()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, users, result)
	repo.AssertExpectations(t)
}

func TestAddUser(t *testing.T) {
	repo := &mock.UserRepoMock{}
	usecase := usecase.NewUserUsecase(repo)
	testCases := []struct {
		name    string
		user    model.User
		mockErr error
		wantErr bool
	}{
		{
			name:    "ValidUser",
			user:    mockUser,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "InvalidUser",
			user:    mockUser,
			mockErr: errors.New("invalid user"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.On("AddUser", tc.user).Return(tc.mockErr).Once()
			err := usecase.AddUser(mockUser)

			if !tc.wantErr {
				assert.NoError(t, err, err)
			} else {
				assert.Error(t, err, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestEditUser(t *testing.T) {
	repo := &mock.UserRepoMock{}
	usecase := usecase.NewUserUsecase(repo)
	testCases := []struct {
		name    string
		user    model.User
		mockErr error
		wantErr bool
	}{
		{
			name:    "ValidUser",
			user:    mockUser,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "InvalidUser",
			user:    mockUser,
			mockErr: errors.New("invalid user"),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo.On("EditUser", tc.user).Return(tc.mockErr).Once()
			err := usecase.EditUser(mockUser)

			if !tc.wantErr {
				assert.NoError(t, err, err)
			} else {
				assert.Error(t, err, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

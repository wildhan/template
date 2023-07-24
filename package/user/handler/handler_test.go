package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"template/package/user/handler"
	"template/package/user/model"
	"template/package/user/usecase/mock"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	mock := &mock.UserUsecaseMock{}
	users := []model.User{mockUser}

	mock.On("GetUsers").Return(users, nil)
	handler := handler.NewUserHandler(mock)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.GET, "/", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.GetUsers(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusOK, rec.Code)
	mock.AssertExpectations(t)
}

func TestAddUser(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("AddUser", mockUser).Return(nil)

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/add", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.AddUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusOK, rec.Code)
	mock.AssertExpectations(t)
}

func TestAddUserDuplicate(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("AddUser", mockUser).Return(errors.New("ERROR: duplicate key value violates unique constraint \"username\" (SQLSTATE 23505)"))

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/add", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.AddUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mock.AssertExpectations(t)
}

func TestAddUserCantNull(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("AddUser", mockUser).Return(errors.New("ERROR: null value in column \"username\" of relation \"user_profile\" violates not-null constraint (SQLSTATE 23502)"))

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/add", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.AddUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mock.AssertExpectations(t)
}

func TestEditUser(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("EditUser", mockUser).Return(nil)

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/edit", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.EditUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusOK, rec.Code)
	mock.AssertExpectations(t)
}

func TestEditUserDuplicate(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("EditUser", mockUser).Return(errors.New("ERROR: duplicate key value violates unique constraint \"username\" (SQLSTATE 23505)"))

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/edit", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.EditUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mock.AssertExpectations(t)
}

func TestEditUserCantNull(t *testing.T) {
	mock := &mock.UserUsecaseMock{}
	mock.On("EditUser", mockUser).Return(errors.New("ERROR: null value in column \"username\" of relation \"user_profile\" violates not-null constraint (SQLSTATE 23502)"))

	handler := handler.NewUserHandler(mock)

	j, err := json.Marshal(mockUser)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequestWithContext(context.TODO(), echo.POST, "/edit", strings.NewReader(string(j)))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	errReq := handler.EditUser(c)
	require.NoError(t, errReq)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mock.AssertExpectations(t)
}

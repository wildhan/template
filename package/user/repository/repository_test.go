package repository

import (
	"template/config/database"
	"template/package/user/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
)

var (
	expectedUser model.User = model.User{
		Id:        "938866d2-b095-4671-9502-cfd21405a57a",
		Username:  "JohnDoe",
		FirstName: "John",
		LastName:  "Doe",
	}
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("SQLMock Connection Failed")
	}
	defer db.Close()

	dbConn, err := database.CreateMockConnection(postgres.Config{Conn: db})
	if err != nil {
		t.Fatalf("DB Connection Failed")
	}

	expectedResult := []model.User{expectedUser}

	mock.ExpectQuery("SELECT id, username, first_name, last_name FROM public.user_profile").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "first_name", "last_name"}).
			AddRow(expectedUser.Id, expectedUser.Username, expectedUser.FirstName, expectedUser.LastName))

	UserRepo := NewUserRepo(dbConn)
	result, err := UserRepo.GetUsers()

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, expectedResult, result)
}

func TestAddUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("SQLMock Connection Failed")
	}
	defer db.Close()

	dbConn, err := database.CreateMockConnection(postgres.Config{Conn: db})
	if err != nil {
		t.Fatalf("DB Connection Failed")
	}

	query := `INSERT INTO public.user_profile`

	mock.ExpectExec(query).
		WithArgs(expectedUser.Username, expectedUser.FirstName, expectedUser.LastName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	UserRepo := NewUserRepo(dbConn)
	err = UserRepo.AddUser(expectedUser)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEditUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("SQLMock Connection Failed")
	}
	defer db.Close()

	dbConn, err := database.CreateMockConnection(postgres.Config{Conn: db})
	if err != nil {
		t.Fatalf("DB Connection Failed")
	}

	query := `UPDATE public.user_profile`

	mock.ExpectExec(query).
		WithArgs(expectedUser.Username, expectedUser.FirstName, expectedUser.LastName, expectedUser.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	UserRepo := NewUserRepo(dbConn)
	err = UserRepo.EditUser(expectedUser)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

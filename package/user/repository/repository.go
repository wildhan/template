package repository

import (
	"template/config/database"
	"template/lib/helper"
	"template/package/user/model"
)

type userRepo struct {
	dbConn *database.DbConnection
}

func NewUserRepo(dbConn *database.DbConnection) UserRepo {
	return &userRepo{dbConn}
}

type UserRepo interface {
	GetUsers() ([]model.User, error)
	AddUser(model.User) error
	EditUser(model.User) error
}

func (r *userRepo) GetUsers() ([]model.User, error) {
	users := make([]model.User, 0)
	db := r.dbConn.DB

	query := `SELECT * FROM public.user_profile;`
	err := db.Raw(query).Scan(&users).Error

	return users, err
}

func (r *userRepo) AddUser(user model.User) error {
	db := r.dbConn.DB
	params := make([]interface{}, 0)

	params = append(params, helper.EmptyStringToNull(user.Username))
	params = append(params, helper.EmptyStringToNull(user.FirstName))
	params = append(params, helper.EmptyStringToNull(user.LastName))

	query := `INSERT INTO public.user_profile
	(username, first_name, last_name)
	VALUES(?, ?, ?);`
	err := db.Exec(query, params...).Error

	return err
}

func (r *userRepo) EditUser(user model.User) error {
	db := r.dbConn.DB
	params := make([]interface{}, 0)

	params = append(params, helper.EmptyStringToNull(user.Username))
	params = append(params, helper.EmptyStringToNull(user.FirstName))
	params = append(params, helper.EmptyStringToNull(user.LastName))

	params = append(params, user.Id)

	query := `UPDATE public.user_profile
	SET username=?, first_name=?, last_name=?
	WHERE id=?::uuid;`
	err := db.Exec(query, params...).Error

	return err
}

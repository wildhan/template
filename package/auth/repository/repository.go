package repository

import (
	"template/config/database"
	"template/lib/helper"
	"template/package/auth/model"
)

type authRepo struct {
	dbConn *database.DbConnection
}

func NewAuthRepo(dbConn *database.DbConnection) AuthRepo {
	return &authRepo{dbConn}
}

type AuthRepo interface {
	InsertUserAuth(model.UserAuth) error
}

func (r *authRepo) InsertUserAuth(user model.UserAuth) error {
	db := r.dbConn.DB
	params := make([]interface{}, 0)

	params = append(params, helper.EmptyStringToNull(user.Username))
	params = append(params, helper.EmptyStringToNull(user.Email))
	params = append(params, helper.EmptyStringToNull(user.HashPassword))

	query := `INSERT INTO public.users_auth
						(username, email, hash_password)
						VALUES(?, ?, ?);`
	return db.Exec(query, params...).Error
}

package repository

import (
	"template/lib/helper"
	"template/package/auth/model"

	"gorm.io/gorm"
)

type authRepo struct{}

func NewAuthRepo() AuthRepo {
	return &authRepo{}
}

type AuthRepo interface {
	InsertUserAuth(db gorm.DB, user model.UserAuth) (model.UserAuth, error)
	InsertUserProfile(db gorm.DB, user model.UserAuth) (model.UserAuth, error)
	ReadUserAuth(db gorm.DB, username string) (*model.UserAuth, error)
}

func (r *authRepo) InsertUserAuth(db gorm.DB, user model.UserAuth) (model.UserAuth, error) {
	var id *string
	params := make([]interface{}, 0)

	params = append(params, helper.EmptyStringToNull(user.Username))
	params = append(params, helper.EmptyStringToNull(user.Email))
	params = append(params, helper.EmptyStringToNull(user.HashPassword))

	query := `INSERT INTO public.users_auth
						(username, email, hash_password)
						VALUES(?, ?, ?)
						RETURNING id;`
	if err := db.Raw(query, params...).Scan(&id).Error; err != nil {
		return user, err
	}

	user.Id = *id

	return user, nil
}

func (r *authRepo) InsertUserProfile(db gorm.DB, user model.UserAuth) (model.UserAuth, error) {
	query := `INSERT INTO public.users_profile
						(id)
						VALUES(?);`
	return user, db.Exec(query, user.Id).Error
}

func (r *authRepo) ReadUserAuth(db gorm.DB, username string) (userAuth *model.UserAuth, err error) {
	query := `SELECT id, username, email, hash_password
						FROM public.users_auth
						WHERE email=?;`
	err = db.Raw(query, username).Scan(&userAuth).Error
	if err != nil {
		return nil, err
	}
	return userAuth, nil
}

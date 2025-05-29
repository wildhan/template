package usecase

import (
	"template/config/database"
	"template/lib/helper"
	"template/package/auth/model"
	"template/package/auth/repository"
)

type authUsecase struct {
	dbConn *database.DbConnection
	repo   repository.AuthRepo
}

func NewAuthUsecase(dbConn *database.DbConnection, repo repository.AuthRepo) AuthUsecase {
	return &authUsecase{
		dbConn: dbConn,
		repo:   repo,
	}
}

type AuthUsecase interface {
	RegistrationUser(user model.UserAuth) error
}

func (uc *authUsecase) RegistrationUser(user model.UserAuth) error {
	tx := uc.dbConn.DB.Begin()

	defer tx.Rollback()

	if hash, err := helper.HashPassword(user.Password); err != nil {
		return err
	} else {
		user.HashPassword = hash
	}

	user, err := uc.repo.InsertUserAuth(*tx, user)
	if err != nil {
		return err
	}

	if _, err := uc.repo.InsertUserProfile(*tx, user); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

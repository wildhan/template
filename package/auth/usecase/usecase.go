package usecase

import (
	"errors"
	"template/config/authorization"
	"template/config/database"
	"template/lib/helper"
	"template/lib/response"
	"template/package/auth/model"
	"template/package/auth/repository"
)

type authUsecase struct {
	dbConn    *database.DbConnection
	authToken authorization.AuthToken
	repo      repository.AuthRepo
}

func NewAuthUsecase(
	dbConn *database.DbConnection,
	authToken authorization.AuthToken,
	repo repository.AuthRepo,
) AuthUsecase {
	return &authUsecase{
		dbConn:    dbConn,
		authToken: authToken,
		repo:      repo,
	}
}

type AuthUsecase interface {
	RegistrationUser(user model.UserAuth) error
	LoginUser(loginParams model.LoginParameter) (*model.LoginResponse, error)
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

	userProfile := model.UserProfile{
		Id:            user.Id,
		CreatedByUser: &user.Id,
	}

	if _, err := uc.repo.InsertUserProfile(*tx, userProfile); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (uc *authUsecase) LoginUser(loginParams model.LoginParameter) (*model.LoginResponse, error) {
	tx := uc.dbConn.DB.Begin()

	defer tx.Rollback()

	user, err := uc.repo.ReadUserAuth(*tx, loginParams.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New(response.ERROR_AUTH_USER_NOT_FOUND)
	}

	isPassMatch := helper.CheckPasswordHash(loginParams.Password, user.HashPassword)
	if !isPassMatch {
		return nil, errors.New(response.ERROR_AUTH_PASS_NOT_MATCH)
	}

	token, err := uc.authToken.GenerateToken(authorization.TokenContainer{
		Sender: user.Id,
	})
	if err != nil {
		return nil, err
	}

	resp := model.LoginResponse{
		Username: user.Username,
		Token:    token,
	}

	return &resp, nil
}

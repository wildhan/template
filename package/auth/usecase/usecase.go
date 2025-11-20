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
	RefreshToken(refreshToken string) (*model.LoginResponse, error)
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

	tc := authorization.TokenContainer{
		Sender: user.Id,
	}

	token, err := uc.authToken.GenerateToken(tc)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.authToken.GenerateRefreshToken(authorization.RefreshTokenContainer{
		TC:    tc,
		Email: loginParams.Email,
	})
	if err != nil {
		return nil, err
	}

	resp := model.LoginResponse{
		Username:     user.Username,
		Token:        token,
		RefreshToken: refreshToken,
	}

	return &resp, nil
}

func (uc *authUsecase) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	tx := uc.dbConn.DB

	rtc, err := uc.authToken.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	user, err := uc.repo.ReadUserAuth(*tx, rtc.Email)
	if err != nil {
		return nil, err
	}

	newToken, err := uc.authToken.GenerateToken(rtc.TC)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := uc.authToken.GenerateRefreshToken(*rtc)
	if err != nil {
		return nil, err
	}
	resp := model.LoginResponse{
		Username:     user.Username,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}

	return &resp, nil
}

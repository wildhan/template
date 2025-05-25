package usecase

import (
	"template/lib/helper"
	"template/package/auth/model"
	"template/package/auth/repository"
)

type authUsecase struct {
	repo repository.AuthRepo
}

func NewAuthUsecase(repo repository.AuthRepo) AuthUsecase {
	return &authUsecase{repo}
}

type AuthUsecase interface {
	RegistrationUser(model.UserAuth) error
}

func (uc *authUsecase) RegistrationUser(user model.UserAuth) error {
	if hash, err := helper.HashPassword(user.Password); err != nil {
		return err
	} else {
		user.HashPassword = hash
	}

	return uc.repo.InsertUserAuth(user)
}

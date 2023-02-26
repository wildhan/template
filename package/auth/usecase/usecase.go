package usecase

import "template/package/auth/repository"

type authUsecase struct {
	repo repository.AuthRepo
}

func NewAuthUsecase(repo repository.AuthRepo) AuthUsecase {
	return &authUsecase{repo}
}

type AuthUsecase interface {
}

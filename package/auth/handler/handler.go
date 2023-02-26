package handler

import (
	"template/package/auth/usecase"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	uc usecase.AuthUsecase
}

func NewAuthHandler(uc usecase.AuthUsecase) *authHandler {
	return &authHandler{uc}
}

func (h *authHandler) Mount(g *echo.Group) {

}

package handler

import (
	"encoding/json"
	"io"
	"strings"
	"template/lib/log"
	"template/lib/response"
	"template/lib/validator"
	"template/package/auth/model"
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
	g.POST("/", h.Registration)
	g.GET("/", h.Login)
	g.GET("/refresh-token", h.RefreshToken)
}

func (h *authHandler) Registration(e echo.Context) error {
	log.Info("Registration ...")

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		return response.ToJson(e).BadRequest("Failed get body")
	}

	user := model.UserAuth{}
	if err = json.Unmarshal(body, &user); err != nil {
		return response.ToJson(e).BadRequest("Failed unmarshal")
	}

	if fieldNotFailed := validator.JsonValidator(user, validator.CUSTOM_VALIDATION_PASSWORD); fieldNotFailed != nil {
		return response.ToJson(e).UnprocessableEntity("Failed Field Validation", fieldNotFailed)
	}

	if err := h.uc.RegistrationUser(user); err != nil {
		return response.GetResponseErorr(e, err)
	}

	return response.ToJson(e).OK(nil, "Registration Success")
}

func (h *authHandler) Login(e echo.Context) error {
	log.Info("Login ...")

	body, err := io.ReadAll(e.Request().Body)
	if err != nil {
		return response.ToJson(e).BadRequest("Failed get body")
	}

	loginParams := model.LoginParameter{}
	if err = json.Unmarshal(body, &loginParams); err != nil {
		return response.ToJson(e).BadRequest("Failed unmarshal")
	}

	if fieldNotFailed := validator.JsonValidator(loginParams); fieldNotFailed != nil {
		return response.ToJson(e).UnprocessableEntity("Failed Field Validation", fieldNotFailed)
	}

	resp, err := h.uc.LoginUser(loginParams)
	if err != nil {
		return response.GetResponseErorr(e, err)
	}

	return response.ToJson(e).OK(resp, "Login Success")
}

func (h *authHandler) RefreshToken(e echo.Context) error {
	log.Info("Refresh Token ...")

	authHeader := e.Request().Header.Get("Authorization")
	if authHeader == "" {
		return response.ToJson(e).BadRequest("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return response.ToJson(e).BadRequest("invalid token format")
	}

	resp, err := h.uc.RefreshToken(parts[1])
	if err != nil {
		return response.ToJson(e).NotAcceptable("refresh token not acceptable: " + err.Error())
	}
	return response.ToJson(e).OK(resp, "Refresh Token Success")
}

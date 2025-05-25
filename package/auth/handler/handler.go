package handler

import (
	"encoding/json"
	"io"
	"template/lib/database"
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
		isDbErr, err := database.GetErrorDatabase(e, err)
		if isDbErr {
			return err
		}
		return response.ToJson(e).InternalServerError(err.Error())
	}

	return response.ToJson(e).OK(nil, "Registration Success")
}

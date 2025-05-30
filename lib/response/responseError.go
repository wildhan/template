package response

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	ERROR_DB_NOT_NULL         = "23502"
	ERROR_DB_DUPLICATE        = "23505"
	ERROR_AUTH_USER_NOT_FOUND = "User Not Found"
	ERROR_AUTH_PASS_NOT_MATCH = "Password Not Match"
)

var errorMessages = map[string]func(e echo.Context, err error) error{
	ERROR_DB_NOT_NULL:         respErrDBNotNull,
	ERROR_DB_DUPLICATE:        respErrDBDuplicate,
	ERROR_AUTH_USER_NOT_FOUND: respErrUserNotFound,
	ERROR_AUTH_PASS_NOT_MATCH: respErrPassNotMatch,
}

func GetResponseErorr(e echo.Context, err error) error {
	errorMsg := err.Error()

	// Check error from DB
	if strings.Contains(errorMsg, "SQLSTATE") {
		indStr := strings.Index(errorMsg, "SQLSTATE ")

		codeStart := indStr + len("SQLSTATE ")
		if len(errorMsg) >= codeStart+5 {
			errorMsg = errorMsg[codeStart : codeStart+5]
		}
	}

	respErr, exist := errorMessages[errorMsg]
	if !exist {
		return ToJson(e).InternalServerError(errorMsg)
	}
	return respErr(e, err)
}

func respErrDBNotNull(e echo.Context, err error) error {
	column := strings.Split(err.Error(), "\"")[1]
	return ToJson(e).BadRequest(fmt.Sprintf(`Can't Null value on column '%v'`, column))
}

func respErrDBDuplicate(e echo.Context, err error) error {
	column := strings.Split(err.Error(), "\"")[1]
	return ToJson(e).Conflict(fmt.Sprintf(`Duplicate value on column '%v'`, column), column)
}

func respErrUserNotFound(e echo.Context, err error) error {
	return ToJson(e).NotAcceptable(`User not found`)
}

func respErrPassNotMatch(e echo.Context, err error) error {
	return ToJson(e).NotAcceptable(`Password not match`)
}

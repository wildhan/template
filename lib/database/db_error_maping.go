package database

import (
	"fmt"
	"strings"
	"template/lib/response"

	"github.com/labstack/echo/v4"
)

func GetErrorDatabase(e echo.Context, err error) (bool, error) {

	switch {
	case strings.Contains(err.Error(), "(SQLSTATE 23502)"):
		column := strings.Split(err.Error(), "\"")[1]
		return true, response.ToJson(e).BadRequest(fmt.Sprintf(`Can't Null value on column '%v'`, column))
	case strings.Contains(err.Error(), "(SQLSTATE 23505)"):
		column := strings.Split(err.Error(), "\"")[1]
		return true, response.ToJson(e).Conflict(fmt.Sprintf(`Duplicate value on column '%v'`, column), column)
	default:
		return false, err
	}
}

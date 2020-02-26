package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Logoff() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusUnauthorized, "Bye!")
	}
}

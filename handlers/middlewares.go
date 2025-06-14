package handlers

import (
	"backend-template/models/user"

	"github.com/labstack/echo/v4"
)

const (
	TokenKeyName = "Robert-Connect-Token"
)

type Header struct {
	TokenID string `header:"Robert-Connect-Token"`
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		header := c.Request().Header.Get(TokenKeyName)

		if len(header) != 36 {
			return next(c)
		}

		Token, err := user.GetUserToken(header)
		if err == nil {
			c.Set("userToken", Token)
		}

		return next(c)
	}
}

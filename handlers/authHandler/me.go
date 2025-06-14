package authHandler

import (
	"backend-template/models/user"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type UpdateUserRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
}

type MeError struct {
	Message string `json:"message" example:"Incorrect token"`
}

// @Summary Get user details
// @Description Get details of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Robert-Connect-Token header string true "Session token"
// @Success 200 {object} UserResponse
// @Failure 401 {object} MeError
// @Router /auth/me [get]
func me(c echo.Context) error {

	var token user.UserToken
	if t := c.Get("userToken"); t != nil {
		token = t.(user.UserToken)
	} else {
		return c.JSON(http.StatusUnauthorized, MeError{Message: "Incorrect token"})
	}

	u, err := user.GetUserById(token.User.ID)
	if err == pgx.ErrNoRows {
		return c.JSON(http.StatusUnauthorized, MeError{Message: "Incorrect token"})
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, u.ToSelfWebDetail())
}

// @Summary Update user details
// @Description Update details of the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Robert-Connect-Token header string true "Session token"
// @Param user body UpdateUserRequest true "User update data"
// @Success 200 {object} UserResponse
// @Failure 400 {object} MeError
// @Failure 401 {object} MeError
// @Router /auth/me [post]
func updateMe(c echo.Context) error {
	var token user.UserToken
	if t := c.Get("userToken"); t != nil {
		token = t.(user.UserToken)
	} else {
		return c.JSON(http.StatusUnauthorized, MeError{Message: "Incorrect token"})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, MeError{Message: "Invalid request"})
	}

	updateUser := user.User{
		ID:        token.User.ID,
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	}

	if !user.UpdateUser(updateUser, false) {
		return c.JSON(http.StatusBadRequest, MeError{Message: "Failed to update user"})
	}

	u, err := user.GetUserById(token.User.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, MeError{Message: "Failed to retrieve update user"})
	}

	return c.JSON(http.StatusOK, u.ToSelfWebDetail())

}

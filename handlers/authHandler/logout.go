package authHandler

import (
	"backend-template/models/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LogoutResponse struct {
	Message string `json:"message" example:"Logged out successfully"`
}

type LogoutError401 struct {
	Message string `json:"message" example:"Invalid session"`
}

// @Summary Logout user
// @Description Logs out the user by revoking their token
// @Tags auth
// @Accept json
// @Produce json
// @Param Robert-Connect-Token header string true "Session token"
// @Success 200 {object} LogoutResponse
// @Failure 401 {object} LogoutError401
// @Router /auth/logout [post]
func logout(c echo.Context) error {

	var id string
	if t := c.Get("userToken"); t != nil {
		id = t.(user.UserToken).TokenID
		user.RevokeUserToken(id)
	}

	return c.JSON(http.StatusOK, LogoutResponse{Message: "Logged out successfully"})
}

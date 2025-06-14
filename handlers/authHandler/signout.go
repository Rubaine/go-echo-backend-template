package authHandler

import (
	"backend-template/models/user"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignoutForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}

type SignoutResponse struct {
	Message string `json:"message" example:"Signout successful"`
}

type SignoutError400 struct {
	Message string `json:"message" example:"Empty password"`
}

type SignoutError403 struct {
	Message string `json:"message" example:"Invalid password"`
}

// @Summary Sign out user
// @Description Sign out the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Param Robert-Connect-Token header string true "Session token"
// @Param signoutForm body SignoutForm true "Signout form"
// @Success 200 {object} SignoutResponse
// @Failure 400 {object} SignoutError400
// @Failure 403 {object} SignoutError403
// @Router /auth/signout [post]
func signout(c echo.Context) error {
	token, err := user.GetTokenFromRequest(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authentication required"})
	}

	var signoutForm SignoutForm
	if err := c.Bind(&signoutForm); err != nil || len(signoutForm.Password) == 0 {
		return c.JSON(http.StatusBadRequest, SignoutError400{Message: "Empty password"})
	}

	if check := user.PasswordCheck(token.User.ID, signoutForm.Password); !check {
		return c.JSON(http.StatusForbidden, SignoutError403{Message: "Invalid password"})
	}

	if err := user.DeleteAccount(token.User.ID); err != "" {
		return fmt.Errorf("impossible de supprimer le compte userID=%d ; error=%s", token.User.ID, err)
	}

	user.RevokeUserToken(token.TokenID)

	return c.JSON(http.StatusOK, SignoutResponse{Message: "Signout successful"})
}

package routes

import (
	"net/http"

	"github.com/shellhub-io/shellhub/api/pkg/gateway"
	"github.com/shellhub-io/shellhub/api/routes/handlers/converter"
	"github.com/shellhub-io/shellhub/pkg/errors"
	"github.com/shellhub-io/shellhub/pkg/models"
)

const (
	UpdateUserDataURL     = "/users/:id/data"
	UpdateUserPasswordURL = "/users/:id/password" //nolint:gosec
)

const (
	ParamUserID   = "id"
	ParamUserName = "username"
)

func (h *Handler) UpdateUserData(c gateway.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return err
	}

	// FIXME: API compatibility
	//
	// The UI uses the fields with error messages to identify if it is invalid or duplicated.
	if fields, err := h.service.UpdateDataUser(c.Ctx(), &user, c.Param(ParamUserID)); err != nil {
		e, ok := err.(errors.Error)
		if !ok {
			return err
		}

		return c.JSON(converter.FromErrServiceToHTTPStatus(e.Code), fields)
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateUserPassword(c gateway.Context) error {
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.UpdatePasswordUser(c.Ctx(), req.CurrentPassword, req.NewPassword, c.Param(ParamUserID)); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

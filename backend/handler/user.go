package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/repository"
	"github.com/tvgelderen/fiscora/types"
)

func (h *Handler) HandleGetMe(c echo.Context) error {
	logger := getLogger(c)
	id := c.Get(UserIdCtxKey)
	if id == nil {
		return c.String(http.StatusBadRequest, "User id missing from request context")
	}

	user, err := h.UserRepository.GetById(c.Request().Context(), id.(uuid.UUID))
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error getting user db: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, types.ToUser(user))
}

package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/api/context"
	"github.com/tvgelderen/fiscora/api/handler"
	"github.com/tvgelderen/fiscora/internal/auth"
)

func AuthorizeEndpoint(h *handler.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := auth.GetId(c.Request())
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			_, err = h.UserRepository.GetById(c.Request().Context(), id)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set(context.UserIdCtxKey, id)

			return next(c)
		}
	}
}

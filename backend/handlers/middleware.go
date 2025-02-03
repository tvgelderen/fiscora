package handlers

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/auth"
)

const UserIdCtxKey = "user_id"
const RequestIdCtxKey = "request_id"
const LoggerCtxKey = "logger"

func (h *Handler) AuthorizeEndpoint(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := auth.GetId(c.Request())
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		_, err = h.UserRepository.GetById(c.Request().Context(), id)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		c.Set(UserIdCtxKey, id)

		return next(c)
	}
}

func AttachLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.NewString()

		logger := slog.Default()

		logger = logger.With(
			slog.String("request_id", requestId),
			slog.String("method", c.Request().Method),
			slog.String("uri", c.Request().RequestURI),
			slog.String("url_host", c.Request().URL.Host),
			slog.String("remote_ip", c.RealIP()),
		)

		c.Set(RequestIdCtxKey, requestId)
		c.Set(LoggerCtxKey, logger)

		return next(c)
	}
}

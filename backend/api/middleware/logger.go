package middleware

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/api/context"
)

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

		c.Set(context.RequestIdCtxKey, requestId)
		c.Set(context.LoggerCtxKey, logger)

		return next(c)
	}
}

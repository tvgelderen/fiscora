package router

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tvgelderen/fiscora/api/context"
	"github.com/tvgelderen/fiscora/api/handler"
	customMiddleware "github.com/tvgelderen/fiscora/api/middleware"
	"github.com/tvgelderen/fiscora/internal/auth"
	"github.com/tvgelderen/fiscora/internal/config"
)

func New(conn *sql.DB) *echo.Echo {
	logger := slog.Default()

	authService := auth.NewAuthService()
	h := handler.NewHandler(conn, authService)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Use(customMiddleware.AttachLogger)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "SUCCESS",
					slog.String("request_id", c.Get(context.RequestIdCtxKey).(string)),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("remote_ip", v.RemoteIP),
				)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "ERROR",
					slog.String("request_id", c.Get(context.RequestIdCtxKey).(string)),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("error", v.Error.Error()),
				)
			}
			return nil
		},
	}))

	e.Use(echoprometheus.NewMiddleware("fiscora"))

	go func() {
		metrics := echo.New()
		metrics.GET("/metrics", echoprometheus.NewHandler())
		if err := metrics.Start(config.Env.PrometheusPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
		}
	}()

	base := e.Group("/_api")
	base.GET("/auth/demo", h.HandleDemoLogin)
	base.GET("/auth/:provider", h.HandleOAuthLogin)
	base.GET("/auth/callback/:provider", h.HandleOAuthCallback)
	base.GET("/auth/logout", h.HandleLogout, customMiddleware.AuthorizeEndpoint(h))

	users := base.Group("/users", customMiddleware.AuthorizeEndpoint(h))
	users.GET("/me", h.HandleGetMe)

	transactions := base.Group("/transactions", customMiddleware.AuthorizeEndpoint(h))
	transactions.GET("", h.HandleGetTransactions)
	transactions.POST("", h.HandleCreateTransaction)
	transactions.PUT("/:id", h.HandleUpdateTransaction)
	transactions.DELETE("/:id", h.HandleDeleteTransaction)
	transactions.DELETE("/:id/budget", h.HandleRemoveTransactionFromBudget)
	transactions.GET("/unassigned", h.HandleGetUnassignedTransactions)
	transactions.GET("/types/intervals", h.HandleGetTransactionIntervals)
	transactions.GET("/types/income", h.HandleGetIncomeTypes)
	transactions.GET("/types/expense", h.HandleGetExpenseTypes)
	transactions.GET("/summary/month", h.HandleGetTransactionMonthInfo)
	transactions.GET("/summary/month/type", h.HandleGetTransactionsPerType)
	transactions.GET("/summary/year", h.HandleGetTransactionYearInfo)
	transactions.GET("/summary/year/type", h.HandleGetTransactionsYearInfoPerType)

	budgets := base.Group("/budgets", customMiddleware.AuthorizeEndpoint(h))
	budgets.GET("", h.HandleGetBudgets)
	budgets.POST("", h.HandleCreateBudget)
	budgets.GET("/:id", h.HandleGetBudget)
	budgets.PUT("/:id", h.HandleUpdateBudget)
	budgets.DELETE("/:id", h.HandleDeleteBudget)
	budgets.DELETE("/:id/expenses/:expense_id", h.HandleDeleteBudgetExpense)
	budgets.POST("/:id/expenses/:expense_id/transactions", h.HandleAddBudgetTransactions)

	return e
}

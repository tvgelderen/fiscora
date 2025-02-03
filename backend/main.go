package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/tvgelderen/fiscora/auth"
	"github.com/tvgelderen/fiscora/config"
	"github.com/tvgelderen/fiscora/handler"
	"github.com/tvgelderen/fiscora/logging"
	"github.com/tvgelderen/fiscora/seed"
)

func main() {
	logger, err := logging.SetupLogger()
	if err != nil {
		panic(fmt.Sprintf("Error setting up logger: %v", err.Error()))
	}

	slog.SetDefault(logger)

	env := config.Env
	if env.DBConnectionString == "" {
		log.Fatalf("No database connection string found")
	}

	conn, err := sql.Open("postgres", env.DBConnectionString)
	if err != nil {
		log.Fatalf("Error establishing database connection: %s", err.Error())
	}

	seed.CheckSeed(conn)

	authService := auth.NewAuthService()
	h := handler.NewHandler(conn, authService)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Use(handler.AttachLogger)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "SUCCESS",
					slog.String("request_id", c.Get(handler.RequestIdCtxKey).(string)),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("remote_ip", v.RemoteIP),
				)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "ERROR",
					slog.String("request_id", c.Get(handler.RequestIdCtxKey).(string)),
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
		if err := metrics.Start(env.PrometheusPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())
		}
	}()

	base := e.Group("/_api")
	base.GET("/auth/demo", h.HandleDemoLogin)
	base.GET("/auth/:provider", h.HandleOAuthLogin)
	base.GET("/auth/callback/:provider", h.HandleOAuthCallback)
	base.GET("/auth/logout", h.HandleLogout, h.AuthorizeEndpoint)

	users := base.Group("/users", h.AuthorizeEndpoint)
	users.GET("/me", h.HandleGetMe)

	transactions := base.Group("/transactions", h.AuthorizeEndpoint)
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

	budgets := base.Group("/budgets", h.AuthorizeEndpoint)
	budgets.GET("", h.HandleGetBudgets)
	budgets.POST("", h.HandleCreateBudget)
	budgets.GET("/:id", h.HandleGetBudget)
	budgets.PUT("/:id", h.HandleUpdateBudget)
	budgets.DELETE("/:id", h.HandleDeleteBudget)
	budgets.DELETE("/:id/expenses/:expense_id", h.HandleDeleteBudgetExpense)
	budgets.POST("/:id/expenses/:expense_id/transactions", h.HandleAddBudgetTransactions)

	e.Logger.Fatal(e.Start(env.Port))
}

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/tvgelderen/fiscora/auth"
	"github.com/tvgelderen/fiscora/config"
	"github.com/tvgelderen/fiscora/handlers"
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

	seedFlag := flag.Bool("seed", false, "Set to true to seed the demo account")
	seedMyAccountFlag := flag.Bool("seed-me", false, "Set to true to seed my account too")
	flag.Parse()
	if *seedFlag {
		seed.Seed(conn)
	}
	if *seedMyAccountFlag {
		seed.SeedMyAccount(conn)
	}

	authService := auth.NewAuthService()
	handler := handlers.NewHandler(conn, authService)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Use(handlers.AttachLogger)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(c.Request().Context(), slog.LevelInfo, "SUCCESS",
					slog.String("request_id", c.Get(handlers.RequestIdCtxKey).(string)),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("remote_ip", v.RemoteIP),
				)
			} else {
				logger.LogAttrs(c.Request().Context(), slog.LevelError, "ERROR",
					slog.String("request_id", c.Get(handlers.RequestIdCtxKey).(string)),
					slog.String("remote_ip", v.RemoteIP),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
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
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()

	base := e.Group("/api")
	base.GET("/auth/demo", handler.HandleDemoLogin)
	base.GET("/auth/:provider", handler.HandleOAuthLogin)
	base.GET("/auth/callback/:provider", handler.HandleOAuthCallback)
	base.GET("/auth/logout", handler.HandleLogout, handler.AuthorizeEndpoint)

	users := base.Group("/users", handler.AuthorizeEndpoint)
	users.GET("/me", handler.HandleGetMe)

	transactions := base.Group("/transactions", handler.AuthorizeEndpoint)
	transactions.GET("", handler.HandleGetTransactions)
	transactions.POST("", handler.HandleCreateTransaction)
	transactions.PUT("/:id", handler.HandleUpdateTransaction)
	transactions.DELETE("/:id", handler.HandleDeleteTransaction)
	transactions.DELETE("/:id/budget", handler.HandleRemoveTransactionFromBudget)
	transactions.GET("/unassigned", handler.HandleGetUnassignedTransactions)
	transactions.GET("/types/intervals", handler.HandleGetTransactionIntervals)
	transactions.GET("/types/income", handler.HandleGetIncomeTypes)
	transactions.GET("/types/expense", handler.HandleGetExpenseTypes)
	transactions.GET("/summary/month", handler.HandleGetTransactionMonthInfo)
	transactions.GET("/summary/month/type", handler.HandleGetTransactionsPerType)
	transactions.GET("/summary/year", handler.HandleGetTransactionYearInfo)
	transactions.GET("/summary/year/type", handler.HandleGetTransactionsYearInfoPerType)

	budgets := base.Group("/budgets", handler.AuthorizeEndpoint)
	budgets.GET("", handler.HandleGetBudgets)
	budgets.POST("", handler.HandleCreateBudget)
	budgets.GET("/:id", handler.HandleGetBudget)
	budgets.PUT("/:id", handler.HandleUpdateBudget)
	budgets.DELETE("/:id", handler.HandleDeleteBudget)
	budgets.DELETE("/:id/expenses/:expense_id", handler.HandleDeleteBudgetExpense)
	budgets.POST("/:id/expenses/:expense_id/transactions", handler.HandleAddBudgetTransactions)

	e.Logger.Fatal(e.Start(env.Port))
}

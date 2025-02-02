package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"

	"github.com/tvgelderen/fiscora/auth"
	"github.com/tvgelderen/fiscora/config"
	"github.com/tvgelderen/fiscora/handlers"
	"github.com/tvgelderen/fiscora/seed"
)

func main() {
	env := config.Envs
	if env.DBConnectionString == "" {
		log.Fatal("No database connection string found")
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
	handler := handlers.NewAPIHandler(conn, authService)

	e := echo.New()

	e.Use(echoprometheus.NewMiddleware("fiscora"))

	go func() {
		metrics := echo.New()
		metrics.GET("/metrics", echoprometheus.NewHandler())
		if err := metrics.Start(env.PrometheusPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	e.Use(middleware.Logger())

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

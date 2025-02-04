package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"

	_ "github.com/lib/pq"

	"github.com/tvgelderen/fiscora/api/router"
	"github.com/tvgelderen/fiscora/internal/config"
	"github.com/tvgelderen/fiscora/internal/logging"
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

	e := router.New(conn)

	e.Logger.Fatal(e.Start(env.Port))
}

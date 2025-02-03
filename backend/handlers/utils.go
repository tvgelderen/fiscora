package handlers

import (
	"log/slog"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/types"
)

func getLogger(c echo.Context) *slog.Logger {
	value := c.Get(LoggerCtxKey)
	if value == nil {
		return slog.Default()
	}

	logger, ok := value.(*slog.Logger)
	if !ok {
		slog.Error("Error parsing logger for the current handler")
		return slog.Default()
	}

	return logger
}

func getUserId(c echo.Context) uuid.UUID {
	return c.Get(UserIdCtxKey).(uuid.UUID)
}

func getMonth(c echo.Context) int {
	monthParam := c.QueryParam("month")
	month, err := strconv.ParseInt(monthParam, 10, 16)
	if err != nil {
		month = int64(time.Now().Month())
	}

	return int(month)
}

func getYear(c echo.Context) int {
	yearParam := c.QueryParam("year")
	year, err := strconv.ParseInt(yearParam, 10, 16)
	if err != nil {
		year = int64(time.Now().Month())
	}

	return int(year)
}

func getStartDate(c echo.Context) (time.Time, error) {
	startDate := c.QueryParam("startDate")
	return time.Parse("2006-01-02", startDate)
}

func getEndDate(c echo.Context) (time.Time, error) {
	endDate := c.QueryParam("endDate")
	return time.Parse("2006-01-02", endDate)
}

func getMonthRange(month int, year int) types.DateRange {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, -1)
	return types.DateRange{
		Start: start,
		End:   end,
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateRandomString(length int) string {
	str := make([]rune, length)
	for idx := range str {
		str[idx] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(str)
}

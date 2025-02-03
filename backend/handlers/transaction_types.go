package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/repository"
)

func (h *Handler) HandleGetTransactionIntervals(c echo.Context) error {
	intervals := make([]string, len(repository.TransactionIntervals))

	for idx, interval := range repository.TransactionIntervals {
		intervals[idx] = interval
	}

	return c.JSON(http.StatusOK, intervals)
}

func (h *Handler) HandleGetIncomeTypes(c echo.Context) error {
	incomeTypes := make([]string, len(repository.IncomeTypes))

	for idx, incomeType := range repository.IncomeTypes {
		incomeTypes[idx] = incomeType
	}

	return c.JSON(http.StatusOK, incomeTypes)
}

func (h *Handler) HandleGetExpenseTypes(c echo.Context) error {
	expenseTypes := make([]string, len(repository.ExpenseTypes))

	for idx, expenseType := range repository.ExpenseTypes {
		expenseTypes[idx] = expenseType
	}

	return c.JSON(http.StatusOK, expenseTypes)
}

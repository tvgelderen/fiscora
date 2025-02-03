package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tvgelderen/fiscora/repository"
	"github.com/tvgelderen/fiscora/types"
)

func (h *Handler) HandleGetTransactions(c echo.Context) error {
	logger := getLogger(c)
	userId := getUserId(c)
	month := getMonth(c)
	year := getYear(c)
	dateRange := getMonthRange(month, year)

	var transactions *[]repository.FullTransaction

	income, err := strconv.ParseBool(c.QueryParam("income"))
	params := repository.GetBetweenDatesParams{
		UserID: userId,
		Start:  dateRange.Start,
		End:    dateRange.End,
	}
	if err != nil {
		transactions, err = h.TransactionRepository.GetBetweenDates(c.Request().Context(), params)
	} else if income {
		transactions, err = h.TransactionRepository.GetIncomeBetweenDates(c.Request().Context(), params)
	} else {
		transactions, err = h.TransactionRepository.GetExpenseBetweenDates(c.Request().Context(), params)
	}
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error getting transactions from db: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, types.ToTransactionReturns(transactions))
}

func (h *Handler) HandleGetUnassignedTransactions(c echo.Context) error {
	logger := getLogger(c)
	userId := getUserId(c)
	startDate, startDateErr := getStartDate(c)
	endDate, endDateErr := getEndDate(c)

	if startDateErr != nil || endDateErr != nil {
		return c.String(http.StatusBadRequest, "Invalid date format")
	}

	params := repository.GetBetweenDatesParams{
		UserID: userId,
		Start:  startDate,
		End:    endDate,
	}

	transactions, err := h.TransactionRepository.GetUnassignedBetweenDates(c.Request().Context(), params)
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error getting transactions from db: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.JSON(http.StatusOK, types.ToBaseTransactionReturns(transactions))
}

func (h *Handler) HandleCreateTransaction(c echo.Context) error {
	logger := getLogger(c)

	decoder := json.NewDecoder(c.Request().Body)
	transaction := types.TransactionForm{}
	err := decoder.Decode(&transaction)
	if err != nil {
		logger.Error(fmt.Sprintf("Error decoding request body: %v", err.Error()))
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	// TODO: validate transaction object

	userId := getUserId(c)

	if transaction.Recurring {
		err = h.TransactionRepository.AddRecurring(c.Request().Context(), repository.AddRecurringParams{
			Params: repository.CreateRecurringTransactionParams{
				UserID:       userId,
				StartDate:    transaction.StartDate.Time,
				EndDate:      transaction.EndDate.Time,
				Interval:     transaction.Interval.String,
				DaysInterval: transaction.DaysInterval.NullInt32,
			},
			Amount:      transaction.Amount,
			Description: transaction.Description,
			Type:        transaction.Type,
		})
	} else {
		_, err = h.TransactionRepository.Add(c.Request().Context(), repository.CreateTransactionParams{
			UserID:      userId,
			Amount:      strconv.FormatFloat(transaction.Amount, 'f', -1, 64),
			Description: transaction.Description,
			Type:        transaction.Type,
			Date:        transaction.StartDate.Time,
		})
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.String(http.StatusCreated, "Transaction created successfully")
}

func (h *Handler) HandleUpdateTransaction(c echo.Context) error {
	logger := getLogger(c)

	decoder := json.NewDecoder(c.Request().Body)
	transactionForm := types.TransactionForm{}
	err := decoder.Decode(&transactionForm)
	if err != nil {
		logger.Error(fmt.Sprintf("Error decoding request body: %v", err.Error()))
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	// TODO: validate transaction object

	userId := getUserId(c)
	transactionId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing transaction id from request: %v", err.Error()))
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	transaction, err := h.TransactionRepository.GetById(c.Request().Context(), userId, int32(transactionId))
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error updating transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	if transaction.RecurringTransactionID.Valid {
		err = h.TransactionRepository.UpdateRecurring(c.Request().Context(), repository.UpdateRecurringParams{
			Params: repository.UpdateRecurringTransactionParams{
				ID:           transaction.RecurringTransactionID.Int32,
				UserID:       userId,
				StartDate:    transactionForm.StartDate.Time,
				EndDate:      transactionForm.EndDate.Time,
				Interval:     transactionForm.Interval.String,
				DaysInterval: transactionForm.DaysInterval.NullInt32,
			},
			Amount:      transactionForm.Amount,
			Description: transactionForm.Description,
			Type:        transactionForm.Type,
		})
		if err != nil {
			logger.Error(fmt.Sprintf("Error updating recurring transaction: %v", err.Error()))
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		return c.NoContent(http.StatusNoContent)
	}

	err = h.TransactionRepository.Update(c.Request().Context(), repository.UpdateTransactionParams{
		ID:          int32(transactionId),
		UserID:      userId,
		Date:        transactionForm.StartDate.Time,
		Amount:      strconv.FormatFloat(transactionForm.Amount, 'f', -1, 64),
		Description: transactionForm.Description,
		Type:        transactionForm.Type,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Error updating transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) HandleRemoveTransactionFromBudget(c echo.Context) error {
	logger := getLogger(c)
	userId := getUserId(c)

	transactionId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing transaction id from request: %v", err.Error()))
		return c.NoContent(http.StatusBadRequest)
	}

	err = h.TransactionRepository.RemoveBudgetId(c.Request().Context(), userId, int32(transactionId))
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error deleting transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) HandleDeleteTransaction(c echo.Context) error {
	logger := getLogger(c)
	userId := getUserId(c)

	transactionId, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		logger.Error(fmt.Sprintf("Error parsing transaction id from request: %v", err.Error()))
		return c.NoContent(http.StatusBadRequest)
	}

	transaction, err := h.TransactionRepository.GetById(c.Request().Context(), userId, int32(transactionId))
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		logger.Error(fmt.Sprintf("Error deleting transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	if transaction.RecurringTransactionID.Valid {
		err = h.TransactionRepository.RemoveRecurring(c.Request().Context(), userId, transaction.RecurringTransactionID.Int32)
		if err != nil {
			if repository.NoRowsFound(err) {
				return c.NoContent(http.StatusNotFound)
			}
			logger.Error(fmt.Sprintf("Error deleting transaction: %v", err.Error()))
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}

		return c.NoContent(http.StatusNoContent)
	}

	err = h.TransactionRepository.Remove(c.Request().Context(), userId, int32(transactionId))
	if err != nil {
		logger.Error(fmt.Sprintf("Error deleting transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}

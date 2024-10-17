package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/tvgelderen/fiscora/repository"
	"github.com/tvgelderen/fiscora/types"
)

func (h *APIHandler) HandleGetTransactions(c echo.Context) error {
	userId := getUserId(c)
	month := getMonth(c)
	year := getYear(c)
	dateRange := getMonthRange(month, year)

	var transactions *[]repository.FullTransaction

	income, err := strconv.ParseBool(c.QueryParam("income"))
	if err != nil {
		transactions, err = h.TransactionRepository.GetBetweenDates(c.Request().Context(), userId, dateRange.Start, dateRange.End)
	} else if income {
		transactions, err = h.TransactionRepository.GetIncomeBetweenDates(c.Request().Context(), userId, dateRange.Start, dateRange.End)
	} else {
		transactions, err = h.TransactionRepository.GetExpenseBetweenDates(c.Request().Context(), userId, dateRange.Start, dateRange.End)
	}
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Error(fmt.Sprintf("Error getting transactions from db: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	returnTransactions := make([]types.TransactionReturn, len(*transactions))
	for idx, tranaction := range *transactions {
		returnTransactions[idx] = types.ToReturnTransaction(tranaction)
	}

	return c.JSON(http.StatusOK, returnTransactions)
}

func (h *APIHandler) HandleCreateTransaction(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	transaction := types.TransactionCreateRequest{}
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Errorf("Error decoding request body: %v", err.Error())
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	// TODO: validate transaction object

	userId := getUserId(c)

	_, err = h.TransactionRepository.Add(c.Request().Context(), repository.CreateTransactionParams{
		UserID:      userId,
		Description: transaction.Description,
		Amount:      strconv.FormatFloat(transaction.Amount, 'f', -1, 64),
		Type:        transaction.Type,
	})
	if err != nil {
		log.Error(fmt.Sprintf("Error creating transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.String(http.StatusCreated, "Transaction created successfully")
}

func (h *APIHandler) HandleUpdateTransaction(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	transaction := types.TransactionUpdateRequest{}
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Errorf("Error decoding request body: %v", err.Error())
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	// TODO: validate transaction object

	userId := getUserId(c)
	transactionIdParam := c.Param("id")
	transactionId, err := strconv.ParseInt(transactionIdParam, 10, 32)
	if err != nil {
		log.Errorf("Error parsing transaction id from request: %v", err.Error())
		return c.String(http.StatusBadRequest, "Error decoding request body")
	}

	err = h.TransactionRepository.Update(c.Request().Context(), repository.UpdateTransactionParams{
		ID:          int32(transactionId),
		UserID:      userId,
		Amount:      strconv.FormatFloat(transaction.Amount, 'f', -1, 64),
		Description: transaction.Description,
		Type:        transaction.Type,
	})
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Error(fmt.Sprintf("Error updating transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(204)
}

func (h *APIHandler) HandleDeleteTransaction(c echo.Context) error {
	userId := getUserId(c)
	transactionIdParam := c.Param("id")

	transactionId, err := strconv.ParseInt(transactionIdParam, 10, 32)
	if err != nil {
		log.Errorf("Error parsing transaction id from request: %v", err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	err = h.TransactionRepository.Remove(c.Request().Context(), int32(transactionId), userId)
	if err != nil {
		if repository.NoRowsFound(err) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Error(fmt.Sprintf("Error deleting transaction: %v", err.Error()))
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(204)
}

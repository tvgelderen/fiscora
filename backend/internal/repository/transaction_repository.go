package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type ITransactionRepository interface {
	GetById(ctx context.Context, userId uuid.UUID, id int32) (*Transaction, error)
	GetByBudgetId(ctx context.Context, userId uuid.UUID, budgetId string) (*[]FullTransaction, error)

	GetUnassignedBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]Transaction, error)

	GetBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error)
	GetIncomeBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error)
	GetExpenseBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error)

	GetAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]float64, error)
	GetIncomeAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]TypeAmount, error)
	GetExpenseAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]TypeAmount, error)

	Add(ctx context.Context, params CreateTransactionParams) (*Transaction, error)
	Update(ctx context.Context, params UpdateTransactionParams) error
	UpdateBudgetId(ctx context.Context, params UpdateTransactionBudgetIdParams) error
	Remove(ctx context.Context, userId uuid.UUID, id int32) error
	RemoveBudgetId(ctx context.Context, userId uuid.UUID, id int32) error

	AddRecurring(ctx context.Context, params AddRecurringParams) error
	UpdateRecurring(ctx context.Context, params UpdateRecurringParams) error
	RemoveRecurring(ctx context.Context, userId uuid.UUID, id int32) error
}

type TransactionRepository struct {
	db *sql.DB
}

func CreateTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (repository *TransactionRepository) GetById(ctx context.Context, userId uuid.UUID, id int32) (*Transaction, error) {
	db := New(repository.db)
	transaction, err := db.GetTransactionById(ctx, GetTransactionByIdParams{
		UserID: userId,
		ID:     id,
	})
	return &transaction, err
}

func (repository *TransactionRepository) GetByBudgetId(ctx context.Context, userId uuid.UUID, budgetId string) (*[]FullTransaction, error) {
	db := New(repository.db)
	transactions, err := db.GetTransactionsByBudgetId(ctx, GetTransactionsByBudgetIdParams{
		UserID:   userId,
		BudgetID: budgetId,
	})
	if err != nil {
		return nil, err
	}

	fullTransactions := make([]FullTransaction, len(transactions))
	for idx, transaction := range transactions {
		fullTransactions[idx] = transaction.FullTransaction
	}

	return &fullTransactions, nil
}

type GetBetweenDatesParams struct {
	UserID uuid.UUID
	Start  time.Time
	End    time.Time
}

func (repository *TransactionRepository) GetUnassignedBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]Transaction, error) {
	db := New(repository.db)
	transactions, err := db.GetUnassignedTransactionsBetweenDates(ctx, GetUnassignedTransactionsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
		Limit:     MaxFetchLimit,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	return &transactions, err
}

func (repository *TransactionRepository) GetBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error) {
	db := New(repository.db)
	transactions, err := db.GetTransactionsBetweenDates(ctx, GetTransactionsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
		Limit:     MaxFetchLimit,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	fullTransactions := make([]FullTransaction, len(transactions))
	for idx, transaction := range transactions {
		fullTransactions[idx] = transaction.FullTransaction
	}

	return &fullTransactions, err
}

func (repository *TransactionRepository) GetIncomeBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error) {
	db := New(repository.db)
	transactions, err := db.GetIncomeTransactionsBetweenDates(ctx, GetIncomeTransactionsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
		Limit:     MaxFetchLimit,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	fullTransactions := make([]FullTransaction, len(transactions))
	for idx, transaction := range transactions {
		fullTransactions[idx] = transaction.FullTransaction
	}

	return &fullTransactions, err
}

func (repository *TransactionRepository) GetExpenseBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]FullTransaction, error) {
	db := New(repository.db)
	transactions, err := db.GetExpenseTransactionsBetweenDates(ctx, GetExpenseTransactionsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
		Limit:     MaxFetchLimit,
		Offset:    0,
	})
	if err != nil {
		return nil, err
	}

	fullTransactions := make([]FullTransaction, len(transactions))
	for idx, transaction := range transactions {
		fullTransactions[idx] = transaction.FullTransaction
	}

	return &fullTransactions, err
}

func (repository *TransactionRepository) GetAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]float64, error) {
	db := New(repository.db)
	amounts, err := db.GetTransactionAmountsBetweenDates(ctx, GetTransactionAmountsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
	})
	if err != nil {
		return nil, err
	}

	floats := make([]float64, len(amounts))
	for idx, typeAmount := range amounts {
		value, err := strconv.ParseFloat(typeAmount, 64)
		if err != nil {
			value = 0
		}

		floats[idx] = value
	}

	return &floats, nil
}

func (repository *TransactionRepository) GetIncomeAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]TypeAmount, error) {
	db := New(repository.db)
	typeAmounts, err := db.GetIncomeTransactionAmountsBetweenDates(ctx, GetIncomeTransactionAmountsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
	})
	if err != nil {
		return nil, err
	}

	returnValues := make([]TypeAmount, len(typeAmounts))
	for idx, typeAmount := range typeAmounts {
		value, err := strconv.ParseFloat(typeAmount.Amount, 64)
		if err != nil {
			value = 0
		}

		returnValues[idx] = TypeAmount{
			Type:   typeAmount.Type,
			Amount: value,
		}
	}

	return &returnValues, nil
}

func (repository *TransactionRepository) GetExpenseAmountsBetweenDates(ctx context.Context, params GetBetweenDatesParams) (*[]TypeAmount, error) {
	db := New(repository.db)
	typeAmounts, err := db.GetExpenseTransactionAmountsBetweenDates(ctx, GetExpenseTransactionAmountsBetweenDatesParams{
		UserID:    params.UserID,
		StartDate: params.Start,
		EndDate:   params.End,
	})
	if err != nil {
		return nil, err
	}

	returnValues := make([]TypeAmount, len(typeAmounts))
	for idx, typeAmount := range typeAmounts {
		value, err := strconv.ParseFloat(typeAmount.Amount, 64)
		if err != nil {
			value = 0
		}

		returnValues[idx] = TypeAmount{
			Type:   typeAmount.Type,
			Amount: value,
		}
	}

	return &returnValues, nil
}

func (repository *TransactionRepository) Add(ctx context.Context, params CreateTransactionParams) (*Transaction, error) {
	db := New(repository.db)
	transaction, err := db.CreateTransaction(ctx, params)
	return &transaction, err
}

func (repository *TransactionRepository) Update(ctx context.Context, params UpdateTransactionParams) error {
	db := New(repository.db)
	return db.UpdateTransaction(ctx, params)
}

func (repository *TransactionRepository) UpdateBudgetId(ctx context.Context, params UpdateTransactionBudgetIdParams) error {
	db := New(repository.db)
	return db.UpdateTransactionBudgetId(ctx, params)
}

func (repository *TransactionRepository) Remove(ctx context.Context, userId uuid.UUID, id int32) error {
	db := New(repository.db)
	return db.DeleteTransaction(ctx, DeleteTransactionParams{
		ID:     id,
		UserID: userId,
	})
}

func (repository *TransactionRepository) RemoveBudgetId(ctx context.Context, userId uuid.UUID, id int32) error {
	db := New(repository.db)
	return db.RemoveTransactionBudgetId(ctx, RemoveTransactionBudgetIdParams{
		ID:     id,
		UserID: userId,
	})
}

type AddRecurringParams struct {
	Params      CreateRecurringTransactionParams
	Amount      float64
	Description string
	Type        string
}

func (repository *TransactionRepository) AddRecurring(ctx context.Context, params AddRecurringParams) error {
	db := New(repository.db)
	recurringTransaction, err := db.CreateRecurringTransaction(ctx, params.Params)
	if err != nil {
		return err
	}

	amountString := strconv.FormatFloat(params.Amount, 'f', -1, 64)
	createParams := getRecurringTransactions(getRecurringTransactionsParams{
		UserID:                 params.Params.UserID,
		RecurringTransactionId: recurringTransaction.ID,
		Description:            params.Description,
		Amount:                 amountString,
		Type:                   params.Type,
		StartDate:              recurringTransaction.StartDate,
		EndDate:                recurringTransaction.EndDate,
		Interval:               recurringTransaction.Interval,
		DaysInterval:           recurringTransaction.DaysInterval.Int32,
	})

	for _, params := range createParams {
		_, err := repository.Add(ctx, params)
		if err != nil {
			return err
		}
	}

	return err
}

type UpdateRecurringParams struct {
	Params      UpdateRecurringTransactionParams
	Amount      float64
	Description string
	Type        string
}

func (repository *TransactionRepository) UpdateRecurring(ctx context.Context, params UpdateRecurringParams) error {
	db := New(repository.db)

	userId := params.Params.UserID
	recurringTransaction, err := db.GetRecurringTransactionById(ctx, GetRecurringTransactionByIdParams{
		UserID: userId,
		ID:     params.Params.ID,
	})
	if err != nil {
		return err
	}

	amountString := strconv.FormatFloat(params.Amount, 'f', -1, 64)
	startDate := params.Params.StartDate
	endDate := params.Params.EndDate
	interval := params.Params.Interval
	daysInterval := params.Params.DaysInterval

	if recurringTransaction.Interval != interval || recurringTransaction.StartDate.UTC() != startDate {
		_, err := db.DeleteTransactionsByRecurringTransactionId(ctx, DeleteTransactionsByRecurringTransactionIdParams{
			UserID:                 userId,
			RecurringTransactionID: recurringTransaction.ID,
		})
		if err != nil {
			return fmt.Errorf("error deleting recurring transactions: %v", err.Error())
		}

		createParams := getRecurringTransactions(getRecurringTransactionsParams{
			UserID:                 userId,
			RecurringTransactionId: recurringTransaction.ID,
			Description:            params.Description,
			Amount:                 amountString,
			Type:                   params.Type,
			StartDate:              startDate,
			EndDate:                endDate,
			Interval:               interval,
			DaysInterval:           daysInterval.Int32,
		})

		for _, params := range createParams {
			_, err := repository.Add(ctx, params)
			if err != nil {
				return fmt.Errorf("error creating transaction: %v", err.Error())
			}
		}
	} else {
		transactions, err := db.GetTransactionsByRecurringTransactionId(ctx, GetTransactionsByRecurringTransactionIdParams{
			UserID:                 userId,
			RecurringTransactionID: params.Params.ID,
		})
		if err != nil {
			return fmt.Errorf("error getting recurring transactions: %v", err.Error())
		}

		if endDate.After(recurringTransaction.EndDate) {
			createParams := getRecurringTransactions(getRecurringTransactionsParams{
				UserID:                 userId,
				RecurringTransactionId: recurringTransaction.ID,
				Description:            params.Description,
				Amount:                 amountString,
				Type:                   params.Type,
				StartDate:              transactions[len(transactions)-1].Date,
				EndDate:                endDate,
				Interval:               interval,
				DaysInterval:           daysInterval.Int32,
			})

			for idx := 1; idx < len(createParams); idx++ {
				_, err := repository.Add(ctx, createParams[idx])
				if err != nil {
					return fmt.Errorf("error creating transaction: %v", err.Error())
				}
			}
		} else if endDate.Before(recurringTransaction.EndDate) {
			_, err := db.DeleteTransactionsByRecurringTransactionIdAndWhereDate(ctx, DeleteTransactionsByRecurringTransactionIdAndWhereDateParams{
				UserID:                 userId,
				RecurringTransactionID: recurringTransaction.ID,
				Date:                   endDate,
			})
			if err != nil {
				return fmt.Errorf("error deleting transactions after new end date: %v", err.Error())
			}

			for idx := len(transactions) - 1; idx >= 0; idx-- {
				if transactions[idx].Date.After(endDate) {
					err := repository.Remove(ctx, userId, transactions[idx].ID)
					if err != nil {
						return fmt.Errorf("error deleting transaction: %v", err.Error())
					}
				} else {
					break
				}
			}
		}

		if transactions[0].Description != params.Description || transactions[0].Amount != amountString || transactions[0].Type != params.Type {
			for _, transaction := range transactions {
				err := db.UpdateTransaction(ctx, UpdateTransactionParams{
					ID:          transaction.ID,
					UserID:      userId,
					Amount:      amountString,
					Description: params.Description,
					Type:        params.Type,
					Date:        transaction.Date,
				})
				if err != nil {
					return fmt.Errorf("error updating transaction: %v", err.Error())
				}
			}
		}
	}

	return db.UpdateRecurringTransaction(ctx, params.Params)
}

func (repository *TransactionRepository) RemoveRecurring(ctx context.Context, userId uuid.UUID, id int32) error {
	db := New(repository.db)
	return db.DeleteRecurringTransaction(ctx, DeleteRecurringTransactionParams{
		ID:     id,
		UserID: userId,
	})
}

type getRecurringTransactionsParams struct {
	UserID                 uuid.UUID
	RecurringTransactionId int32
	Description            string
	Amount                 string
	Type                   string
	StartDate              time.Time
	EndDate                time.Time
	Interval               string
	DaysInterval           int32
}

func getRecurringTransactions(params getRecurringTransactionsParams) []CreateTransactionParams {
	date := params.StartDate
	endDate := params.EndDate

	var createParams []CreateTransactionParams
	for date.Before(endDate) {
		createParams = append(createParams, CreateTransactionParams{
			UserID:                 params.UserID,
			RecurringTransactionID: sql.NullInt32{Valid: true, Int32: params.RecurringTransactionId},
			Amount:                 params.Amount,
			Description:            params.Description,
			Type:                   params.Type,
			Date:                   date,
		})

		switch params.Interval {
		case TransactionIntervalDaily:
			date = addDays(date, 1)
		case TransactionIntervalWeekly:
			date = addWeeks(date, 1)
		case TransactionIntervalMonthly:
			date = addMonths(date, 1)
		case TransactionIntervalOther:
			date = addDays(date, int(params.DaysInterval))
		}
	}

	return createParams
}

func addDays(date time.Time, days int) time.Time {
	return date.AddDate(0, 0, days)
}

func addWeeks(date time.Time, weeks int) time.Time {
	return date.AddDate(0, 0, weeks*7)
}

func addMonths(date time.Time, months int) time.Time {
	return date.AddDate(0, months, 0)
}

type TypeAmount struct {
	Type   string
	Amount float64
}

// Transaction interval
const (
	TransactionIntervalDaily   string = "Daily"
	TransactionIntervalWeekly         = "Weekly"
	TransactionIntervalMonthly        = "Monthly"
	TransactionIntervalOther          = "Other"
)

var TransactionIntervals = []string{
	TransactionIntervalDaily,
	TransactionIntervalWeekly,
	TransactionIntervalMonthly,
	TransactionIntervalOther,
}

// Income type
const (
	IncomeTypeSalary            string = "Salary"
	IncomeTypePassive                  = "Passive"
	IncomeTypeCapitalGains             = "Capital Gains"
	IncomeTypeDividend                 = "Dividend"
	IncomeTypeGovernmentPayment        = "Government Payment"
	IncomeTypeOther                    = "Other"
)

var IncomeTypes = []string{
	IncomeTypeSalary,
	IncomeTypePassive,
	IncomeTypeCapitalGains,
	IncomeTypeDividend,
	IncomeTypeGovernmentPayment,
	IncomeTypeOther,
}

// Expense type
const (
	ExpenseTypeMortgage      string = "Mortgage"
	ExpenseTypeRent                 = "Rent"
	ExpenseTypeUtilities            = "Utilities"
	ExpenseTypeFixed                = "Fixed"
	ExpenseTypeGroceries            = "Groceries"
	ExpenseTypeInsurance            = "Insurance"
	ExpenseTypeTravel               = "Travel"
	ExpenseTypeTaxes                = "Taxes"
	ExpenseTypeInterest             = "Interest"
	ExpenseTypeSubscriptions        = "Subscriptions"
	ExpenseTypeOther                = "Other"
)

var ExpenseTypes = []string{
	ExpenseTypeMortgage,
	ExpenseTypeRent,
	ExpenseTypeUtilities,
	ExpenseTypeFixed,
	ExpenseTypeGroceries,
	ExpenseTypeInsurance,
	ExpenseTypeTravel,
	ExpenseTypeTaxes,
	ExpenseTypeInterest,
	ExpenseTypeSubscriptions,
	ExpenseTypeOther,
}

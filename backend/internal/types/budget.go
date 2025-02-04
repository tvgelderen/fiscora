package types

import (
	"strconv"
	"time"

	"github.com/tvgelderen/fiscora/internal/repository"
)

type BaseBudget struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
}

type BudgetCreateRequest struct {
	BaseBudget
}

type BudgetUpdateRequest struct {
	BaseBudget
}

type BudgetReturn struct {
	BaseBudget
	ID           string                 `json:"id"`
	Created      time.Time              `json:"created"`
	Updated      time.Time              `json:"updated"`
	Expenses     *[]BudgetExpenseReturn `json:"expenses"`
	Transactions *[]TransactionReturn   `json:"transactions"`
}

type BudgetForm struct {
	BaseBudget
	Expenses []BudgetExpenseForm `json:"expenses"`
}

type BaseBudgetExpense struct {
	Name            string  `json:"name"`
	AllocatedAmount float64 `json:"allocatedAmount"`
	CurrentAmount   float64 `json:"currentAmount"`
}

type BudgetExpenseCreateRequest struct {
	BaseBudgetExpense
}

type BudgetExpenseUpdateRequest struct {
	BaseBudgetExpense
}

type BudgetExpenseForm struct {
	BaseBudgetExpense
	ID int32 `json:"id"`
}

type BudgetExpenseReturn struct {
	BaseBudgetExpense
	ID int32 `json:"id"`
}

func ToBudgetReturn(budget *repository.BudgetWithExpenses) BudgetReturn {
	amount, _ := strconv.ParseFloat(budget.Amount, 64)

	expenses := make([]BudgetExpenseReturn, len(budget.Expenses))
	for idx, expense := range budget.Expenses {
		expenses[idx] = ToBudgetExpenseReturn(&expense)
	}

	return BudgetReturn{
		ID:           budget.ID,
		Created:      budget.Created,
		Updated:      budget.Updated,
		Expenses:     &expenses,
		Transactions: nil,
		BaseBudget: BaseBudget{
			Name:        budget.Name,
			Description: budget.Description,
			Amount:      amount,
			StartDate:   budget.StartDate,
			EndDate:     budget.EndDate,
		},
	}
}

func ToBudgetExpenseReturn(expense *repository.BudgetExpense) BudgetExpenseReturn {
	allocatedAmount, _ := strconv.ParseFloat(expense.AllocatedAmount, 64)
	currentAmount, _ := strconv.ParseFloat(expense.CurrentAmount, 64)

	return BudgetExpenseReturn{
		ID: expense.ID,
		BaseBudgetExpense: BaseBudgetExpense{
			Name:            expense.Name,
			AllocatedAmount: allocatedAmount,
			CurrentAmount:   currentAmount,
		},
	}
}

const (
	BudgetTypeWeekly  string = "Weekly"
	BudgetTypeMonthly        = "Monthly"
	BudgetTypeYearly         = "Yearly"
	BudgetTypeCustom         = "Custom"
)

var BudgetTypes = []string{
	BudgetTypeWeekly,
	BudgetTypeMonthly,
	BudgetTypeYearly,
	BudgetTypeCustom,
}

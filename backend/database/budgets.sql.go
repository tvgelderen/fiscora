// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: budgets.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createBudget = `-- name: CreateBudget :one
INSERT INTO budgets (id, user_id, name, description, amount, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, name, description, amount, start_date, end_date, created, updated
`

type CreateBudgetParams struct {
	ID          string
	UserID      uuid.UUID
	Name        string
	Description string
	Amount      string
	StartDate   time.Time
	EndDate     time.Time
}

func (q *Queries) CreateBudget(ctx context.Context, arg CreateBudgetParams) (Budget, error) {
	row := q.db.QueryRowContext(ctx, createBudget,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.Description,
		arg.Amount,
		arg.StartDate,
		arg.EndDate,
	)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Amount,
		&i.StartDate,
		&i.EndDate,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const createBudgetExpense = `-- name: CreateBudgetExpense :one
INSERT INTO budget_expenses (budget_id, name, allocated_amount)
VALUES ($1, $2, $3)
RETURNING id, budget_id, name, allocated_amount, current_amount, created, updated
`

type CreateBudgetExpenseParams struct {
	BudgetID        string
	Name            string
	AllocatedAmount string
}

func (q *Queries) CreateBudgetExpense(ctx context.Context, arg CreateBudgetExpenseParams) (BudgetExpense, error) {
	row := q.db.QueryRowContext(ctx, createBudgetExpense, arg.BudgetID, arg.Name, arg.AllocatedAmount)
	var i BudgetExpense
	err := row.Scan(
		&i.ID,
		&i.BudgetID,
		&i.Name,
		&i.AllocatedAmount,
		&i.CurrentAmount,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const deleteBudget = `-- name: DeleteBudget :exec
DELETE FROM budgets 
WHERE id = $1 AND user_id = $2
`

type DeleteBudgetParams struct {
	ID     string
	UserID uuid.UUID
}

func (q *Queries) DeleteBudget(ctx context.Context, arg DeleteBudgetParams) error {
	_, err := q.db.ExecContext(ctx, deleteBudget, arg.ID, arg.UserID)
	return err
}

const deleteBudgetExpense = `-- name: DeleteBudgetExpense :exec
DELETE FROM budget_expenses
WHERE id = $1 AND budget_id = $2
`

type DeleteBudgetExpenseParams struct {
	ID       int32
	BudgetID string
}

func (q *Queries) DeleteBudgetExpense(ctx context.Context, arg DeleteBudgetExpenseParams) error {
	_, err := q.db.ExecContext(ctx, deleteBudgetExpense, arg.ID, arg.BudgetID)
	return err
}

const getBudget = `-- name: GetBudget :one
SELECT id, user_id, name, description, amount, start_date, end_date, created, updated FROM budgets
WHERE id = $1
`

func (q *Queries) GetBudget(ctx context.Context, id string) (Budget, error) {
	row := q.db.QueryRowContext(ctx, getBudget, id)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Amount,
		&i.StartDate,
		&i.EndDate,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const getBudgetExpenses = `-- name: GetBudgetExpenses :many
SELECT id, budget_id, name, allocated_amount, current_amount, created, updated FROM budget_expenses
WHERE budget_id = $1
`

func (q *Queries) GetBudgetExpenses(ctx context.Context, budgetID string) ([]BudgetExpense, error) {
	rows, err := q.db.QueryContext(ctx, getBudgetExpenses, budgetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BudgetExpense
	for rows.Next() {
		var i BudgetExpense
		if err := rows.Scan(
			&i.ID,
			&i.BudgetID,
			&i.Name,
			&i.AllocatedAmount,
			&i.CurrentAmount,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBudgets = `-- name: GetBudgets :many
SELECT id, user_id, name, description, amount, start_date, end_date, created, updated FROM budgets
WHERE user_id = $1
ORDER BY created DESC
LIMIT $2
OFFSET $3
`

type GetBudgetsParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) GetBudgets(ctx context.Context, arg GetBudgetsParams) ([]Budget, error) {
	rows, err := q.db.QueryContext(ctx, getBudgets, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Budget
	for rows.Next() {
		var i Budget
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.Amount,
			&i.StartDate,
			&i.EndDate,
			&i.Created,
			&i.Updated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBudgetsExpenses = `-- name: GetBudgetsExpenses :many
SELECT e.id, e.budget_id, e.name, e.allocated_amount, e.current_amount FROM budgets b JOIN budget_expenses e ON b.id = e.budget_id
WHERE b.user_id = $1
LIMIT $2
OFFSET $3
`

type GetBudgetsExpensesParams struct {
	UserID uuid.UUID
	Limit  int32
	Offset int32
}

type GetBudgetsExpensesRow struct {
	ID              int32
	BudgetID        string
	Name            string
	AllocatedAmount string
	CurrentAmount   string
}

func (q *Queries) GetBudgetsExpenses(ctx context.Context, arg GetBudgetsExpensesParams) ([]GetBudgetsExpensesRow, error) {
	rows, err := q.db.QueryContext(ctx, getBudgetsExpenses, arg.UserID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBudgetsExpensesRow
	for rows.Next() {
		var i GetBudgetsExpensesRow
		if err := rows.Scan(
			&i.ID,
			&i.BudgetID,
			&i.Name,
			&i.AllocatedAmount,
			&i.CurrentAmount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBudget = `-- name: UpdateBudget :one
UPDATE budgets
SET name = $3, description = $4, amount = $5, start_date = $6, end_date = $7, updated = (now() at time zone 'utc')
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, name, description, amount, start_date, end_date, created, updated
`

type UpdateBudgetParams struct {
	ID          string
	UserID      uuid.UUID
	Name        string
	Description string
	Amount      string
	StartDate   time.Time
	EndDate     time.Time
}

func (q *Queries) UpdateBudget(ctx context.Context, arg UpdateBudgetParams) (Budget, error) {
	row := q.db.QueryRowContext(ctx, updateBudget,
		arg.ID,
		arg.UserID,
		arg.Name,
		arg.Description,
		arg.Amount,
		arg.StartDate,
		arg.EndDate,
	)
	var i Budget
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Amount,
		&i.StartDate,
		&i.EndDate,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

const updateBudgetExpense = `-- name: UpdateBudgetExpense :one
UPDATE budget_expenses 
SET name = $2, allocated_amount = $3, current_amount = $4
WHERE id = $1
RETURNING id, budget_id, name, allocated_amount, current_amount, created, updated
`

type UpdateBudgetExpenseParams struct {
	ID              int32
	Name            string
	AllocatedAmount string
	CurrentAmount   string
}

func (q *Queries) UpdateBudgetExpense(ctx context.Context, arg UpdateBudgetExpenseParams) (BudgetExpense, error) {
	row := q.db.QueryRowContext(ctx, updateBudgetExpense,
		arg.ID,
		arg.Name,
		arg.AllocatedAmount,
		arg.CurrentAmount,
	)
	var i BudgetExpense
	err := row.Scan(
		&i.ID,
		&i.BudgetID,
		&i.Name,
		&i.AllocatedAmount,
		&i.CurrentAmount,
		&i.Created,
		&i.Updated,
	)
	return i, err
}

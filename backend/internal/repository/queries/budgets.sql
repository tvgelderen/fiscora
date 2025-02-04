-- name: CreateBudget :one
INSERT INTO budgets (id, user_id, name, description, amount, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateBudget :one
UPDATE budgets
SET name = $3, description = $4, amount = $5, start_date = $6, end_date = $7, updated = (now() at time zone 'utc')
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: GetBudgets :many
SELECT * FROM budgets
WHERE user_id = $1
ORDER BY created DESC
LIMIT $2
OFFSET $3;

-- name: GetBudget :one
SELECT * FROM budgets
WHERE id = $1 AND user_id = $2;

-- name: GetBudgetsExpenses :many
SELECT sqlc.embed(be) FROM budgets b JOIN budget_expenses be ON b.id = be.budget_id
WHERE b.user_id = $1
LIMIT $2
OFFSET $3;

-- name: DeleteBudget :exec
DELETE FROM budgets 
WHERE id = $1 AND user_id = $2;


-- name: CreateBudgetExpense :one
INSERT INTO budget_expenses (budget_id, name, allocated_amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateBudgetExpense :one
UPDATE budget_expenses 
SET name = $2, allocated_amount = $3
WHERE id = $1
RETURNING *;

-- name: GetBudgetExpenses :many
SELECT * FROM budget_expenses
WHERE budget_id = $1;

-- name: DeleteBudgetExpense :exec
DELETE FROM budget_expenses
WHERE id = $1 AND budget_id = $2;

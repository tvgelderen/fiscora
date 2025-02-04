-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    provider VARCHAR(16) NOT NULL,
    provider_id VARCHAR(32) NOT NULL,
    username VARCHAR(32) NOT NULL,
    email VARCHAR(64) NOT NULL,
    avatar VARCHAR(128),
    created TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    UNIQUE(provider, provider_id),
    UNIQUE(provider, username),
    UNIQUE(provider, email)
);

CREATE TABLE IF NOT EXISTS recurring_transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    interval VARCHAR(16) NOT NULL,
    days_interval INTEGER,
    created TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS budgets (
    id VARCHAR(16) PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(256) NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS budget_expenses (
    id SERIAL PRIMARY KEY,
    budget_id VARCHAR(16) NOT NULL,
    name VARCHAR(64) NOT NULL,
    allocated_amount DECIMAL(19, 2) NOT NULL,
    current_amount DECIMAL(19, 2) NOT NULL DEFAULT 0,
    created TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    FOREIGN KEY (budget_id) REFERENCES budgets(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    budget_id VARCHAR(16) DEFAULT NULL,
    budget_expense_id INT DEFAULT NULL,
    recurring_transaction_id INT DEFAULT NULL,
    description VARCHAR(512) NOT NULL,
    amount DECIMAL(19, 4) NOT NULL,
    type VARCHAR(32) NOT NULL,
    date TIMESTAMP NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),
    updated TIMESTAMP NOT NULL DEFAULT (now() at time zone 'utc'),

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(budget_id) REFERENCES budgets(id),
    FOREIGN KEY(budget_expense_id) REFERENCES budget_expenses(id),
    FOREIGN KEY(recurring_transaction_id) REFERENCES recurring_transactions(id) ON DELETE CASCADE
);

CREATE VIEW full_transaction AS (
    SELECT t.*, rt.start_date, rt.end_date, rt.interval, rt.days_interval, rt.created as recurring_created, rt.updated as recurring_updated, b.name as budget_name, be.name as budget_expense_name
    FROM transactions t 
        LEFT OUTER JOIN recurring_transactions rt ON t.recurring_transaction_id = rt.id 
        LEFT OUTER JOIN budget_expenses be ON t.budget_expense_id = be.id 
        LEFT JOIN budgets b ON be.budget_id = b.id
);

-- +goose Down
DROP VIEW full_transaction;

DROP TABLE transactions;
DROP TABLE budget_expenses;
DROP TABLE budgets;
DROP TABLE recurring_transactions;
DROP TABLE users;

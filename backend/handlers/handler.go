package handlers

import (
	"database/sql"

	"github.com/tvgelderen/fiscora/auth"
	"github.com/tvgelderen/fiscora/repository"
)

type Handler struct {
	UserRepository        repository.IUserRepository
	TransactionRepository repository.ITransactionRepository
	BudgetRepository      repository.IBudgetRepository
	AuthService           *auth.AuthService
}

func NewHandler(db *sql.DB, auth *auth.AuthService) *Handler {
	return &Handler{
		UserRepository:        repository.CreateUserRepository(db),
		TransactionRepository: repository.CreateTransactionRepository(db),
		BudgetRepository:      repository.CreateBudgetRepository(db),
		AuthService:           auth,
	}
}

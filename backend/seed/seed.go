package seed

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/tvgelderen/fiscora/repository"
	"github.com/tvgelderen/fiscora/types"
)

var userId uuid.UUID

func Seed(conn *sql.DB) {
	slog.Info("Seeding repository.")

	userRepository := repository.CreateUserRepository(conn)
	transactionRepository := repository.CreateTransactionRepository(conn)
	budgetRepository := repository.CreateBudgetRepository(conn)

	userId, _ = uuid.NewUUID()

	slog.Info("Creating demo user.")
	createDemoUser(userRepository)

	slog.Info("Creating demo user transactions.")
	createTransactions(transactionRepository)

	slog.Info("Creating demo user budgets.")
	createBudgets(budgetRepository)
}

func SeedMyAccount(conn *sql.DB) {
	slog.Info("Seeding repository.")

	userRepository := repository.CreateUserRepository(conn)
	transactionRepository := repository.CreateTransactionRepository(conn)
	budgetRepository := repository.CreateBudgetRepository(conn)

	user, err := userRepository.GetByEmail(context.Background(), "thvangelderen@gmail.com")
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting user from db: ", err.Error()))
	}

	userId = user.ID

	slog.Info("Creating transactions.")
	createTransactions(transactionRepository)

	slog.Info("Creating bugets.")
	createBudgets(budgetRepository)
}

func createDemoUser(userRepository *repository.UserRepository) {
	_, _ = userRepository.Add(context.Background(), repository.CreateUserParams{
		ID:         userId,
		Provider:   "demo",
		ProviderID: "demo",
		Username:   "demo",
		Email:      "demo",
	})
}

func createTransactions(transactionRepository *repository.TransactionRepository) {
	for _, recurringTransaction := range recurringTransactions {
		_ = transactionRepository.AddRecurring(context.Background(), repository.AddRecurringParams{
			Params: repository.CreateRecurringTransactionParams{
				UserID:       userId,
				StartDate:    recurringTransaction.StartDate.Time,
				EndDate:      recurringTransaction.EndDate.Time,
				Interval:     recurringTransaction.Interval.String,
				DaysInterval: recurringTransaction.DaysInterval.NullInt32,
			},
			Amount:      recurringTransaction.Amount,
			Description: recurringTransaction.Description,
			Type:        recurringTransaction.Type,
		})
	}

	for _, transaction := range transactions {
		transaction.UserID = userId
		_, _ = transactionRepository.Add(context.Background(), transaction)
	}
}

func createBudgets(repository *repository.BudgetRepository) {
	for _, budget := range budgets {
		budget.UserID = userId
		_, _ = repository.Add(context.Background(), budget)
	}

	for _, budgetExpense := range budgetExpenses {
		_, _ = repository.AddExpense(context.Background(), budgetExpense)
	}
}

func randomTime() time.Time {
	return time.Now().UTC().AddDate(0, 0, -rand.Intn(120))
}

var budgetId1 = "4HYmJAkwkEcBWNlb"
var budgetId2 = "nlbNr8xaGAQNEvNX"

var budgets = []repository.CreateBudgetParams{
	{
		ID:          budgetId1,
		Name:        "Monthly Budget",
		Description: "Overall monthly budget for household expenses",
		Amount:      "2024",
		StartDate:   time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 10, 31, 0, 0, 0, 0, time.UTC),
	},
	{
		ID:          budgetId2,
		Name:        "Vacation Budget",
		Description: "Saving for summer vacation",
		Amount:      "1250",
		StartDate:   time.Date(2024, 8, 7, 0, 0, 0, 0, time.UTC),
		EndDate:     time.Date(2024, 8, 21, 0, 0, 0, 0, time.UTC),
	},
}

var budgetExpenses = []repository.CreateBudgetExpenseParams{
	{
		BudgetID:        budgetId1,
		Name:            "Groceries",
		AllocatedAmount: "500",
	},
	{
		BudgetID:        budgetId1,
		Name:            "Utilities",
		AllocatedAmount: "300",
	},
	{
		BudgetID:        budgetId1,
		Name:            "Entertainment",
		AllocatedAmount: "100",
	},
	{
		BudgetID:        budgetId1,
		Name:            "Savings",
		AllocatedAmount: "500",
	},
	{
		BudgetID:        budgetId2,
		Name:            "Accomodation",
		AllocatedAmount: "600",
	},
	{
		BudgetID:        budgetId2,
		Name:            "Transportation",
		AllocatedAmount: "400",
	},
	{
		BudgetID:        budgetId2,
		Name:            "Activities",
		AllocatedAmount: "300",
	},
	{
		BudgetID:        budgetId2,
		Name:            "Food",
		AllocatedAmount: "200",
	},
}

var recurringTransactions = []types.TransactionForm{
	// Salary (Monthly, starting July 25, 2024)
	{
		Amount:       3789,
		Description:  "Salary",
		Type:         repository.IncomeTypeSalary,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 25, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Weekly income (Starting May 12, 2024)
	{
		Amount:       128,
		Description:  "Weekly income",
		Type:         repository.IncomeTypePassive,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 5, 12, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2025, 4, 12, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalWeekly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Mortgage (Monthly, starting January 1, 2024)
	{
		Amount:       -632,
		Description:  "Mortgage",
		Type:         repository.ExpenseTypeMortgage,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 1, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Utilities (Monthly, starting January 1, 2024)
	{
		Amount:       -215,
		Description:  "Utilities",
		Type:         repository.ExpenseTypeUtilities,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 1, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Internet (Monthly, starting January 6, 2024)
	{
		Amount:       -64,
		Description:  "Internet",
		Type:         repository.ExpenseTypeUtilities,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 6, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// HOA (Monthly, starting January 1, 2024)
	{
		Amount:       -224,
		Description:  "HOA",
		Type:         repository.ExpenseTypeFixed,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 1, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Health insurance (Monthly, starting January 10, 2024)
	{
		Amount:       -124,
		Description:  "Health insurance",
		Type:         repository.ExpenseTypeInsurance,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 10, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	// Various subscriptions (Monthly, starting January 10, 2024)
	{
		Amount:       -54,
		Description:  "Streaming services bundle",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 10, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	{
		Amount:       -12.99,
		Description:  "Music streaming service",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 15, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	{
		Amount:       -9.99,
		Description:  "Cloud storage",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 5, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	{
		Amount:       -29.99,
		Description:  "Online fitness membership",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 1, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	{
		Amount:       -9.99,
		Description:  "Software license",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 1, 15, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
	{
		Amount:       -4.99,
		Description:  "App subscription",
		Type:         repository.ExpenseTypeSubscriptions,
		StartDate:    types.NewNullTimeFromTime(time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)),
		EndDate:      types.NewNullTimeFromTime(time.Date(2032, 12, 20, 0, 0, 0, 0, time.UTC)),
		Interval:     types.NewNullStringFromString(repository.TransactionIntervalMonthly),
		DaysInterval: types.NewNullIntFromInt(0),
	},
}

var transactions = []repository.CreateTransactionParams{
	// Groceries (Non-recurring, spread throughout the year)
	{
		UserID:      userId,
		Amount:      "-120",
		Description: "New Year groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-85",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-100",
		Description: "Monthly stock up",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 1, 29, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-130",
		Description: "Super Bowl party supplies",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-123",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 2, 18, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-79",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 2, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-115",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 3, 8, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-80",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 3, 14, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-95",
		Description: "Spring sale groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 3, 22, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-110",
		Description: "Easter holiday groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 4, 5, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-125",
		Description: "Groceries for BBQ party",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 4, 19, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-96",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 4, 27, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-130",
		Description: "Cinco de Mayo groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 5, 3, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-85",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 5, 17, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-115",
		Description: "Memorial Day groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 5, 30, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-102",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 6, 7, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-90",
		Description: "Groceries for June",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 6, 12, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-95",
		Description: "Father's Day groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 6, 21, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-115",
		Description: "Fourth of July groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 7, 3, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-110",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 7, 18, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-79",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 7, 28, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-150",
		Description: "Back to school groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 8, 5, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-100",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 8, 18, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-78",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 8, 28, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-120",
		Description: "Labor Day groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 9, 1, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-65",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 9, 24, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-130",
		Description: "Fall season groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 9, 15, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-89",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 10, 02, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-125",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 10, 14, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-140",
		Description: "Halloween groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 10, 27, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-95",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 11, 11, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-160",
		Description: "Thanksgiving groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 11, 23, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-86",
		Description: "Groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 12, 01, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-135",
		Description: "Weekly groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 12, 10, 0, 0, 0, 0, time.UTC),
	},
	{
		UserID:      userId,
		Amount:      "-165",
		Description: "Christmas holiday groceries",
		Type:        repository.ExpenseTypeGroceries,
		Date:        time.Date(2024, 12, 24, 0, 0, 0, 0, time.UTC),
	},
}

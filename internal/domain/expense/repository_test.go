package expense

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := "postgres://postgres:***@localhost:5432/mumix_test?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		year VARCHAR(10) NOT NULL,
		month VARCHAR(20) NOT NULL,
		income_main BIGINT DEFAULT 0,
		income_other BIGINT DEFAULT 0,
		transfer_to_nda BIGINT DEFAULT 0,
		transfer_to_mimi BIGINT DEFAULT 0,
		credit_card BIGINT DEFAULT 0,
		fixed_savings BIGINT DEFAULT 0,
		other_expense BIGINT DEFAULT 0,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP
	)`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	_, _ = db.Exec("DELETE FROM expenses")
	return db
}

func TestParseRupiah(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Rp0", 0},
		{"Rp8.800.000", 8800000},
		{"Rp1.000", 1000},
		{"-Rp26.000", -26000},
		{"(Rp26.000)", -26000},
		{"", 0},
		{"Rp 5.000.000", 5000000},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseRupiah(tt.input)
			if result != tt.expected {
				t.Errorf("parseRupiah(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatRupiah(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "Rp0"},
		{8800000, "Rp8.800.000"},
		{-26000, "-Rp26.000"},
		{1000, "Rp1.000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatRupiah(tt.input)
			if result != tt.expected {
				t.Errorf("formatRupiah(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpenseCRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	e := &Expense{
		Year:           "2026",
		Month:          "Januari",
		IncomeMain:     10000000,
		IncomeOther:    2000000,
		TransferToNda:  3000000,
		TransferToMimi: 2000000,
		CreditCard:     1000000,
		FixedSavings:   2000000,
		OtherExpense:   500000,
	}

	created, err := repo.Create(ctx, e)
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}
	if created.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	expenses, err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed to find all: %v", err)
	}
	if len(expenses) != 1 {
		t.Errorf("expected 1 expense, got %d", len(expenses))
	}

	created.IncomeMain = 12000000
	updated, err := repo.Update(ctx, created.ID, created)
	if err != nil {
		t.Fatalf("failed to update: %v", err)
	}
	if updated.IncomeMain != 12000000 {
		t.Errorf("expected income_main 12000000, got %d", updated.IncomeMain)
	}

	err = repo.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	expenses, _ = repo.FindAll(ctx)
	if len(expenses) != 0 {
		t.Errorf("expected 0 expenses after delete, got %d", len(expenses))
	}
}

func TestGetTotals(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	_, _ = repo.Create(ctx, &Expense{
		Year: "2026", Month: "Januari",
		IncomeMain: 10000000, IncomeOther: 2000000,
		TransferToNda: 3000000, FixedSavings: 2000000,
	})
	_, _ = repo.Create(ctx, &Expense{
		Year: "2026", Month: "Februari",
		IncomeMain: 10000000, IncomeOther: 1000000,
		CreditCard: 1500000, OtherExpense: 500000,
	})

	totals, err := repo.GetTotals(ctx)
	if err != nil {
		t.Fatalf("failed to get totals: %v", err)
	}

	if totals["pendapatan_utama"] != 20000000 {
		t.Errorf("expected pendapatan_utama 20000000, got %d", totals["pendapatan_utama"])
	}
	if totals["fixed_savings"] != 2000000 {
		t.Errorf("expected fixed_savings 2000000, got %d", totals["fixed_savings"])
	}
}

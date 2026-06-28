-- +goose Up
CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    year VARCHAR(4) NOT NULL,
    month VARCHAR(20) NOT NULL,
    income_main INT NOT NULL DEFAULT 0,
    income_other INT NOT NULL DEFAULT 0,
    transfer_to_nda INT NOT NULL DEFAULT 0,
    transfer_to_mimi INT NOT NULL DEFAULT 0,
    credit_card INT NOT NULL DEFAULT 0,
    fixed_savings INT NOT NULL DEFAULT 0,
    other_expense INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_expenses_year_month ON expenses(year, month);

-- +goose Down
DROP TABLE IF EXISTS expenses;

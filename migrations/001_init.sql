-- +goose Up

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    password TEXT NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    token_version INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    completed BOOLEAN DEFAULT false,
    priority VARCHAR(20) DEFAULT 'medium',
    due_date DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
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
);

CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE OR REPLACE FUNCTION month_order(month_name VARCHAR)
RETURNS INTEGER AS $$
BEGIN
    RETURN CASE month_name
        WHEN 'Januari' THEN 1
        WHEN 'Februari' THEN 2
        WHEN 'Maret' THEN 3
        WHEN 'April' THEN 4
        WHEN 'Mei' THEN 5
        WHEN 'Juni' THEN 6
        WHEN 'Juli' THEN 7
        WHEN 'Agustus' THEN 8
        WHEN 'September' THEN 9
        WHEN 'Oktober' THEN 10
        WHEN 'November' THEN 11
        WHEN 'Desember' THEN 12
        ELSE 0
    END;
END;
$$ LANGUAGE plpgsql;

-- +goose Down
DROP FUNCTION IF EXISTS month_order(VARCHAR);
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;

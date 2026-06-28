# Mumix Backend

Go backend for Mumix application.

## Quick Start

\`\`\`bash
go mod download
cp .env.example .env
go run ./cmd/server
\`\`\`

## API Endpoints

### Auth (Public)
- `POST /api/auth/register` - Register
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout

### Todos (Authenticated)
- `GET /api/todos` - List (paginated)
- `POST /api/todos` - Create
- `GET /api/todos/:id` - Detail
- `PUT /api/todos/:id` - Update
- `PATCH /api/todos/:id` - Partial update
- `DELETE /api/todos/:id` - Delete

### Expenses (Authenticated)
- `GET /api/expenses` - List
- `POST /api/expenses` - Create
- `GET /api/expenses/:id` - Detail
- `PUT /api/expenses/:id` - Update
- `DELETE /api/expenses/:id` - Delete
- `GET /api/expenses/totals` - Totals
- `GET /api/expenses/export` - Export CSV

### Users (Admin)
- `GET /api/users` - List
- `POST /api/users` - Create
- `GET /api/users/:id` - Detail
- `PUT /api/users/:id` - Update
- `DELETE /api/users/:id` - Delete

## Tech Stack

- Go 1.22
- PostgreSQL (pgx/v5)
- JWT authentication
- Standard library net/http

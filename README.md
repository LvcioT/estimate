# Scrum Planning Poker (estimate)

Simple Planning Poker web app using Go, Gin, HTMX and SQLite.

## Configuration

The application can be configured using environment variables or a `.env` file. Copy `.env.example` to `.env` and adjust the values:

```bash
cp .env.example .env
```

Available configuration options:

- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (debug/release)
- `DB_PATH`: SQLite database file path
- `SESSION_SECRET`: Secret key for session cookies
- `HASH_COST`: BCrypt hash cost (default: 10)
- `DEFAULT_ADMIN_USERNAME`: Default admin username when no users exist
- `DEFAULT_ADMIN_PASSWORD`: Default admin password when no users exist

## Running

1. Install dependencies:
```bash
go mod tidy
```

2. Configure the application:
```bash
cp .env.example .env
# Edit .env with your settings
```

3. Run the server:
```bash
go run ./
```

The application will be available at http://localhost:8080 (or the configured port).

**Important**: In production:
- Set `GIN_MODE=release`
- Use a strong `SESSION_SECRET`
- Change the default admin password
- Back up your `estimate.db` file regularly
# estimate
Scrum Planning Poker

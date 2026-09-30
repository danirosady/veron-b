.PHONY: help migrate-up migrate-down migrate-reset migrate-status db-console

# Load .env variables
include .env
export

# Connection string
PSQL := psql "host=$(DB_HOST) port=$(DB_PORT) dbname=$(DB_NAME) user=$(DB_USER) password=$(DB_PASSWORD) sslmode=$(DB_SSLMODE)"

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "  migrate-up        Run pending migrations (auto on server start)"
	@echo "  migrate-down     Rollback last migration"
	@echo "  migrate-reset    Drop all tables & re-run from scratch"
	@echo "  migrate-status   Show migration status"
	@echo "  db-console       Open psql console"

# Run pending migrations (called automatically by server start)
migrate-up:
	@echo "Running migrations..."
	@go run cmd/server/main.go --migrate-up

# Rollback last migration (runs .down.sql)
migrate-down:
	@echo "Rolling back last migration..."
	@go run cmd/server/main.go --migrate-down

# Fresh reset: mark all migrations as not applied + drop all tables
migrate-reset:
	@$(PSQL) -c "TRUNCATE schema_migrations RESTART IDENTITY CASCADE;" 2>/dev/null || true
	@$(PSQL) -c 'DO $dopla$ DECLARE r RECORD; BEGIN FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = '\''public'\'') LOOP EXECUTE '\''DROP TABLE IF EXISTS '\'' || quote_ident(r.tablename) || '\'' CASCADE'\''; END LOOP; END $dopla$;' 2>/dev/null || true
	@echo "Database cleared. Run 'make migrate-up' to re-apply."

# Show migration status
migrate-status:
	@echo "=== Applied Migrations ==="
	@$(PSQL) -c "SELECT version, applied_at FROM schema_migrations ORDER BY applied_at;" 2>/dev/null || echo "schema_migrations table not found"
	@echo ""
	@echo "=== Available Migrations ==="
	@ls migrations/*.up.sql 2>/dev/null | xargs -I{} basename {} .up.sql

# Open psql console
db-console:
	@$(PSQL)

# Development server (migrations auto-run on start)
server:
	go run cmd/server/main.go

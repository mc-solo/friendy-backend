# =============================================================================
# MIGRATION MAKEFILE
# =============================================================================
# This is an AI-Gen Makefile that manages database migrations using golang-migrate.
# It includes commands for creating, applying, rolling back, and fixing dirty
# migration states.
# [perfect for dev exp]
# =============================================================================

# -----------------------------------------------------------------------------
# Config (you can override this config to match your local db)
# -----------------------------------------------------------------------------
DB_DRIVER   ?= postgres
DB_USER     ?= postgres
DB_PASSWORD ?= postgres
DB_HOST     ?= localhost
DB_PORT     ?= 5432
DB_NAME     ?= friendy
DB_SSL      ?= disable

# Build the database URL
DB_URL       = "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL)"

MIGRATIONS_PATH ?= ./migrations
MIGRATE_CMD  ?= migrate

# Colors for pretty output
RED    = \033[0;31m
GREEN  = \033[0;32m
YELLOW = \033[1;33m
BLUE   = \033[0;34m
NC    = \033[0m # No Color

# -----------------------------------------------------------------------------
# Helper to check if migrate is installed
# -----------------------------------------------------------------------------
.PHONY: check-migrate
check-migrate:
	@command -v $(MIGRATE_CMD) >/dev/null 2>&1 || { \
		echo "$(RED)migrate is not installed.$(NC)"; \
		echo "Please install it via:"; \
		echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		exit 1; \
	}

# -----------------------------------------------------------------------------
# Show help
# -----------------------------------------------------------------------------
.PHONY: help
help:
	@echo "$(BLUE)Available commands:$(NC)"
	@echo "  $(GREEN)make create name=<name>$(NC)        Create new migration files (up/down)"
	@echo "  $(GREEN)make up$(NC)                         Apply all pending migrations"
	@echo "  $(GREEN)make down$(NC)                       Rollback the last migration"
	@echo "  $(GREEN)make down-all$(NC)                   Rollback all migrations"
	@echo "  $(GREEN)make reset$(NC)                      Rollback all and apply all (fresh start)"
	@echo "  $(GREEN)make redo$(NC)                       Rollback the last migration and reapply it"
	@echo "  $(GREEN)make version$(NC)                    Show current migration version"
	@echo "  $(GREEN)make status$(NC)                     Show detailed migration status (including dirty flag)"
	@echo "  $(GREEN)make force-version v=<version>$(NC)  Force set migration version (dangerous)"
	@echo "  $(GREEN)make fix-dirty$(NC)                  Help to fix a dirty migration state"
	@echo "  $(GREEN)make db-shell$(NC)                   Open a psql shell to the database"
	@echo "  $(GREEN)make help$(NC)                       Show this help"

# -----------------------------------------------------------------------------
# Create a new migration
# -----------------------------------------------------------------------------
.PHONY: create
create: check-migrate
	@if [ -z "$(name)" ]; then \
		echo "$(RED)Missing migration name. Use: make create name=<name>$(NC)"; \
		exit 1; \
	fi
	@$(MIGRATE_CMD) create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
	@echo "$(GREEN)Created migration: $(name)$(NC)"

# -----------------------------------------------------------------------------
# Apply all pending migrations
# -----------------------------------------------------------------------------
.PHONY: up
up: check-migrate
	@echo "$(BLUE)Applying migrations...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) up
	@echo "$(GREEN)Migrations applied$(NC)"

# -----------------------------------------------------------------------------
# Rollback the last migration
# -----------------------------------------------------------------------------
.PHONY: down
down: check-migrate
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) down 1
	@echo "$(GREEN)Rollback complete$(NC)"

# -----------------------------------------------------------------------------
# Rollback all migrations
# -----------------------------------------------------------------------------
.PHONY: down-all
down-all: check-migrate
	@echo "$(YELLOW)Rolling back ALL migrations...$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) down
	@echo "$(GREEN)All migrations rolled back$(NC)"

# -----------------------------------------------------------------------------
# Reset database: down-all then up
# -----------------------------------------------------------------------------
.PHONY: reset
reset: check-migrate
	@echo "$(YELLOW)Resetting database (down-all + up)...$(NC)"
	@$(MAKE) down-all
	@$(MAKE) up
	@echo "$(GREEN)Database reset complete$(NC)"

# -----------------------------------------------------------------------------
# Redo last migration: down then up
# -----------------------------------------------------------------------------
.PHONY: redo
redo: check-migrate
	@echo "$(YELLOW)Redoing last migration...$(NC)"
	@$(MAKE) down
	@$(MAKE) up
	@echo "$(GREEN)Redo complete$(NC)"

# -----------------------------------------------------------------------------
# Show current migration version
# -----------------------------------------------------------------------------
.PHONY: version
version: check-migrate
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) version

# -----------------------------------------------------------------------------
# Show detailed status (including dirty flag)
# -----------------------------------------------------------------------------
.PHONY: status
status: check-migrate
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) version 2>&1 | grep -q "dirty" && \
		echo "$(RED)Database is in DIRTY state!$(NC)" || \
		echo "$(GREEN)Database is clean.$(NC)"
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) version

# -----------------------------------------------------------------------------
# Force set migration version (use only to recover from dirty state)
# -----------------------------------------------------------------------------
.PHONY: force-version
force-version: check-migrate
	@if [ -z "$(v)" ]; then \
		echo "$(RED)Missing version. Use: make force-version v=<version>$(NC)"; \
		exit 1; \
	fi
	@echo "$(RED)DANGER: Forcing migration version to $(v). This can break your database!$(NC)"
	@read -p "Are you absolutely sure? (type 'yes' to confirm) " confirmation; \
	if [ "$$confirmation" != "yes" ]; then \
		echo "$(YELLOW)Aborted.$(NC)"; \
		exit 1; \
	fi
	@$(MIGRATE_CMD) -path $(MIGRATIONS_PATH) -database $(DB_URL) force $(v)
	@echo "$(GREEN)Version forced to $(v)$(NC)"

# -----------------------------------------------------------------------------
# Help for fixing dirty migration state
# -----------------------------------------------------------------------------
.PHONY: fix-dirty
fix-dirty:
	@echo "$(YELLOW)How to fix a dirty migration:$(NC)"
	@echo ""
	@echo "1. First, check the current version and error:"
	@echo "   $(GREEN)make status$(NC)"
	@echo ""
	@echo "2. Identify the last successful migration version (from your migration files)."
	@echo "   For example, if the dirty version is 3, and you know migration 2 was fine,"
	@echo "   you might want to force the version back to 2."
	@echo ""
	@echo "3. Force the version (replace X with the correct version number):"
	@echo "   $(GREEN)make force-version v=X$(NC)"
	@echo ""
	@echo "4. Then re-run migrations:"
	@echo "   $(GREEN)make up$(NC)"
	@echo ""
	@echo "$(RED)Warning: Forcing a version should only be done if you understand the consequences.$(NC)"
	@echo "   It marks the dirty version as clean without actually applying any changes."

# -----------------------------------------------------------------------------
# Open a psql shell to the database
# -----------------------------------------------------------------------------
.PHONY: db-shell
db-shell:
	@echo "$(BLUE)Connecting to database...$(NC)"
	@psql "$(DB_URL)"

# -----------------------------------------------------------------------------
# (Optional) Drop the database – use with extreme caution!
# -----------------------------------------------------------------------------
# .PHONY: db-drop
# db-drop:
# 	@echo "$(RED)DANGER: This will DROP the database '$(DB_NAME)'.$(NC)"
# 	@read -p "Type the database name to confirm: " confirmation; \
# 	if [ "$$confirmation" != "$(DB_NAME)" ]; then \
# 		echo "$(YELLOW)Aborted.$(NC)"; \
# 		exit 1; \
# 	fi
# 	@dropdb --if-exists $(DB_NAME)
# 	@echo "$(GREEN)Database dropped$(NC)"
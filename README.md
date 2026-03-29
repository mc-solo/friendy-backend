# friendy-backend

**Friendy** is an Ethiopian localized dating and friendship platform; built to connect people across Ethiopia through meaningful relationships. This repository contains the backend infrastructure that powers all Friendy clients: **Web, Android, and iOS**.

Written with the powerful and efficient **Go 1.24** programming language, following a clean layered architecture.

> **🚧 This project is currently under active development.** APIs and features are subject to change.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| HTTP Router | chi v5 |
| ORM | GORM v1.31.1 |
| Database | PostgreSQL 15 |
| Auth | JWT (Access & Refresh tokens), OAuth2, 2FA, Argon2/Bcrypt, and more |
| Config | Viper (`.env` / YAML / env vars) |
| Migrations | golang-migrate v4 |
| Dev tooling | Air (live reload), Docker Compose |

---

## Project Structure

```
.
├── cmd/
│   └── api/v1/         # Application entry point (main.go)
├── internal/
│   ├── app/            # App container, router wiring
│   ├── config/         # Config loading & DB connection
│   ├── database/       # GORM models & migration runner
│   ├── delivery/       # HTTP handlers (controllers)
│   ├── repository/     # Data access layer (stores)
│   ├── service/        # Business logic layer
│   └── utils/          # Shared helpers (token, password)
├── migrations/         # SQL migration files (up/down)
├── docker-compose.yml  # Local PostgreSQL container
├── Makefile            # DB & migration automation
└── .air.toml           # Air live-reload config
```

---

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Air](https://github.com/cosmtrek/air) — `go install github.com/air-verse/air@latest`
- [golang-migrate](https://github.com/golang-migrate/migrate) — `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

### 1. Clone & configure

```bash
git clone https://github.com/mc-solo/friendy.git
cd friendy-backend
cp env.example .env
# Edit .env with your local settings
```

### 2. Start the database

```bash
docker compose up -d
```

### 3. Run migrations

```bash
make up
```

### 4. Start the dev server

```bash
air
```

The API will be available at `http://localhost:8080` (or whichever port is configured).

---

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `friendy` |

---

## Makefile Reference

```bash
make up                       # Apply all pending migrations
make down                     # Roll back the last migration
make down-all                 # Roll back all migrations
make reset                    # down-all + up (fresh start)
make redo                     # Roll back and reapply the last migration
make create name=<name>       # Create a new migration (up/down files)
make version                  # Show current migration version
make status                   # Show migration status + dirty flag
make force-version v=<N>      # Force migration version (recovery only)
make fix-dirty                # Guide to fix a dirty migration state
make db-shell                 # Open a psql shell to the database
```

---

## Architecture

The application follows a **Layered / Clean Architecture**:

```
Client → Router (chi) → Handlers (Delivery) → Services (Business Logic) → Repositories (Data Access) → PostgreSQL
```

For a full breakdown of layers, request lifecycle, and how to add a new feature, see [ARCHITECTURE.md](./ARCHITECTURE.md).

---

## License

MIT

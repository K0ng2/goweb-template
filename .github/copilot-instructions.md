# AI Coding Instructions for GoWeb Template

## Architecture Overview

This is a full-stack Go + Vue 3 (Nuxt 4) web application with a clean separation of concerns:

### Backend (Go)
- **Framework**: Fiber v3 (fast, Express-like HTTP server)
- **Database Layer**: Abstract `Repository` interface with pluggable SQLite/PostgreSQL implementations
- **Query Building**: go-jet ORM for type-safe SQL generation
- **Documentation**: Swagger/OpenAPI via swaggo (generates from handler godoc comments)

### Frontend (Vue 3)
- **Framework**: Nuxt 4 (SSR disabled - SPA mode)
- **Styling**: Tailwind CSS v4 + DaisyUI
- **Icons**: FontAwesome v7
- **Package Manager**: Bun (not npm)

### Build & Task Management
- **Task Runner**: Mise (Go-based task framework in `mise.toml`)
- **Database Migrations**: SQL files in `backend/repo/{sqlite,postgres}/queries/init.sql`

---

## Critical Developer Workflows

### Running the Application
```bash
mise run dev        # Start Nuxt dev server (port 3000)
mise run run-server # Build and run Go server (port 8080)
mise run build-server # Build only
```

### Database & Code Generation
```bash
mise run jet [sqlite|postgres]  # Generate database models from schema
mise run swag                    # Generate Swagger docs from handler comments
```

### Version Management
```bash
mise run version-bump [patch|minor|major|prerelease]  # Bumps version in package.json
```

### Configuration
- Backend config: `backend/config.yml` (DSN and port via Viper)
- Default: SQLite at `./db.sqlite` (edit DSN in config.yml for PostgreSQL)
- Environment: `mise.file = ".env"` loads variables

---

## Key Architectural Patterns

### Database Abstraction
1. **Interface-based design** (`backend/repo/repo.go`): `Repository` interface allows swapping implementations
2. **Per-driver implementations**: `backend/repo/sqlite/` and `backend/repo/postgres/` mirror functionality
3. **Initialization pattern**: `repo.NewRepo(db)` auto-detects driver from DSN prefix and initializes schema
4. **Executor abstraction** (`database/database.go`): Supports transactions via `QueryContext`, `ExecContext`, `QueryRowContext`

**DSN Detection Logic**:
- SQLite: `sqlite://file:./path.db`
- PostgreSQL: `postgres://user:pass@host/db`

### Handler & Dependency Injection
- Handlers receive `*database.Database` → create `Repository` → pass to handler methods
- Pattern: `NewHandler(db) → Handler.repo → Repository methods`
- All handlers use Fiber's `c fiber.Ctx` for requests

### Swagger Documentation
- **Generation source**: `backend/server/server.go` (set as `-g` flag)
- **Handler godoc comments**: Each handler must have Swagger annotations (see `DatabaseHealth` example)
- **Endpoint**: `/api/swagger/*` (auto-configured)
- Command: `swag init -g server/server.go -o ./docs --ot go`

### API Response Format
All responses use `model.APIResponse[T]` generic struct:
```go
APIResponse{
  Data: <T>,
  Meta: &Meta{Total, Limit, Offset},  // optional, for pagination
  Error: "message"  // only if error
}
```

### Web Asset Embedding
- Frontend files: `web/public/` → embedded in binary via `backend/web/embed.go`
- Served statically at `/` (index.html) and `/200.html` (fallback)
- No separate web server needed in production

---

## Project-Specific Conventions

### Go Package Layout
- `cmd/`: CLI commands (if needed)
- `config/`: Configuration loading (Viper-based singleton `config.C`)
- `database/`: Database connection abstraction
- `handler/`: HTTP endpoint handlers
- `model/`: Go structs for JSON serialization (tagged with `json:"field"`)
- `repo/`: Data layer interface + implementations
- `server/`: Fiber app setup and routes
- `docs/`: Auto-generated Swagger files (don't edit)
- `web/`: Frontend asset embedding

### Adding a New API Endpoint
1. Create handler method in `handler/handler.go` with Swagger godoc
2. Register route in `server/server.go`: `api.Get("/path", h.HandlerMethod)`
3. Define request/response models in `model/model.go`
4. Run `mise run swag` to regenerate docs

### Database Queries with go-jet
- Schemas auto-generated from `init.sql` files
- Use generated package for type-safe column/table references
- Import from `github.com/go-jet/jet/v2`

### Frontend Development
- Entry point: `web/app/app.vue` (sets up layouts and routes)
- Routes auto-discovered from `web/app/pages/` (Nuxt file routing)
- Layouts: `web/app/layouts/default.vue` (auto-applied)
- CSS: `web/app/assets/css/main.css` (imported by Tailwind)
- FontAwesome available globally via `fontawesome.ts` plugin

### Configuration Hierarchy
- Viper reads `backend/config.yml` on startup
- Unmarshals into `config.C` global struct
- Environment variables can override (set in `.env` which is loaded by Mise)

---

## Integration Points & External Dependencies

### Backend Dependencies
- **gofiber/fiber/v3**: HTTP framework with middleware
- **go-jet/jet/v2**: Type-safe SQL builder
- **lib/pq**: PostgreSQL driver
- **modernc.org/sqlite**: SQLite driver (pure Go)
- **spf13/viper**: Config management
- **swaggo/swag**: Swagger doc generation

### Frontend Dependencies
- **nuxt**: Meta-framework (handles routing, SSR setup, etc.)
- **tailwindcss**: Utility-first CSS
- **@fortawesome/vue-fontawesome**: Icon library

### Build System
- **Mise**: Cross-platform task runner (replaces Make/npm scripts)
- **go build**: Standard Go compilation
- **nuxt build**: Vite-based bundling
- **swag**: Swagger doc generator (Go tool)
- **jet**: go-jet CLI for schema generation

### Cross-Component Data Flow
1. Frontend (Nuxt) → HTTP request → Fiber server
2. Fiber handler → Repository interface call
3. SQLite/PostgreSQL implementation → Database
4. Result → `model.APIResponse[T]` JSON response
5. Frontend receives and renders

---

## Notes for AI Agents

- **Database schema changes**: Edit `backend/repo/{sqlite,postgres}/queries/init.sql`, then run `mise run jet` to regenerate models
- **New dependencies**: Update `backend/go.mod` or `web/package.json`, then re-run build tasks
- **Fiber middleware**: Added in `server.New()` before route definition (CORS, logging, recovery already configured)
- **Error handling**: Fiber handlers return `fiber.Ctx.Send*()` or JSON via `fiber.Ctx.JSON()`
- **Testing**: No tests found yet - set up test files in `*_test.go` alongside implementation files

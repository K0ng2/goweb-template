# GoWeb Template

A modern, full-stack web application template combining Go backend (Fiber) with Vue 3 frontend (Nuxt 4). Ships with built-in database abstraction, Swagger API docs, and type-safe database queries.

## ✨ Features

- **Go Backend** - Fiber v3 HTTP framework with middleware (CORS, logging, recovery)
- **Vue 3 Frontend** - Nuxt 4 SSR-disabled SPA with file-based routing
- **Database Abstraction** - Pluggable SQLite/PostgreSQL with go-jet type-safe queries
- **API Documentation** - Auto-generated Swagger/OpenAPI docs from handler comments
- **Styling** - Tailwind CSS v4 + DaisyUI component library
- **Icons** - FontAwesome v7 integration
- **Task Automation** - Mise task runner (cross-platform Make replacement)
- **Production Ready** - Single binary deployment with embedded frontend assets

## 🚀 Quick Start

### Prerequisites
- Go 1.25.4+
- Bun (JavaScript runtime)
- Mise (task runner) - [install here](https://mise.jdx.dev/)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd goweb-template

# Install Mise (if not already installed)
curl https://mise.jdx.dev/install.sh | sh

# Run development servers
mise run dev        # Nuxt frontend on http://localhost:3000
mise run run-server # Go backend on http://localhost:8080
```

### Verify Setup

Navigate to `http://localhost:3000` in your browser. The frontend should connect to the backend at `http://localhost:8080`.

## 📁 Project Structure

```
goweb-template/
├── backend/                  # Go application
│   ├── main.go              # Entry point
│   ├── config.yml           # Configuration (DSN, port)
│   ├── go.mod               # Go dependencies
│   ├── config/              # Configuration loading (Viper)
│   ├── database/            # Database connection abstraction
│   ├── handler/             # HTTP endpoint handlers
│   ├── model/               # API request/response models
│   ├── repo/                # Data layer interface
│   │   ├── repo.go          # Repository interface
│   │   ├── sqlite/          # SQLite implementation
│   │   └── postgres/        # PostgreSQL implementation
│   ├── server/              # Fiber app setup & routes
│   ├── docs/                # Auto-generated Swagger docs
│   └── web/                 # Frontend asset embedding
│
├── web/                     # Nuxt application
│   ├── nuxt.config.ts       # Nuxt configuration
│   ├── package.json         # Frontend dependencies
│   ├── tsconfig.json        # TypeScript config
│   ├── public/              # Static assets (embedded in Go binary)
│   └── app/
│       ├── app.vue          # Root component
│       ├── pages/           # Auto-routed pages
│       ├── layouts/         # Reusable layouts
│       ├── assets/          # CSS & images
│       └── plugins/         # Nuxt plugins (FontAwesome, etc.)
│
├── mise.toml                # Task definitions
├── Dockerfile               # Container build
└── README.md                # This file
```

## 🛠️ Common Commands

### Development

```bash
# Start both servers in separate terminals
mise run dev                 # Nuxt dev server (port 3000)
mise run run-server          # Go server (port 8080)

# Only build without running
mise run build-server        # Compile Go binary
```

### Database

```bash
# Generate database models from schema (choose one)
mise run jet sqlite          # SQLite
mise run jet postgres        # PostgreSQL
```

### API Documentation

```bash
# Generate Swagger docs from handler comments
mise run swag
# View at http://localhost:8080/api/swagger/
```

### Deployment

```bash
# Build production binary with optimizations
cd backend
go build -ldflags="-s -w" -v -o dist/app
./dist/app
```

### Versioning

```bash
mise run version-bump patch     # e.g., 1.0.0 → 1.0.1
mise run version-bump minor     # e.g., 1.0.0 → 1.1.0
mise run version-bump major     # e.g., 1.0.0 → 2.0.0
```

## 📊 Architecture

### Data Flow

1. **Frontend Request** - Vue component → HTTP request to `/api/*`
2. **Backend Handler** - Fiber route → Handler function in `backend/handler/`
3. **Repository Pattern** - Handler → `Repository` interface → SQLite/PostgreSQL implementation
4. **Database Query** - go-jet type-safe query → Result
5. **API Response** - `model.APIResponse[T]` → JSON response
6. **Frontend Render** - Vue component displays result

### Database Abstraction

The `Repository` interface (`backend/repo/repo.go`) provides a consistent API regardless of database:

```go
// Implementations in backend/repo/sqlite/ and backend/repo/postgres/
type Repository interface {
    Ping(ctx context.Context) error
    Init(ctx context.Context) error
    // Add your custom methods here
}
```

**Database Detection** (from `backend/config.yml` DSN):
- `sqlite://file:./db.sqlite` → Uses SQLite
- `postgres://user:pass@host/db` → Uses PostgreSQL

### API Response Format

All endpoints return a consistent response structure:

```go
type APIResponse[T any] struct {
    Data  T     `json:"data,omitempty"`
    Meta  *Meta `json:"meta,omitempty"`        // Optional pagination
    Error string `json:"error,omitempty"`       // Only on error
}

type Meta struct {
    Total  int64 `json:"total"`
    Limit  int64 `json:"limit"`
    Offset int64 `json:"offset"`
}
```

**Example response:**
```json
{
  "data": { "status": "healthy", "database": "healthy", "uptime": "2h30m" },
  "meta": null
}
```

## 🔧 Configuration

### Backend (`backend/config.yml`)

```yaml
DSN: "sqlite://file:./db.sqlite"  # Database connection string
Port: ":8080"                      # Server port
```

**Environment overrides** - Set in `.env` (loaded by Mise):
```bash
DSN=postgres://user:pass@localhost/mydb
PORT=:3000
```

### Frontend (`web/nuxt.config.ts`)

```typescript
export default defineNuxtConfig({
  ssr: false,  // SPA mode (no server-side rendering)
  vite: {
    plugins: [tailwindcss()],
  },
})
```

## 📝 Adding API Endpoints

### 1. Define Request/Response Models

Edit `backend/model/model.go`:

```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type GetUserRequest struct {
    ID int `query:"id,default:1"`
}
```

### 2. Create Handler

Edit `backend/handler/handler.go`:

```go
// GetUser godoc
// @Summary Get user by ID
// @Description Fetch a user from database
// @Tags users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} model.APIResponse[model.User]
// @Router /users/{id} [get]
func (h *Handler) GetUser(c fiber.Ctx) error {
    var req model.GetUserRequest
    if err := c.QueryParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.APIResponse[any]{
            Error: "Invalid query params",
        })
    }

    user, err := h.repo.GetUser(c.RequestCtx(), req.ID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.APIResponse[any]{
            Error: err.Error(),
        })
    }

    return c.JSON(model.APIResponse[model.User]{
        Data: user,
    })
}
```

### 3. Add Repository Methods

Implement in both `backend/repo/sqlite/sqlite.go` and `backend/repo/postgres/postgres.go`:

```go
func (r *Sqlite) GetUser(ctx context.Context, id int) (*model.User, error) {
    // Use go-jet for type-safe queries
    var user model.User
    err := r.exec.QueryRowContext(ctx,
        "SELECT id, name FROM users WHERE id = ?", id,
    ).Scan(&user.ID, &user.Name)
    return &user, err
}
```

### 4. Register Route

Edit `backend/server/server.go`:

```go
api := app.Group("/api")
api.Get("/users/:id", h.GetUser)
```

### 5. Generate Swagger Docs

```bash
mise run swag
```

Visit `http://localhost:8080/api/swagger/` to see updated documentation.

## 🎨 Frontend Development

### Creating Pages

Pages in `web/app/pages/` are auto-routed. Example:

```vue
<!-- web/app/pages/users.vue -->
<template>
  <div class="p-8">
    <h1 class="text-3xl font-bold">Users</h1>
    <div class="mt-4" v-if="user">
      <p>{{ user.name }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
const user = ref(null)

onMounted(async () => {
  const response = await fetch('/api/users/1')
  const json = await response.json()
  user.value = json.data
})
</script>
```

### Using Tailwind & DaisyUI

```vue
<template>
  <button class="btn btn-primary">Click me</button>
  <div class="card bg-base-100 shadow-xl">
    <div class="card-body">
      <h2 class="card-title">Card title</h2>
    </div>
  </div>
</template>
```

### Using FontAwesome Icons

Icons are globally available via the `fontawesome.ts` plugin:

```vue
<template>
  <font-awesome-icon icon="fa-heart" />
  <font-awesome-icon :icon="['fab', 'github']" />
</template>
```

## 🗄️ Database Setup

### SQLite (Default)

No setup needed. Database auto-creates at `./db.sqlite` on first run.

**Schema:** Edit `backend/repo/sqlite/queries/init.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
```

Then run:
```bash
mise run jet sqlite
```

### PostgreSQL

Update `backend/config.yml`:

```yaml
DSN: "postgres://user:password@localhost:5432/mydb"
```

**Schema:** Edit `backend/repo/postgres/queries/init.sql`

```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
```

Then run:
```bash
mise run jet postgres
```

## 🐳 Docker Deployment

Build and run in a container:

```bash
docker build -t goweb-template .
docker run -p 8080:8080 goweb-template
```

The Dockerfile builds both frontend and backend, embedding frontend assets in the final Go binary.

## 🧪 Testing

Currently, no test framework is configured. To add tests:

```bash
# Create a test file alongside implementation
touch backend/handler/handler_test.go
```

Example test structure:

```go
package handler

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestDatabaseHealth(t *testing.T) {
    // Setup
    // Execute
    // Assert
}
```

## 📚 Dependencies

### Backend
- `gofiber/fiber/v3` - HTTP framework
- `go-jet/jet/v2` - Type-safe SQL builder
- `lib/pq` - PostgreSQL driver
- `modernc.org/sqlite` - SQLite driver
- `spf13/viper` - Configuration management
- `swaggo/swag` - Swagger documentation

### Frontend
- `nuxt` - Meta-framework
- `tailwindcss` - CSS framework
- `daisyui` - Component library
- `@fortawesome/vue-fontawesome` - Icons

## 🔐 Security Considerations

- **CORS** - Enabled in `backend/server/server.go`
- **SQL Injection** - Protected via go-jet type-safe queries
- **Environment Variables** - Use `.env` for sensitive config (never commit)
- **HTTPS** - Add reverse proxy (nginx, Caddy) for production TLS

## 📖 Additional Resources

- [Fiber Documentation](https://docs.gofiber.io/)
- [Nuxt Documentation](https://nuxt.com/docs)
- [go-jet Query Builder](https://github.com/go-jet/jet)
- [Swagger/OpenAPI](https://swagger.io/)
- [Tailwind CSS](https://tailwindcss.com/)
- [DaisyUI Components](https://daisyui.com/)

## 📄 License

This template is open source and available under the MIT License.

## 💡 Tips

- **Hot reload** - Both frontend (`mise run dev`) and backend (file watcher in Mise) support hot reload
- **API Testing** - Use Swagger UI at `/api/swagger/` or tools like `curl`, Postman, or `rest-client` VS Code extension
- **Database inspection** - For SQLite, use `sqlite3 db.sqlite` CLI or GUI tools like DBeaver
- **Type safety** - Leverage go-jet generated code and TypeScript for maximum safety

---

**Questions or issues?** Check `.github/copilot-instructions.md` for AI agent guidelines or refer to the project structure above.

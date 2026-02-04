# CSV Processor - Progress & Context

## Project Goal
Educational project - learn systems design and backend engineering.
Build async CSV processing system with own infrastructure (no managed services).

## Architecture (Approved in Stages 1-5)
```
[Client] → [API Server] → [RabbitMQ] → [Worker(s)] → [Storage]
                ↓                           ↓
           [PostgreSQL] ←──────────────────┘
```

### Key Decisions
- **Queue:** RabbitMQ (ACK/NACK semantics for reliability)
- **Database:** PostgreSQL (task state storage)
- **API & Worker:** Separate binaries (scale independently)
- **Communication:** Polling (not WebSocket/SSE)
- **Tech Stack:** Go + Gin framework

### Task Lifecycle
```
PENDING → PROCESSING → COMPLETED
                    → FAILED
```

---

## Current Implementation Status

### ✅ Completed

#### 1. Project Structure
```
/csv-processor
  /cmd
    /api
      main.go          ← API server entry point
  /internal
    /api
      server.go        ← Server struct, routes, handlers
    /config
      config.go        ← Config loading from .env
  .env                 ← Environment variables
  go.mod
  go.sum
```

#### 2. Configuration System
- **File:** `internal/config/config.go`
- **Method:** `.env` file with `godotenv`
- **Fields:**
  - `PORT` (API server port)
  - `DB_NAME`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD` (PostgreSQL)

#### 3. API Server
- **File:** `internal/api/server.go`
- **Features:**
  - Server struct with Gin router
  - `NewServer()` constructor
  - `Run()` method
  - `setupRoutes()` for endpoint registration
- **Current endpoints:**
  - `GET /health` → `{"status": "ok"}`

#### 4. Running Infrastructure (Docker)
```
✓ PostgreSQL:  localhost:5555 (user: admin, pass: admin123, db: csv_processor)
✓ Redis:       localhost:6380
✓ MinIO:       localhost:9002/9003
```

#### 5. Working Server
```bash
go run cmd/api/main.go
curl localhost:8080/health  # Works!
```

---

## TODO - Next Steps

### Immediate (MVP)
- [ ] **Database Connection**
  - Add PostgreSQL driver (`pgx` or `database/sql` + `pq`)
  - Create DB connection in `internal/database/postgres.go`
  - Add DB to Server struct
  - Test connection on startup

- [ ] **Task Model & Schema**
  - Define `Task` struct in `internal/models/task.go`
  - Create SQL schema (tasks table)
  - Fields: `id`, `status`, `created_at`, `updated_at`, `error_message`

- [ ] **Task Endpoints (API Contract)**
  - `POST /tasks` - create task, return task_id
  - `GET /tasks/:id` - get task status
  - `GET /tasks/:id/result` - download result (later)
  - Mock responses first, then wire to DB

- [ ] **Database Operations**
  - `CreateTask()` - insert new task as PENDING
  - `GetTask()` - fetch task by ID
  - `UpdateTaskStatus()` - update status (PROCESSING/COMPLETED/FAILED)

- [ ] **File Upload Handling**
  - Accept CSV file in `POST /tasks`
  - Save to MinIO (object storage)
  - Store file_path in DB

### Later (Post-MVP)
- [ ] **RabbitMQ Integration**
  - Setup RabbitMQ connection
  - Publisher in API (send task to queue)
  - Consumer in Worker

- [ ] **Worker Implementation**
  - `/cmd/worker/main.go`
  - Pull tasks from RabbitMQ
  - Process CSV
  - Update task status in DB

- [ ] **Reliability Features**
  - Graceful shutdown (SIGTERM/SIGINT handling)
  - Task timeout mechanism (stuck PROCESSING → FAILED)
  - Retry logic
  - Idempotency

- [ ] **Observability**
  - Logging library (zerolog)
  - Metrics
  - Health checks (DB, Queue)

---

## Key Design Patterns Used

### 1. Dependency Injection
```go
// main.go creates dependencies, passes to server
cfg := config.LoadConfig()
server := api.NewServer(cfg.Port)
```

### 2. Struct Methods
```go
type Server struct { ... }
func (s *Server) Run() error { ... }
```

### 3. Constructor Pattern
```go
func NewServer(port string) *Server {
    srv := &Server{...}
    srv.setupRoutes()
    return srv
}
```

---

## Conventions & Standards

- **Config variables:** Use `cfg` (not `config` or `con`)
- **Exported fields:** CamelCase with capital letter (`DBHost`, not `db_host`)
- **Error handling:** Always check errors, use `log.Fatal()` in main
- **Imports:** Group stdlib, external, internal
- **Port format:** String with colon (`:8080`)

---

## Commands Reference

```bash
# Run API server
go run cmd/api/main.go

# Test health endpoint
curl localhost:8080/health

# Format code
go fmt ./...

# Update dependencies
go mod tidy

# Connect to PostgreSQL (via Docker)
docker exec -it csv-processor-db psql -U admin -d csv_processor
```

---

## Questions to Revisit Later

1. Should routes be split into separate file when we add more endpoints?
2. How to organize handlers - by entity (`task.go`) or by operation?
3. When to add graceful shutdown?
4. Logging strategy - structured logs from day 1 or add later?

---

**Last Updated:** 2026-02-04
**Current Stage:** Stage 6 - Implementation (Database integration next)

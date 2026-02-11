# CSV Processor - Progress & Context

## Project Goal
Educational project - learn systems design and backend engineering.
Build async CSV processing system with own infrastructure (no managed services).

## Architecture
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
PENDING → PROCESSING → FINISHED
                     → FAILED
```
Terminal states: FINISHED, FAILED (retry creates new task via original_task_id)

---

## Project Structure
```
/csv-processor
  /cmd
    /api/main.go                     ← API server entry point
    /worker/main.go                  ← Worker entry point (TBD)
  /internal
    /api
      server.go                      ← Server struct, routes, Run()
      /handlers
        tasks.go                     ← Task CRUD handlers
        health.go                    ← Health endpoint
        repository.go                ← TaskRepository interface
    /config
      config.go                      ← Config from .env
    /db
      db.go                          ← DB connection, setup, pool config
      task.go                        ← TaskRepo (implements TaskRepository)
      errors.go                      ← Sentinel errors
      init_task.go                   ← Schema initialization
    /models
      task.go                        ← Task struct, input DTOs
  .env
  docker-compose.yml
```

## Infrastructure (Docker)
```
PostgreSQL:  localhost:5555
MinIO:       localhost:9002 (API) / 9003 (Console)
RabbitMQ:    commented out in docker-compose (not needed yet)
```

---

## Completed

### Database Integration
- [x] PostgreSQL driver (`database/sql` + `lib/pq`)
- [x] Connection with `Ping()` on startup (fail fast)
- [x] Connection pool config (max open/idle conns, lifetime)
- [x] Schema init with `CREATE TABLE IF NOT EXISTS`
- [x] DB indexes on `status` and `created_at`

### Task Model & Schema
- [x] `Task` struct with JSON/DB tags, nullable fields as pointers
- [x] `TaskCreateInput`, `TaskUpdateStatusInput` with Gin binding validation
- [x] CHECK constraint on status values in DB
- [x] UUID primary keys

### Task Endpoints (API)
- [x] `GET /health`
- [x] `GET /tasks` (dev helper, no pagination)
- [x] `GET /tasks/:task_id`
- [x] `POST /tasks`
- [x] `PUT /tasks/:task_id`
- [x] `DELETE /tasks/:task_id`
- [x] Standardized response format `{"status": "...", "data": "..."}`

### Database Operations (TaskRepo)
- [x] `GetTask()` / `GetTasks()` / `CreateTask()` / `UpdateTask()` / `DeleteTask()`
- [x] All methods use `context.Context` + `*Context` SQL variants
- [x] Full `RETURNING` clauses (no partial/zero-value responses)
- [x] Sentinel errors: `ErrTaskNotFound`, `ErrInvalidTransition`

### Status State Machine
- [x] Atomic UPDATE with transition validation in SQL WHERE clause
- [x] Allowed: `pending→processing`, `processing→finished`, `processing→failed`
- [x] Disambiguation: ErrNoRows → SELECT to distinguish not-found vs invalid transition
- [x] Handler returns 409 Conflict on invalid transitions

### Code Quality & Patterns
- [x] Dependency injection: `main.go` wires cfg → db → repo → handler → server
- [x] Negative space programming: panic in constructors for nil dependencies
- [x] `TaskRepository` interface in handlers package (consumer-side)
- [x] Handler depends on interface, not concrete `*db.TaskRepo`
- [x] Internal errors hidden from client (log server-side, generic 500 response)
- [x] `parseUUID` helper with early return

---

## TODO - Next Steps

### Immediate
- [ ] **File Upload to MinIO**
  - Accept CSV file in `POST /tasks` (multipart form)
  - Connect to MinIO from API
  - Upload file, generate S3 path server-side
  - Remove `s3_input_path` from `TaskCreateInput` (client shouldn't control this)
  - Store generated path in DB

### Post-MVP
- [ ] **RabbitMQ Integration**
  - Setup connection
  - Publisher in API (send task_id to queue after create)
  - Consumer in Worker

- [ ] **Worker Implementation**
  - Fix `package worker` → `package main`
  - Pull tasks from RabbitMQ
  - Download CSV from MinIO, process, upload result
  - Update task status in DB

- [ ] **Reliability**
  - Graceful shutdown (SIGTERM/SIGINT)
  - Task timeout (stuck PROCESSING → FAILED)
  - Retry logic (new task with original_task_id)
  - Idempotency

- [ ] **Observability**
  - Structured logging (zerolog or slog)
  - Health check with DB ping (liveness vs readiness)
  - Metrics

### Known Debt (conscious decisions)
- Schema managed by `CREATE TABLE IF NOT EXISTS`, no migration tool yet
- `GetTasks()` has no pagination, no `rows.Err()` check
- Errors from `db` package imported in handlers (breaks interface isolation)
- No tests yet (deferred until worker/async flow exists)

---

## Conventions

- **Config:** use `cfg`
- **Exported fields:** CamelCase (`DBHost`)
- **Error handling:** check all errors, `log.Fatal()` in main, `panic` for nil deps
- **Imports:** group stdlib / external / internal
- **Port format:** string with colon (`:8080`)
- **Context:** always first param in repo methods, from `c.Request.Context()`
- **Interfaces:** defined where consumed (handler package), not where implemented

---

## Commands

```bash
go run cmd/api/main.go                                          # Run API
curl localhost:8080/health                                      # Health check
go fmt ./...                                                    # Format
go mod tidy                                                     # Deps
docker exec -it csv-processor-db psql -U admin -d csv_processor # DB shell
```

---

**Last Updated:** 2026-02-10
**Current Stage:** File Upload to MinIO

## MVP - Next Steps

### 1. Database Integration
- [x] Add PostgreSQL driver (pgx or database/sql + pq)
- [x] Create `internal/database/postgres.go` with connection logic
- [-] Add DB to Server struct
- [x] Test connection on startup

### 2. Task Model & Schema
- [x] Define `Task` struct in `internal/models/task.go`
- [x] Create SQL schema (tasks table: id, status, created_at, updated_at, error_message)
- [x] Run migrations or manual schema creation

### 3. Task Endpoints
- [x] `POST /tasks` - create task, return task_id (mock first)
- [x] `GET /tasks/:id` - get task status (mock first)
- [x] Wire endpoints to DB operations

### 4. Database Operations
- [x] `CreateTask()` - insert new task as PENDING
- [x] `GetTask()` - fetch task by ID
- [x] `UpdateTaskStatus()` - update task status

### 5. File Upload
- [ ] Accept CSV in `POST /tasks`
- [ ] Save to MinIO
- [ ] Store file_path in DB

---

## Post-MVP

### RabbitMQ
- [ ] Setup RabbitMQ connection
- [ ] Publisher in API server
- [ ] Consumer in Worker

### Worker
- [ ] `/cmd/worker/main.go`
- [ ] Pull tasks from queue
- [ ] Process CSV
- [ ] Update DB

### Reliability
- [ ] Graceful shutdown (SIGTERM/SIGINT)
- [ ] Task timeout mechanism
- [ ] Retry logic
- [ ] Idempotency

### Observability
- [ ] Logging library (zerolog)
- [ ] Metrics
- [ ] Health checks (DB, Queue)

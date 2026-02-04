## MVP - Next Steps

### 1. Database Integration
- [ ] Add PostgreSQL driver (pgx or database/sql + pq)
- [ ] Create `internal/database/postgres.go` with connection logic
- [ ] Add DB to Server struct
- [ ] Test connection on startup

### 2. Task Model & Schema
- [ ] Define `Task` struct in `internal/models/task.go`
- [ ] Create SQL schema (tasks table: id, status, created_at, updated_at, error_message)
- [ ] Run migrations or manual schema creation

### 3. Task Endpoints
- [ ] `POST /tasks` - create task, return task_id (mock first)
- [ ] `GET /tasks/:id` - get task status (mock first)
- [ ] Wire endpoints to DB operations

### 4. Database Operations
- [ ] `CreateTask()` - insert new task as PENDING
- [ ] `GetTask()` - fetch task by ID
- [ ] `UpdateTaskStatus()` - update task status

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

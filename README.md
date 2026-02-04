# CSV Processor

> **Educational Project** - Learning systems design and backend engineering through hands-on implementation

## About This Project

This is an **educational project** focused on understanding distributed systems, async processing, and backend architecture from first principles. The goal is not to build the fastest or most feature-rich solution, but to **deeply understand how systems work** by building them manually.

### Learning Philosophy

- **Hand-written code** - Every line is written manually to understand what's happening under the hood. Doesn't apply to .md files hehe
- **No shortcuts** - Deliberately avoiding managed services (AWS Lambda, SQS, etc.) to see the mechanics
- **Mentorship-driven** - Built with guidance from a senior systems engineer acting as a technical mentor (see `claude.md`)
- **Trade-offs over features** - Focus on understanding why decisions matter, not just making things work

This project prioritizes:
- ✅ Understanding failure modes
- ✅ Learning systems thinking
- ✅ Grasping architectural trade-offs
- ❌ Not production-ready code
- ❌ Not rapid feature development

## What It Does

Async CSV file processing system with the following flow:

1. User uploads a CSV file via API
2. System creates a task and returns a `task_id`
3. Worker processes the CSV asynchronously
4. User polls for status using `task_id`
5. When complete, user downloads the processed result

### Architecture

```
[Client] → [API Server] → [RabbitMQ] → [Worker(s)] → [Storage (MinIO)]
                ↓                          ↓
           [PostgreSQL] ←──────────────────┘
```

**Components:**
- **API Server** - HTTP API for upload, status checking, result download
- **PostgreSQL** - Source of truth for task state
- **RabbitMQ** - Message queue for task distribution
- **Worker(s)** - Processes CSV files from the queue
- **MinIO** - Object storage for input/output files

### Design Decisions

**Queue:** RabbitMQ (not Redis)
- ACK/NACK semantics for reliability
- At-least-once delivery
- Learned trade-off: Complexity vs reliability

**Communication:** Polling (not WebSockets/SSE)
- Simpler implementation
- Good enough for 100 concurrent users
- Learned trade-off: Real-time UX vs simplicity

**Separate API & Worker binaries**
- Independent scaling
- Clear separation of concerns
- More complex deployment

## Tech Stack

- **Language:** Go 1.25
- **Web Framework:** Gin
- **Database:** PostgreSQL 15
- **Queue:** RabbitMQ 3
- **Storage:** MinIO
- **Config:** Environment variables with `.env` file

## Project Structure

```
/csv-processor
  /cmd
    /api          # API server entry point
      main.go
    /worker       # Worker entry point (coming soon)
      main.go
  /internal       # Private application code
    /api
      server.go   # Server struct, routes, handlers
    /config
      config.go   # Configuration loading
    /database     # (coming soon)
    /queue        # (coming soon)
    /models       # (coming soon)
    /processor    # (coming soon)
  docker-compose.yml
  .env.example
  README.md
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Setup

1. **Clone the repository**
```bash
git clone <repo-url>
cd csv-processor
```

2. **Copy environment variables**
```bash
cp .env.example .env
# Edit .env if needed (default values work for local dev)
```

3. **Start infrastructure**
```bash
docker-compose up -d
```

This starts:
- PostgreSQL on `localhost:5555`
- RabbitMQ on `localhost:5672` (management UI: `localhost:15672`)
- MinIO on `localhost:9002` (console: `localhost:9003`)

4. **Install Go dependencies**
```bash
go mod download
```

5. **Run API server**
```bash
go run cmd/api/main.go
```

6. **Test it**
```bash
curl localhost:8080/health
# Should return: {"status":"ok"}
```

### RabbitMQ Management UI

Access at `http://localhost:15672`
- Username: `admin`
- Password: `admin123`

### MinIO Console

Access at `http://localhost:9003`
- Username: `minioadmin`
- Password: `minioadmin123`

## Development Status

🚧 **Work in Progress** - Currently implementing MVP

**Completed:**
- ✅ Project structure
- ✅ Configuration system (.env)
- ✅ API server setup
- ✅ Health endpoint
- ✅ Docker infrastructure

**In Progress:**
- 🔄 Database connection
- 🔄 Task model & schema
- 🔄 Task endpoints (POST /tasks, GET /tasks/:id)

**Coming Soon:**
- ⏳ RabbitMQ integration
- ⏳ Worker implementation
- ⏳ CSV processing logic
- ⏳ File upload/download
- ⏳ Graceful shutdown
- ⏳ Observability (logging, metrics)

See `TODO.md` for detailed task list.

## API Endpoints (Planned)

### Create Task
```
POST /tasks
Content-Type: multipart/form-data

file: <CSV file>

Response:
{
  "task_id": "uuid",
  "status": "pending"
}
```

### Get Task Status
```
GET /tasks/:id

Response:
{
  "task_id": "uuid",
  "status": "pending|processing|completed|failed",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "error": "error message if failed"
}
```

### Download Result
```
GET /tasks/:id/result

Response: CSV file download
```

## Learning Resources

This project explores concepts like:
- Async task processing patterns
- Message queue semantics (at-most-once vs at-least-once)
- Database as source of truth
- Worker pool patterns
- Failure handling (worker crashes, queue unavailable, timeouts)
- Graceful degradation
- Idempotency
- Horizontal scaling

## Why Hand-Written?

**This code is intentionally written manually** (not AI-generated) to ensure deep understanding of:
- How HTTP servers work
- How message queues operate
- How database connections are managed
- How errors propagate through a system
- How to structure Go applications

AI is used as a **mentor** for:
- Asking guiding questions ("what happens if...")
- Code reviews and architectural feedback
- Pointing out failure modes and edge cases
- Explaining trade-offs and alternatives

See `claude.md` for the full mentorship philosophy.

## Non-Goals

This project is **not**:
- ❌ Production-ready (no auth, limited error handling, dev credentials)
- ❌ Optimized for performance (learning > speed)
- ❌ Following all best practices (intentionally exploring trade-offs)
- ❌ Using managed services (building from scratch for education)

## License

MIT - This is educational code, use at your own risk.

## Contributing

This is a personal learning project. If you're building something similar, feel free to fork and adapt, but note this is optimized for learning, not production use.

---

**Built with:** Go, PostgreSQL, RabbitMQ, MinIO, and a lot of "what happens if..." questions.

# Solace

A personal habit tracker and journal app — built for clarity, calm, and consistency.

## What it does

- **Habit tracking** — create habits, check in daily, track streaks
- **Journal** — write entries, log mood, attach images via S3, paginate through history
- **Auth** — JWT-based user accounts, everything is private per user

## Tech Stack

| Layer | Tech |
|---|---|
| Backend | Go (net/http, clean architecture) |
| Database | PostgreSQL |
| Auth | JWT (golang-jwt) |
| Storage | AWS S3 (presigned URLs, direct client upload) |
| Frontend | React + TypeScript + Tailwind |
| Container | Docker + Docker Compose |

## Project Structure

```
solace/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── habit/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── journal/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── middleware/
│   │   └── auth.go
│   └── db/
│       └── postgres.go
├── migrations/
├── docker-compose.yml
├── Dockerfile
└── .env
```

## Database Schema

```
users           — id, name, email, password, created_at, updated_at
habits          — id, user_id(fk), title, image_url, current_streak, last_checked_at, created_at, updated_at
habit_checking  — id, habit_id(fk), checked_date, created_at
moods           — id, label (happy/sad/anxious/calm/neutral)
journal         — id, user_id(fk), mood_id(fk), status, title, description, image_key, created_at, updated_at, deleted_at
```

## API Endpoints

### Auth
```
POST /api/v1/register       — create account (name, email, password)
POST /api/v1/login          — login, returns JWT token
GET  /api/v1/me             — get current user info (protected)
```

### Habits
```
GET  /api/v1/habits              — list user's habits with current streak
POST /api/v1/habits              — create habit
POST /api/v1/habits/:id/check-in — check in for today, returns updated streak
```

### Journal
```
GET    /api/v1/journals           — list published entries (offset/limit pagination)
GET    /api/v1/journals/drafts    — list draft entries
POST   /api/v1/journals           — create entry (draft or published)
GET    /api/v1/journals/:id       — get single entry (image resolved to fresh presigned URL)
PATCH  /api/v1/journals/:id       — update entry (partial update)
DELETE /api/v1/journals/:id       — soft delete entry
```

### Image Upload (S3)
```
POST /api/v1/journals/presign-upload      — get presigned PUT URL for direct S3 upload
POST /api/v1/journals/:id/confirm-upload  — verify upload landed in S3, persist key to DB
```

## Image Upload Flow

```
1. Client → POST /presign-upload        → backend returns S3 key + presigned PUT URL
2. Client → PUT directly to S3 URL      → image lands in S3 (backend not involved)
3. Client → POST /confirm-upload        → backend verifies via HeadObject, saves key to DB
4. Client → GET /journals/:id           → backend generates fresh presigned GET URL per request
```

Key decisions:
- S3 key is stored in DB, never the presigned URL (URLs expire, keys don't)
- Backend verifies object existence via HeadObject before trusting client's confirm
- Ownership verified by parsing userID embedded in the S3 key path
- Private bucket, Block Public Access enabled — all reads/writes via signed URLs only

## Getting Started

```bash
# copy env file
cp .env.example .env

# start everything
docker compose up --build

# server runs at
http://localhost:8000
```

## Environment Variables

```
DSN=postgres://user:password@localhost:5432/solace
JWT_SECRET=your-jwt-secret
TOKEN_VALID_PERIOD=8
PORT=8000
AWS_ACCESS_KEY_ID=...
AWS_SECRET_ACCESS_KEY=...
AWS_REGION=ap-south-1
S3_BUCKET=your-bucket-name
```

## Architecture Decisions

**Three-layer pattern** per feature: handler → service → repository. Handlers parse HTTP, services contain business logic, repositories do DB queries only. Layers never skip — handlers never touch the DB.

**JWT auth** — stateless, no sessions. Token validated in middleware, user ID extracted from claims and passed via context to service layer.

**Moods as lookup table** — fixed set of values seeded once at migration time. Journal entries reference mood by ID, not by string value.

**Habit streaks stored, not recalculated** — current_streak is updated transactionally on each check-in alongside the habit_checking insert. Both succeed or both roll back. Tradeoff: faster reads, but streak can drift if data is manually edited.

**Soft delete on journals** — deleted_at timestamp instead of hard delete. Data is recoverable, and all list queries filter with WHERE deleted_at IS NULL.

**Offset/limit pagination on journals** — simple and sufficient at current scale. Cursor-based pagination would be needed at millions of rows to avoid full table scans.

**Postgres over SQLite** — multiple app instances need a networked database. File-based SQLite causes write locking under concurrent load from multiple containers.

**Direct-to-S3 upload** — images never pass through the backend server. Presigned URLs let the client upload directly to S3, keeping the backend stateless and avoiding memory/bandwidth overhead.

## Learning Context

Built as a hands-on systems design project covering:
- Go backend from scratch with clean architecture
- PostgreSQL schema design and migrations
- JWT authentication and middleware
- AWS S3 presigned URLs and IAM
- Docker multi-stage builds (~27MB final image)
- Docker Compose orchestration
- Horizontal scaling with a shared database
- Concurrency, race conditions, and mutex patterns

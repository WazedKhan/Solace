# Solace

A personal habit tracker and journal app — built for clarity, calm, and consistency.

## What it does

- **Habit tracking** — create habits, check in daily, track streaks
- **Journal** — write entries, log mood, attach images
- **Auth** — JWT-based user accounts, everything is private per user

## Tech Stack

| Layer | Tech |
|---|---|
| Backend | Go |
| Database | PostgreSQL |
| Auth | JWT (golang-jwt) |
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
users         — id, name, email, password, created_at, updated_at
habits        — id, user_id(fk), title, image_url, created_at, updated_at
habit_checking — id, habit_id(fk), checked_date, created_at
moods         — id, label (happy/sad/anxious/calm/neutral)
journal       — id, user_id(fk), mood_id(fk), status, description, image_url, created_at, updated_at
```

## API Endpoints

### Auth
```
POST /auth/register   — create account (name, email, password)
POST /auth/login      — login, returns JWT token
```

### Habits
```
GET    /habits          — list user's habits
POST   /habits          — create habit
DELETE /habits/:id      — archive habit
POST   /habits/:id/checkin — check in for today
GET    /habits/:id/streak  — get current streak
```

### Journal
```
GET    /journal         — list user's entries
POST   /journal         — create entry
GET    /journal/:id     — get single entry
PUT    /journal/:id     — update entry
DELETE /journal/:id     — delete entry
```

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
API_KEY=your-api-key
DATABASE_URL=postgres://user:password@localhost:5432/solace
JWT_SECRET=your-jwt-secret
```

## Architecture Decisions

- **Three-layer pattern** per feature: handler → service → repository. Handlers parse HTTP, services contain business logic, repositories do DB queries only.
- **JWT auth** — stateless, no sessions. Token validated in middleware on every protected route.
- **Moods as lookup table** — fixed set of values seeded once, never changed by users.
- **Habit streaks calculated on query** — no denormalized streak counter that can go out of sync. Query `habit_checking` for consecutive days on each request.
- **Postgres over SQLite** — multiple app instances need a networked database. Learned this the hard way with file locking issues under concurrent load.

## Learning Context

Built as a hands-on systems design project covering:
- Go backend from scratch
- PostgreSQL schema design
- JWT authentication
- Docker multi-stage builds
- Docker Compose orchestration
- Horizontal scaling with a shared database

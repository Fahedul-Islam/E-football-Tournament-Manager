# E-Football Tournament Manager

A production-ready **RESTful API** backend for managing e-football and esports tournaments. Built with **Go** following **Clean Architecture** principles — covering everything from user registration and participant management to automated match scheduling, real-time WebSocket notifications, and a full announcement system with social interactions.

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![WebSocket](https://img.shields.io/badge/WebSocket-Gorilla-00ADD8)
![License](https://img.shields.io/badge/License-MIT-green)
![Architecture](https://img.shields.io/badge/Architecture-Clean-blueviolet)

---

## Table of Contents

- [Project Highlights](#project-highlights)
- [Architecture](#architecture)
- [Features](#features)
- [Database Schema](#database-schema)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Tournament Flow](#tournament-flow)
- [Example API Usage](#example-api-usage)
- [Development](#development)
- [Roadmap](#roadmap)

---

## Project Highlights

- **Clean Architecture** — Strict separation across Domain, Repository, Service, and Handler layers
- **Security-First** — JWT authentication, bcrypt password hashing, role-based access control, and per-IP rate limiting
- **Real-Time Notifications** — WebSocket integration for instant push notifications to connected clients
- **Automated Tournament Logic** — Group generation, round-robin scheduling, and automatic knockout progression
- **Announcement System** — Announcements with comments, threaded replies, reactions, and read-receipt tracking
- **Production-Ready** — Graceful shutdown, health check endpoint, status-code logging, and database migrations
- **Tested** — Unit tests for business logic and security-critical code paths
- **Developer-Friendly** — Makefile for common tasks, golangci-lint config, and GitHub Actions CI

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Layer                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │     Middleware: RateLimit → Auth → Logger → CORS          │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Handlers (REST + WebSocket)                  │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                       Service Layer                             │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ UserSvc │ TournamentSvc │ ParticipantSvc │ AnnouncementSvc│  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                      Repository Layer                           │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  UserRepo │ TournamentRepo │ ParticipantRepo │ AnnRepo    │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                       Domain Layer                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ User │ Tournament │ Participant │ Match │ Notification    │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                      Infrastructure                             │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │     PostgreSQL │ Migrations │ WebSocket Hub │ Config      │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Features

### Authentication & Security

- User registration with **strong password validation** (min 8 chars, upper, lower, digit, special char)
- Passwords stored using **bcrypt** hashing — plaintext is never stored or returned
- **JWT authentication** with access tokens (24h) and refresh tokens (7 days)
- **Role-based access control** — Admin and Player roles with protected routes
- **Per-IP rate limiting** — 10 req/s with a burst of 20 on auth endpoints (login, register)
- WebSocket connections authenticated with JWT

### Tournament Management

- Create, update, and delete tournaments
- Three tournament formats:
  - **Group Stage + Knockout** (World Cup style)
  - **League** (full round-robin, every team vs every team)
- Configurable max players and date scheduling

### Participant Management

- Players request to join tournaments
- Admin approval/rejection workflow
- Team name registration
- Participant status tracking: `pending` / `approved` / `rejected`

### Match System

- **Automated Group Generation** — Random balanced distribution of participants into groups
- **Round-Robin Scheduling** — Every team plays every other team within their group
- **Live Score Updates** — Match scores updated by admin
- **Automatic Stats Calculation** after every score update:
  - Wins, Draws, Losses
  - Goals Scored / Conceded / Difference
  - Points (3 for win, 1 for draw, 0 for loss)

### Leaderboard & Standings

- Group-wise leaderboards sorted by: Points → Goal Difference → Goals Scored
- Live standings update after each match result

### Knockout Stage Automation

- Top 2 from each group automatically qualify
- Automatic bracket generation through: Round of 16 → Quarterfinals → Semifinals → Final
- Winner determined and advanced after each match

### Announcement System

- Admin posts tournament announcements with types: `general`, `update`, `reminder`, `result`, `urgent`
- Announcements can be **pinned** or have **comments disabled**
- **Threaded Comments** — Nested replies with parent-child structure
- **Reactions** — Like/dislike on announcements and comments
- **Read Receipts** — Track which participants have seen each announcement

### Real-Time Notifications

- **WebSocket Hub** — Persistent bidirectional connection for instant updates
- Notifications pushed immediately when announcements are published
- **Notification Center** — All notifications stored persistently with read/unread status
- Paginated notification history
- Mark individual or all notifications as read

### Operations

- **Health Check** — `GET /health` checks DB connectivity; used by deployment tools
- **Graceful Shutdown** — Handles `SIGTERM`/`SIGINT`, waits up to 30s for in-flight requests
- **Request Logging** — Every request logged with method, path, status code, and duration

---

## Database Schema

### Core Tables

```
┌──────────────┐     ┌─────────────────┐     ┌──────────────────┐
│    Users     │     │   Tournaments   │     │   Participants   │
├──────────────┤     ├─────────────────┤     ├──────────────────┤
│ id           │     │ id              │     │ id               │
│ username     │────>│ created_by (FK) │<────│ tournament_id    │
│ email        │     │ name            │     │ user_id (FK)     │
│ password     │     │ tournament_type │     │ team_name        │
│ role         │     │ max_players     │     │ status           │
│ created_at   │     │ start_date      │     │ created_at       │
└──────────────┘     │ end_date        │     └──────────────────┘
                     └─────────────────┘              │
                              │                       │
                     ┌────────┴────────┐              │
                     │                 │              │
              ┌──────▼───────┐  ┌──────▼──────┐  ┌───▼────────────┐
              │    Groups    │  │   Matches   │  │  Player Stats  │
              ├──────────────┤  ├─────────────┤  ├────────────────┤
              │ id           │  │ id          │  │ participant_id │
              │ tournament_id│  │ tournament  │  │ matches_played │
              │ name (A,B,..)│  │ group_id    │  │ wins/draws     │
              └──────────────┘  │ round       │  │ goals_scored   │
                     │          │ score_a/b   │  │ goal_difference│
              ┌──────▼───────┐  │ status      │  │ points         │
              │ Group Teams  │  │ winner_id   │  └────────────────┘
              ├──────────────┤  └─────────────┘
              │ group_id     │
              │ participant  │
              └──────────────┘
```

### Announcement Tables

```
┌─────────────────┐     ┌────────────────────────┐     ┌─────────────────────┐
│  Announcements  │     │  Announcement Comments │     │  Announcement Seen  │
├─────────────────┤     ├────────────────────────┤     ├─────────────────────┤
│ id              │     │ id                     │     │ id                  │
│ tournament_id   │<────│ announcement_id (FK)   │     │ announcement_id(FK) │
│ author_id (FK)  │     │ user_id (FK)           │     │ user_id (FK)        │
│ title           │     │ parent_comment_id (FK) │     │ is_seen             │
│ content         │     │ content                │     │ seen_at             │
│ type            │     │ is_edited              │     └─────────────────────┘
│ is_pinned       │     │ created_at             │
│ is_commentable  │     └────────────────────────┘
│ likes_count     │              │
│ dislikes_count  │              ▼
│ comments_count  │     ┌────────────────────────────┐
│ created_at      │     │  Comment Reactions         │
└─────────────────┘     ├────────────────────────────┤
        │               │ comment_id (FK)            │
        ▼               │ user_id (FK)               │
┌───────────────────┐   │ reaction_type              │
│ Announcement      │   └────────────────────────────┘
│   Reactions       │
├───────────────────┤
│ announcement_id   │
│ user_id (FK)      │
│ reaction_type     │
└───────────────────┘
```

### All 13 Tables

| Table | Description |
|-------|-------------|
| `users` | User accounts with roles (admin/player) |
| `tournaments` | Tournament definitions and settings |
| `participants` | Player registrations for tournaments |
| `groups` | Group stage groupings (A, B, C, ...) |
| `group_teams` | Many-to-many: participants in groups |
| `matches` | Match records with scores and results |
| `player_stats` | Aggregated player statistics |
| `announcements` | Tournament announcements |
| `announcement_reactions` | Likes/dislikes on announcements |
| `announcement_comments` | Comments with threaded replies |
| `announcement_comment_reactions` | Likes/dislikes on comments |
| `announcement_seen` | Read receipts for announcements |
| `notifications` | User notifications with read status |

---

## Tech Stack

| Category | Technology | Purpose |
|----------|-----------|---------|
| **Language** | Go 1.24 | Core backend |
| **HTTP Server** | Go `net/http` (stdlib) | Routing and request handling |
| **Database** | PostgreSQL 16 | Relational data storage |
| **Authentication** | `golang-jwt/jwt/v4` | JWT token generation and validation |
| **Password Hashing** | `golang.org/x/crypto` (bcrypt) | Secure password storage |
| **Rate Limiting** | `golang.org/x/time/rate` | Per-IP request throttling |
| **Real-Time** | `gorilla/websocket` | Bidirectional WebSocket communication |
| **Migrations** | `golang-migrate/migrate` | Versioned schema management |
| **Config** | `godotenv` | `.env` file loading |
| **Containerization** | Docker & Docker Compose | PostgreSQL container |

---

## Project Structure

```
backend/
├── cmd/
│   └── serve.go                 # Server bootstrap, DI, graceful shutdown
├── config/
│   └── config.go                # Environment config (server, DB, JWT)
├── infra/
│   ├── db/
│   │   ├── connections.go       # PostgreSQL connection
│   │   ├── migrate.go           # Auto-migration on startup
│   │   └── migrations/          # 13 versioned SQL migration files
│   └── ws/
│       ├── ws_hub.go            # WebSocket hub — client registry & broadcast
│       └── client.go            # WebSocket client read/write pumps
├── internal/
│   ├── delivery/http/
│   │   ├── handler/
│   │   │   ├── user/            # Register, login, get users
│   │   │   ├── tournament/      # Tournament, match, participant, leaderboard
│   │   │   ├── participant/     # Player join, group view, match schedule
│   │   │   ├── announcement/    # Announcements, comments, reactions, notifications
│   │   │   └── ws/              # WebSocket upgrade handler
│   │   └── middleware/
│   │       ├── auth.go          # JWT validation, role enforcement
│   │       ├── rate_limiter.go  # Per-IP rate limiting (10 req/s, burst 20)
│   │       ├── logger.go        # Request logging with status code + duration
│   │       ├── cors_with_preflight.go
│   │       └── manager.go       # Chainable middleware manager
│   ├── domain/                  # Structs, DTOs, repository interfaces
│   ├── repository/              # SQL implementations of domain interfaces
│   └── service/                 # Business logic layer
├── utils/
│   ├── password.go              # bcrypt hash, compare, validate
│   ├── password_test.go         # Unit tests for password utilities
│   ├── email.go                 # Email format validation
│   └── sendData.go              # JSON response helper
├── .golangci.yml                # Linter configuration
├── docker-compose.yml           # PostgreSQL container
├── Makefile                     # Developer commands
├── go.mod
└── main.go
```

**Tests:**
```
utils/password_test.go                        # HashPassword, ValidatePassword, CheckPasswordHash
internal/service/user/user_service_test.go    # Register: validates before hashing, rejects weak passwords
```

---

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose

### Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/fahedul-islam/e-football-tournament-manager.git
cd e-football-tournament-manager/backend

# 2. Start PostgreSQL
make docker-up

# 3. Configure environment
cp .env.example .env   # then edit with your values

# 4. Run the server
make run
```

The server starts at `http://localhost:8080`. Database migrations run automatically on startup.

### Environment Variables

```env
# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=tournament_manager
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-strong-secret-key

# Environment
ENV=development
```

> **Note:** `JWT_SECRET` should be a long, random string in production. Never commit your `.env` file.

---

## API Endpoints

### System

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/health` | Database connectivity check | - |

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/register` | Register a new user | - |
| `POST` | `/login` | Login and receive JWT tokens | - |
| `GET` | `/users` | Get all users | Admin |

> `/register` and `/login` are rate-limited to 10 requests/second per IP.

### Password Requirements

Passwords must be at least **8 characters** and contain: uppercase letter, lowercase letter, digit, and special character.

Example valid password: `Secure@123`

### Tournament Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/tournaments/create` | Create a tournament |
| `GET` | `/tournaments` | List your tournaments |
| `PUT` | `/tournaments?id={id}` | Update a tournament |
| `DELETE` | `/tournaments?id={id}` | Delete a tournament |
| `GET` | `/tournaments/create_match_schedules?tournament_id={id}&group_count={n}` | Generate groups & matches |
| `GET` | `/tournaments/matches?tournament_id={id}` | List all matches |
| `PATCH` | `/tournaments/matche-score/update` | Update a match score |
| `GET` | `/tournaments/leaderboard?tournament_id={id}` | Group stage standings |

### Participant Management

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/join-tournament?tournament_id={id}&team_name={name}` | Request to join | Player |
| `PATCH` | `/tournaments/approve` | Approve a join request | Admin |
| `PATCH` | `/tournaments/reject` | Reject a join request | Admin |
| `POST` | `/tournaments/addparticipant` | Add participant directly | Admin |
| `POST` | `/tournaments/removeparticipant` | Remove participant | Admin |
| `GET` | `/tournaments/participants?tournament_id={id}` | List all participants | Admin |

### Player Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/tournament/group-distribution?tournament_id={id}` | View group assignments | Player |
| `GET` | `/tournament/match-schedule?tournament_id={id}` | View match schedule | Player |

### Announcement Management (Admin)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/tournaments/announcements?tournament_id={id}` | Create announcement |
| `PUT` | `/tournaments/announcements/update?tournament_id={id}&announcement_id={id}` | Update announcement |
| `DELETE` | `/tournaments/announcements/delete?tournament_id={id}&announcement_id={id}` | Delete announcement |
| `GET` | `/tournaments/announcements/seen_status?tournament_id={id}&announcement_id={id}` | Get read receipts |

### Announcements & Comments (All Authenticated)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/tournaments/announcements?tournament_id={id}` | List all announcements |
| `GET` | `/tournaments/announcements/get?tournament_id={id}&announcement_id={id}` | Get one announcement |
| `POST` | `/tournament/announcement/react?tournament_id={id}&announcement_id={id}&reaction={like\|dislike}` | React to announcement |
| `POST` | `/tournaments/announcements/comments?tournament_id={id}&announcement_id={id}` | Post a comment |
| `GET` | `/tournaments/announcements/comments?tournament_id={id}&announcement_id={id}` | Get comments |
| `GET` | `/tournaments/announcements/comments?...&parent_comment_id={id}` | Get replies to a comment |
| `PUT` | `/tournaments/announcements/comments/edit?tournament_id={id}&comment_id={id}` | Edit your comment |
| `DELETE` | `/tournaments/announcements/comments/delete?tournament_id={id}&comment_id={id}` | Delete your comment |
| `POST` | `/tournaments/announcements/comments/react?tournament_id={id}&comment_id={id}&reaction={like\|dislike}` | React to a comment |

### Notifications (All Authenticated)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/notifications?page={n}` | Paginated notification history |
| `POST` | `/notifications/mark_read?notification_id={id}` | Mark one as read |
| `POST` | `/notifications/mark_all_read` | Mark all as read |

### WebSocket

| Protocol | Endpoint | Description | Auth |
|----------|----------|-------------|------|
| `WebSocket` | `/ws` | Real-time notification stream | JWT |

```javascript
const ws = new WebSocket("ws://localhost:8080/ws", [], {
  headers: { Authorization: "Bearer <token>" },
});
ws.onmessage = (event) => {
  console.log("Notification:", event.data);
};
```

---

## Tournament Flow

```
1. Admin creates tournament
         ↓
2. Players request to join
         ↓
3. Admin approves / rejects participants
         ↓
4. Admin generates match schedules
   (groups auto-created, round-robin matches inserted)
         ↓
5. Admin updates match scores
   (points, goals, standings recalculated automatically)
         ↓
6. Group stage completes
   (top 2 from each group qualify)
         ↓
7. Knockout bracket auto-generated
   (Ro16 → QF → SF → Final)
         ↓
8. Champion crowned
```

---

## Example API Usage

### Health Check

```bash
curl http://localhost:8080/health
# {"status":"ok","db":"connected"}
```

### Register a User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@example.com",
    "password": "Secure@123",
    "role": "player"
  }'
```

> Passwords must meet the requirements above. `Secure@123` is a valid example.

### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "player1@example.com",
    "password": "Secure@123",
    "role": "player"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": "24h0m0s",
  "token_type": "bearer"
}
```

### Create a Tournament (Admin)

```bash
curl -X POST http://localhost:8080/tournaments/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "name": "Champions League 2025",
    "description": "Annual e-football championship",
    "tournament_type": "group+knockout",
    "max_players": 32,
    "start_date": "2025-06-01",
    "end_date": "2025-06-30"
  }'
```

### Join a Tournament (Player)

```bash
curl -X POST "http://localhost:8080/join-tournament?tournament_id=1&team_name=RedDragons" \
  -H "Authorization: Bearer <player_token>"
```

### Generate Match Schedules (Admin)

```bash
curl -X GET "http://localhost:8080/tournaments/create_match_schedules?tournament_id=1&group_count=4" \
  -H "Authorization: Bearer <admin_token>"
```

### Update Match Score (Admin)

```bash
curl -X PATCH http://localhost:8080/tournaments/matche-score/update \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "tournament_id": 1,
    "participant_a_id": 5,
    "participant_b_id": 8,
    "round": "Group Stage",
    "score_a": 3,
    "score_b": 1
  }'
```

### Get Leaderboard

```bash
curl -X GET "http://localhost:8080/tournaments/leaderboard?tournament_id=1" \
  -H "Authorization: Bearer <admin_token>"
```

### Create an Announcement (Admin)

```bash
curl -X POST "http://localhost:8080/tournaments/announcements?tournament_id=1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "title": "Semifinal Schedule",
    "content": "Semifinals scheduled for June 25th at 8PM.",
    "announcement_type": "update",
    "is_pinned": true,
    "is_commentable": true
  }'
```

### Get Notifications

```bash
curl -X GET "http://localhost:8080/notifications?page=1" \
  -H "Authorization: Bearer <token>"
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 2,
    "notification_type": "announcement",
    "reference_id": 5,
    "message": "New announcement: Semifinal Schedule",
    "is_read": false,
    "created_at": "2025-06-20T17:16:22.431558Z"
  }
]
```

---

## Development

### Available Make Commands

```bash
make run          # Start the server
make build        # Compile to bin/server
make test         # Run all tests with coverage
make lint         # Run golangci-lint
make docker-up    # Start PostgreSQL container
make docker-down  # Stop PostgreSQL container
make migrate-up   # Apply pending migrations
make migrate-down # Roll back last migration
make clean        # Remove compiled binaries
```

### Running Tests

```bash
make test

# Or run specific packages:
go test ./utils/... -v
go test ./internal/service/user/... -v
```

### Running the Linter

```bash
# Install golangci-lint first:
# https://golangci-lint.run/usage/install/

make lint
```

### Running Migrations Manually

```bash
make migrate-up
make migrate-down

# Or directly:
migrate -path infra/db/migrations \
  -database "postgres://postgres:secret@localhost:5434/tournament_manager?sslmode=disable" up
```

### Project Conventions

- **Handlers** — HTTP request parsing and response writing only. No business logic.
- **Services** — All business logic and validation lives here.
- **Repositories** — All SQL queries live here. No business logic.
- **Domain** — Plain struct definitions and repository interfaces. No dependencies.

---

## Roadmap

- [x] Announcement system with comments, reactions, and threaded replies
- [x] WebSocket real-time notifications
- [x] Notification center with read/unread status
- [x] Rate limiting middleware
- [x] Health check endpoint
- [x] Graceful shutdown
- [x] Unit tests for business logic
- [x] Makefile and golangci-lint configuration
- [x] GitHub Actions CI pipeline
- [ ] Refresh token endpoint (`POST /auth/refresh`)
- [ ] Swagger / OpenAPI documentation (`/swagger/`)
- [ ] Email notifications for match schedules
- [ ] Tournament bracket visualization API
- [ ] Player rankings across multiple tournaments
- [ ] Double-elimination bracket support
- [ ] Match rescheduling functionality

---

## Contributing

Contributions are welcome. Please open an issue first to discuss what you'd like to change.

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature`
3. Commit your changes: `git commit -m 'Add your feature'`
4. Push and open a Pull Request

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Author

**Fahedul Islam**

GitHub: [@Fahedul-Islam](https://github.com/Fahedul-Islam)

---

<p align="center">Built with Go</p>

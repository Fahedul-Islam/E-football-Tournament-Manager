# ⚽ E-Football Tournament Manager

A robust, production-ready **RESTful API** backend for managing e-football (or any sports) tournaments, built with **Go** following **Clean Architecture** principles. This system handles everything from user registration to knockout stage generation.

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Architecture](https://img.shields.io/badge/Architecture-Clean-blueviolet)

---

## 🎯 Project Highlights

- **Clean Architecture** - Separation of concerns with Domain, Repository, Service, and Handler layers
- **Role-Based Access Control** - JWT authentication with admin/player roles
- **Automated Tournament Logic** - Group generation, match scheduling, and knockout progression
- **Database Migrations** - Version-controlled schema with golang-migrate
- **Production-Ready** - Middleware for logging, CORS, and authentication

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │           Middleware (Auth, CORS, Logger)            │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Handlers                           │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │   UserService │ TournamentService │ ParticipantSvc  │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                   Repository Layer                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │     UserRepo  │  TournamentRepo  │  ParticipantRepo │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                    Domain Layer                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  User │ Tournament │ Participant │ Match │ Group    │   │
│  └─────────────────────────────────────────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                  Infrastructure                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │         PostgreSQL │ Migrations │ Config            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## ✨ Features

### 🔐 Authentication & Authorization

- User registration with secure password hashing (bcrypt)
- JWT-based authentication with access & refresh tokens
- Role-based access control (Admin/Player)
- Protected routes with middleware

### 🏆 Tournament Management

- Create, update, delete tournaments
- Support for multiple tournament types:
  - **Group Stage + Knockout** (World Cup style)
  - **League Format** (Round-robin)
- Configurable max players per tournament
- Date-based scheduling

### 👥 Participant Management

- Players can request to join tournaments
- Admin approval/rejection workflow
- Team name registration
- Participant status tracking (pending/approved/rejected)

### 📊 Match System

- **Automated Group Generation** - Random distribution of participants
- **Round-Robin Scheduling** - Every team plays each other in group stage
- **Live Score Updates** - Real-time match score tracking
- **Automatic Stats Calculation**:
  - Wins, Draws, Losses
  - Goals Scored/Conceded
  - Goal Difference
  - Points (3 for win, 1 for draw)

### 🏅 Leaderboard & Standings

- Group-wise leaderboards sorted by:
  1. Points
  2. Goal Difference
  3. Goals Scored
- Real-time standings updates after each match

### ⚡ Knockout Stage Automation

- Auto-qualification of top 2 from each group
- Supports progression through:
  - Round of 16
  - Quarterfinals
  - Semifinals
  - Final
- Winner determination and advancement

---

## 🗄️ Database Schema

```sql
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

---

## 🛠️ Tech Stack

| Category              | Technology                     |
| --------------------- | ------------------------------ |
| **Language**          | Go 1.24                        |
| **Database**          | PostgreSQL 16                  |
| **Authentication**    | JWT (golang-jwt/jwt/v4)        |
| **Password Hashing**  | bcrypt (golang.org/x/crypto)   |
| **Migrations**        | golang-migrate/migrate         |
| **Config Management** | godotenv                       |
| **Containerization**  | Docker & Docker Compose        |
| **HTTP Server**       | Go Standard Library (net/http) |

---

## 📁 Project Structure

```
backend/
├── cmd/
│   └── serve.go                 # Application entry point & DI
├── config/
│   └── config.go                # Environment configuration
├── infra/
│   └── db/
│       ├── connections.go       # Database connection
│       ├── migrate.go           # Migration runner
│       └── migrations/          # SQL migration files (11 migrations)
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/         # HTTP handlers (controllers)
│   │       │   ├── user/
│   │       │   ├── participant/
│   │       │   └── tournamentManager/
│   │       └── middleware/      # Auth, CORS, Logger middlewares
│   ├── domain/                  # Business entities & DTOs
│   │   ├── user.go
│   │   ├── tournament.go
│   │   ├── participant.go
│   │   ├── match.go
│   │   ├── group.go
│   │   └── player_stat.go
│   ├── repository/              # Data access layer
│   │   ├── user-repo/
│   │   ├── participant_repo/
│   │   └── tournament_manager_repo/
│   └── service/                 # Business logic layer
│       ├── user/
│       ├── participant/
│       └── tournament/
├── utils/                       # Helper functions
├── docker-compose.yml
├── go.mod
└── main.go
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 16 (or use Docker)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/e-football-tournament-manager.git
   cd e-football-tournament-manager/backend
   ```

2. **Set up environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

   ```env
   DB_HOST=localhost
   DB_PORT=5434
   DB_USER=postgres
   DB_PASSWORD=secret
   DB_NAME=tournament_manager
   DB_SSLMODE=disable
   JWT_SECRET=your-super-secret-key
   SERVER_PORT=8080
   ```

3. **Start PostgreSQL with Docker**

   ```bash
   docker-compose up -d
   ```

4. **Run the application**

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`

---

## 📡 API Endpoints

### Authentication

| Method | Endpoint    | Description         | Auth  |
| ------ | ----------- | ------------------- | ----- |
| `POST` | `/register` | Register a new user | -     |
| `POST` | `/login`    | User login          | -     |
| `GET`  | `/users`    | Get all users       | Admin |

### Tournament Management (Admin Only)

| Method   | Endpoint                                      | Description               |
| -------- | --------------------------------------------- | ------------------------- |
| `POST`   | `/tournaments/create`                         | Create a new tournament   |
| `GET`    | `/tournaments`                                | Get all tournaments       |
| `PUT`    | `/tournaments?id={id}`                        | Update tournament         |
| `DELETE` | `/tournaments?id={id}`                        | Delete tournament         |
| `GET`    | `/tournaments/create_match_schedules?id={id}` | Generate groups & matches |
| `GET`    | `/tournaments/matches?id={id}`                | Get all matches           |
| `PATCH`  | `/tournaments/matche-score/update`            | Update match score        |
| `GET`    | `/tournaments/leaderboard?id={id}`            | Get group standings       |

### Participant Management

| Method  | Endpoint                              | Description              | Auth   |
| ------- | ------------------------------------- | ------------------------ | ------ |
| `POST`  | `/join-tournament`                    | Request to join          | Player |
| `PATCH` | `/tournaments/approve?p_id={id}`      | Approve participant      | Admin  |
| `PATCH` | `/tournaments/reject?p_id={id}`       | Reject participant       | Admin  |
| `POST`  | `/tournaments/addparticipant`         | Add participant directly | Admin  |
| `POST`  | `/tournaments/removeparticipant`      | Remove participant       | Admin  |
| `GET`   | `/tournaments/participants?t_id={id}` | Get all participants     | Admin  |

### Player Endpoints

| Method | Endpoint                                 | Description            | Auth   |
| ------ | ---------------------------------------- | ---------------------- | ------ |
| `GET`  | `/tournament/group-distribution?id={id}` | View group assignments | Player |
| `GET`  | `/tournament/match-schedule?id={id}`     | View match schedule    | Player |

---

## 🔄 Tournament Flow

```
1. Admin creates tournament
         ↓
2. Players request to join
         ↓
3. Admin approves/rejects participants
         ↓
4. Admin generates match schedules
   (Auto-creates groups + round-robin matches)
         ↓
5. Admin updates match scores
   (Stats auto-calculated)
         ↓
6. Group stage completes
   (Top 2 from each group qualify)
         ↓
7. Knockout stage auto-generated
   (Ro16 → QF → SF → Final)
         ↓
8. Champion crowned! 🏆
```

---

## 🧪 Example API Usage

### Register a User

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "player1",
    "email": "player1@email.com",
    "password": "securePass123"
  }'
```

### Create a Tournament (Admin)

```bash
curl -X POST http://localhost:8080/tournaments/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "name": "Champions League 2024",
    "description": "Annual e-football championship",
    "tournament_type": "group_knockout",
    "max_players": 32,
    "start_date": "2024-03-01",
    "end_date": "2024-03-31"
  }'
```

### Update Match Score

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

---

## 🔧 Development

### Running Migrations Manually

```bash
# Up
migrate -path infra/db/migrations -database "postgres://user:pass@localhost:5434/tournament_manager?sslmode=disable" up

# Down
migrate -path infra/db/migrations -database "postgres://user:pass@localhost:5434/tournament_manager?sslmode=disable" down
```

### Project Conventions

- **Handlers**: HTTP request/response handling only
- **Services**: Business logic and validation
- **Repositories**: Database operations
- **Domain**: Entity definitions and DTOs

---

## 🗺️ Roadmap

- [ ] WebSocket support for live match updates
- [ ] Tournament bracket visualization API
- [ ] Email notifications for match schedules
- [ ] Player rankings across tournaments
- [ ] Support for double-elimination brackets
- [ ] Match rescheduling functionality
- [ ] Tournament statistics & analytics API
- [ ] Rate limiting middleware
- [ ] Swagger/OpenAPI documentation

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Fhedul**

- GitHub: [@fhedul](https://github.com/fhedul)

---

<p align="center">
  <b>Built with ❤️ and Go</b>
</p>

# ⚽ E-Football Tournament Manager

A comprehensive, production-ready **RESTful API** backend for managing e-football (and esports) tournaments. Built with **Go** following **Clean Architecture** principles, this platform delivers enterprise-grade tournament management capabilities—from user registration and participant coordination to automated match scheduling, real-time notifications via **WebSocket**, and a full-featured announcement system with social interactions.

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![WebSocket](https://img.shields.io/badge/WebSocket-Gorilla-00ADD8?logo=websocket&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Architecture](https://img.shields.io/badge/Architecture-Clean-blueviolet)

---

## 🎯 Project Highlights

- **Clean Architecture** — Strict separation of concerns across Domain, Repository, Service, and Handler layers
- **Role-Based Access Control** — JWT authentication with granular admin/player permissions
- **Real-Time Notifications** — WebSocket integration for instant push notifications to connected clients
- **Automated Tournament Logic** — Intelligent group generation, round-robin scheduling, and knockout progression
- **Announcement System** — Full-featured announcements with comments, reactions, threaded replies, and read tracking
- **Notification Center** — Persistent notification storage with read/unread status management
- **Database Migrations** — Version-controlled schema evolution with golang-migrate (13 migrations)
- **Production-Ready** — Enterprise middleware stack for logging, CORS, and authentication

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                      HTTP Layer                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │           Middleware (Auth, CORS, Logger)                 │  │
│  └───────────────────────────────────────────────────────────┘  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Handlers (REST + WebSocket)                  │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                    Service Layer                                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ UserSvc │ TournamentSvc │ ParticipantSvc │ AnnouncementSvc│  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                   Repository Layer                              │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  UserRepo │ TournamentRepo │ ParticipantRepo │ AnnRepo    │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                    Domain Layer                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ User │ Tournament │ Participant │ Match │ Notification    │  │
│  └───────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                  Infrastructure                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │     PostgreSQL │ Migrations │ WebSocket Hub │ Config      │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## ✨ Features

### 🔐 Authentication & Authorization

- User registration with secure password hashing (bcrypt)
- JWT-based authentication with access & refresh tokens
- Role-based access control (Admin/Player)
- Protected routes with middleware
- WebSocket authentication support

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

- **Automated Group Generation** — Random distribution of participants into balanced groups
- **Round-Robin Scheduling** — Every team plays each other within group stage
- **Live Score Updates** — Real-time match score tracking
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

### 📢 Announcement System

- **Create Announcements** — Admins can post tournament announcements
- **Announcement Types** — general, update, reminder, result, urgent, other
- **Pinned Announcements** — Important announcements can be pinned to top
- **Commentable Toggle** — Enable/disable comments per announcement
- **Reactions** — Like/dislike on announcements and comments
- **Threaded Comments** — Nested reply support with parent comments
- **Edit/Delete Comments** — Users can manage their own comments
- **Seen Status Tracking** — Track which participants have seen announcements

### 🔔 Real-Time Notifications

- **WebSocket Integration** — Persistent bidirectional connection for instant updates
- **Push Notifications** — Instant delivery when announcements are published
- **Notification Center** — Persistent storage of all notifications per user
- **Read Status Management** — Mark individual or all notifications as read
- **Paginated History** — Fetch notification history with pagination support
- **User-Targeted Delivery** — Notifications routed only to relevant tournament participants

---

## 🗄️ Database Schema

### Core Tables

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

### Announcement Tables

```sql
┌─────────────────┐     ┌────────────────────────┐     ┌─────────────────────┐
│  Announcements  │     │  Announcement Comments │     │  Announcement Seen  │
├─────────────────┤     ├────────────────────────┤     ├─────────────────────┤
│ id              │     │ id                     │     │ id                  │
│ tournament_id   │<────│ announcement_id (FK)   │     │ announcement_id(FK) │
│ author_id (FK)  │     │ user_id (FK)           │     │ user_id (FK)        │
│ title           │     │ parent_comment_id (FK) │     │ is_seen             │
│ content         │     │ content                │     │ seen_at             │
│ announcement_   │     │ is_edited              │     └─────────────────────┘
│   type          │     │ likes_count            │
│ is_pinned       │     │ dislikes_count         │
│ is_commentable  │     │ created_at             │
│ likes_count     │     │ updated_at             │
│ dislikes_count  │     └────────────────────────┘
│ comments_count  │              │
│ created_at      │              │
│ updated_at      │              ▼
└─────────────────┘     ┌────────────────────────────┐
        │               │ Announcement Comment       │
        │               │      Reactions             │
        ▼               ├────────────────────────────┤
┌───────────────────┐   │ id                         │
│ Announcement      │   │ comment_id (FK)            │
│   Reactions       │   │ user_id (FK)               │
├───────────────────┤   │ reaction_type (like/       │
│ id                │   │   dislike)                 │
│ announcement_id   │   │ created_at                 │
│ user_id (FK)      │   └────────────────────────────┘
│ reaction_type     │
│ created_at        │
└───────────────────┘
```

### Notification Tables

```sql
┌─────────────────────────────────────────────────────────────────┐
│                     Notifications                               │
├─────────────────────────────────────────────────────────────────┤
│ id                    SERIAL PRIMARY KEY                        │
│ user_id               INTEGER NOT NULL (FK → users)             │
│ notification_type     VARCHAR(50) NOT NULL                      │
│ reference_id          INTEGER NOT NULL                          │
│ message               TEXT NOT NULL                             │
│ is_read               BOOLEAN DEFAULT FALSE                     │
│ created_at            TIMESTAMPTZ DEFAULT NOW()                 │
├─────────────────────────────────────────────────────────────────┤
│ Indexes:                                                        │
│  • idx_notifications_user_id (user_id)                          │
│  • idx_notifications_user_created (user_id, created_at DESC)    │
│  • idx_notifications_unread (user_id) WHERE is_read = FALSE     │
└─────────────────────────────────────────────────────────────────┘
```

### All Database Tables (13 Migrations)

| Table                            | Description                             |
| -------------------------------- | --------------------------------------- |
| `users`                          | User accounts with roles (admin/player) |
| `tournaments`                    | Tournament definitions and settings     |
| `participants`                   | Player registrations for tournaments    |
| `groups`                         | Group stage groupings (A, B, C, etc.)   |
| `group_teams`                    | Many-to-many: participants in groups    |
| `matches`                        | Match records with scores and results   |
| `player_stats`                   | Aggregated player statistics            |
| `announcements`                  | Tournament announcements                |
| `announcement_reactions`         | Likes/dislikes on announcements         |
| `announcement_comments`          | Comments with threaded replies          |
| `announcement_comment_reactions` | Likes/dislikes on comments              |
| `announcement_seen`              | Read receipts for announcements         |
| `notifications`                  | User notifications with read status     |

---

## 🛠️ Tech Stack

| Category              | Technology                     |
| --------------------- | ------------------------------ |
| **Language**          | Go 1.24                        |
| **Database**          | PostgreSQL 16                  |
| **Real-Time**         | WebSocket (gorilla/websocket)  |
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
│   ├── db/
│   │   ├── connections.go       # Database connection
│   │   ├── migrate.go           # Migration runner
│   │   └── migrations/          # SQL migration files (13 migrations)
│   └── ws/
│       ├── ws_hub.go            # WebSocket hub (client registry & broadcast)
│       └── client.go            # WebSocket client connection handler
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/         # HTTP handlers (controllers)
│   │       │   ├── user/
│   │       │   ├── participant/
│   │       │   ├── tournament/
│   │       │   ├── announcement/
│   │       │   └── ws/          # WebSocket handler
│   │       └── middleware/      # Auth, CORS, Logger middlewares
│   ├── domain/                  # Business entities & DTOs
│   │   ├── user.go
│   │   ├── tournament.go
│   │   ├── participant.go
│   │   ├── match.go
│   │   ├── group.go
│   │   ├── player_stat.go
│   │   ├── announcement.go
│   │   └── notification.go
│   ├── repository/              # Data access layer
│   │   ├── user/
│   │   ├── participant_repo/
│   │   ├── tournament_manager_repo/
│   │   └── announcement/
│   └── service/                 # Business logic layer
│       ├── user/
│       ├── participant/
│       ├── tournament/
│       └── announcement/
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

| Method   | Endpoint                                                                     | Description               |
| -------- | ---------------------------------------------------------------------------- | ------------------------- |
| `POST`   | `/tournaments/create`                                                        | Create a new tournament   |
| `GET`    | `/tournaments`                                                               | Get all tournaments       |
| `PUT`    | `/tournaments?id={id}`                                                       | Update tournament         |
| `DELETE` | `/tournaments?id={id}`                                                       | Delete tournament         |
| `GET`    | `/tournaments/create_match_schedules?tournament_id={id}&group_count={count}` | Generate groups & matches |
| `GET`    | `/tournaments/matches?tournament_id={id}`                                    | Get all matches           |
| `PATCH`  | `/tournaments/matche-score/update`                                           | Update match score        |
| `GET`    | `/tournaments/leaderboard?tournament_id={id}`                                | Get group standings       |

### Participant Management

| Method  | Endpoint                                               | Description              | Auth   |
| ------- | ------------------------------------------------------ | ------------------------ | ------ |
| `POST`  | `/join-tournament?tournament_id={id}&team_name={name}` | Request to join          | Player |
| `PATCH` | `/tournaments/approve`                                 | Approve participant      | Admin  |
| `PATCH` | `/tournaments/reject`                                  | Reject participant       | Admin  |
| `POST`  | `/tournaments/addparticipant`                          | Add participant directly | Admin  |
| `POST`  | `/tournaments/removeparticipant`                       | Remove participant       | Admin  |
| `GET`   | `/tournaments/participants?tournament_id={id}`         | Get all participants     | Admin  |

### Player Endpoints

| Method | Endpoint                                            | Description            | Auth   |
| ------ | --------------------------------------------------- | ---------------------- | ------ |
| `GET`  | `/tournament/group-distribution?tournament_id={id}` | View group assignments | Player |
| `GET`  | `/tournament/match-schedule?tournament_id={id}`     | View match schedule    | Player |

### Announcement Management (Admin)

| Method   | Endpoint                                                                         | Description         | Auth  |
| -------- | -------------------------------------------------------------------------------- | ------------------- | ----- |
| `POST`   | `/tournaments/announcements?tournament_id={id}`                                  | Create announcement | Admin |
| `PUT`    | `/tournaments/announcements/update?tournament_id={id}&announcement_id={id}`      | Update announcement | Admin |
| `DELETE` | `/tournaments/announcements/delete?tournament_id={id}&announcement_id={id}`      | Delete announcement | Admin |
| `GET`    | `/tournaments/announcements/seen_status?tournament_id={id}&announcement_id={id}` | Get seen status     | Admin |

### Announcement Viewing (All Authenticated Users)

| Method | Endpoint                                                                 | Description             | Auth |
| ------ | ------------------------------------------------------------------------ | ----------------------- | ---- |
| `GET`  | `/tournaments/announcements?tournament_id={id}`                          | Get all announcements   | Any  |
| `GET`  | `/tournaments/announcements/get?tournament_id={id}&announcement_id={id}` | Get single announcement | Any  |

### Announcement Reactions (Player)

| Method | Endpoint                                                                                          | Description           | Auth   |
| ------ | ------------------------------------------------------------------------------------------------- | --------------------- | ------ |
| `POST` | `/tournament/announcement/react?tournament_id={id}&announcement_id={id}&reaction={like\|dislike}` | React to announcement | Player |

### Announcement Comments (All Authenticated Users)

| Method   | Endpoint                                                                                                | Description      | Auth |
| -------- | ------------------------------------------------------------------------------------------------------- | ---------------- | ---- |
| `POST`   | `/tournaments/announcements/comments?tournament_id={id}&announcement_id={id}`                           | Create comment   | Any  |
| `GET`    | `/tournaments/announcements/comments?tournament_id={id}&announcement_id={id}`                           | Get comments     | Any  |
| `GET`    | `/tournaments/announcements/comments?tournament_id={id}&announcement_id={id}&parent_comment_id={id}`    | Get replies      | Any  |
| `PUT`    | `/tournaments/announcements/comments/edit?tournament_id={id}&comment_id={id}`                           | Edit comment     | Any  |
| `DELETE` | `/tournaments/announcements/comments/delete?tournament_id={id}&comment_id={id}`                         | Delete comment   | Any  |
| `POST`   | `/tournaments/announcements/comments/react?tournament_id={id}&comment_id={id}&reaction={like\|dislike}` | React to comment | Any  |

### Notifications (All Authenticated Users)

| Method | Endpoint                                        | Description                      | Auth |
| ------ | ----------------------------------------------- | -------------------------------- | ---- |
| `GET`  | `/notifications?page={page}`                    | Get paginated notifications      | Any  |
| `POST` | `/notifications/mark_read?notification_id={id}` | Mark single notification as read | Any  |
| `POST` | `/notifications/mark_all_read`                  | Mark all notifications as read   | Any  |

### WebSocket (Real-Time Notifications)

| Protocol    | Endpoint | Description                       | Auth |
| ----------- | -------- | --------------------------------- | ---- |
| `WebSocket` | `/ws`    | Real-time notification connection | JWT  |

**WebSocket Connection Example:**

```javascript
const ws = new WebSocket("ws://localhost:8080/ws", [], {
  headers: { Authorization: "Bearer <token>" },
});

ws.onmessage = (event) => {
  console.log("Notification:", event.data);
  // Example: "New announcement: Tournament Schedule Update"
};
```

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
    "email": "player1@gmail.com",
    "password": "securePass123",
    "role": "player"
  }'
```

### Create a Tournament (Admin)

```bash
curl -X POST http://localhost:8080/tournaments/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "name": "Champions League 2024",
    "description": "Annual e-football championship with group stage and knockout rounds",
    "tournament_type": "group+knockout",
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

### Join Tournament (Player)

```bash
curl -X POST "http://localhost:8080/join-tournament?tournament_id=1&team_name=MyTeam" \
  -H "Authorization: Bearer <player_token>"
```

### Generate Match Schedules (Admin)

```bash
curl -X GET "http://localhost:8080/tournaments/create_match_schedules?tournament_id=1&group_count=4" \
  -H "Authorization: Bearer <admin_token>"
```

### Get Group Stage Leaderboard (Admin)

```bash
curl -X GET "http://localhost:8080/tournaments/leaderboard?tournament_id=1" \
  -H "Authorization: Bearer <admin_token>"
```

### View Match Schedule (Player)

```bash
curl -X GET "http://localhost:8080/tournament/match-schedule?tournament_id=1" \
  -H "Authorization: Bearer <player_token>"
```

### Create Announcement (Admin)

```bash
curl -X POST "http://localhost:8080/tournaments/announcements?tournament_id=1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin_token>" \
  -d '{
    "title": "Tournament Schedule Update",
    "content": "The semifinal matches have been rescheduled to March 15th.",
    "announcement_type": "update",
    "is_pinned": true,
    "is_commentable": true
  }'
```

### Get All Announcements

```bash
curl -X GET "http://localhost:8080/tournaments/announcements?tournament_id=1" \
  -H "Authorization: Bearer <token>"
```

### React to Announcement (Player)

```bash
curl -X POST "http://localhost:8080/tournament/announcement/react?tournament_id=1&announcement_id=1&reaction=like" \
  -H "Authorization: Bearer <player_token>"
```

### Comment on Announcement

```bash
curl -X POST "http://localhost:8080/tournaments/announcements/comments?tournament_id=1&announcement_id=1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "content": "Great update! Looking forward to the matches."
  }'
```

### Reply to a Comment

```bash
curl -X POST "http://localhost:8080/tournaments/announcements/comments?tournament_id=1&announcement_id=1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "content": "I agree with this comment!",
    "parent_comment_id": 5
  }'
```

### Get Comments on Announcement

```bash
curl -X GET "http://localhost:8080/tournaments/announcements/comments?tournament_id=1&announcement_id=1" \
  -H "Authorization: Bearer <token>"
```

### React to Comment

```bash
curl -X POST "http://localhost:8080/tournaments/announcements/comments/react?tournament_id=1&comment_id=5&reaction=like" \
  -H "Authorization: Bearer <token>"
```

### Edit Comment

```bash
curl -X PUT "http://localhost:8080/tournaments/announcements/comments/edit?tournament_id=1&comment_id=5" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "content": "Updated comment text here"
  }'
```

### Delete Comment

```bash
curl -X DELETE "http://localhost:8080/tournaments/announcements/comments/delete?tournament_id=1&comment_id=5" \
  -H "Authorization: Bearer <token>"
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
    "message": "New announcement: Tournament Schedule Update",
    "is_read": false,
    "created_at": "2026-03-08T17:16:22.431558Z"
  }
]
```

### Mark Notification as Read

```bash
curl -X POST "http://localhost:8080/notifications/mark_read?notification_id=1" \
  -H "Authorization: Bearer <token>"
```

### Mark All Notifications as Read

```bash
curl -X POST "http://localhost:8080/notifications/mark_all_read" \
  -H "Authorization: Bearer <token>"
```

### Connect to WebSocket for Real-Time Notifications

```python
import websocket

def on_message(ws, message):
    print(f"Notification: {message}")

ws = websocket.WebSocketApp(
    "ws://localhost:8080/ws",
    header={"Authorization": "Bearer <token>"},
    on_message=on_message
)
ws.run_forever()
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

- [x] Announcement system with comments and reactions
- [x] WebSocket support for real-time notifications
- [x] Notification center with read/unread status
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

**Fahedul Islam**

- GitHub: [@fahedul-islam](https://github.com/Fahedul-Islam)

---

<p align="center">
  <b>Built with ❤️ and Go</b>
</p>

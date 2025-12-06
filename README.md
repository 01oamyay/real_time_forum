# Real-Time Forum

A full-stack forum built with Go, SQLite, and a vanilla JavaScript single-page application. Users can register, create posts, react through likes/dislikes, leave comments, and chat in real time over WebSockets. The backend bootstraps automatically, applies SQL migrations, serves the SPA, and exposes a REST API that the front-end consumes.

## Feature Highlights

- **Authentication & sessions** – custom JWT implementation backed by a `sessions` table; cookies are HTTP-only while the SPA stores light-weight user metadata in `localStorage`.
- **Post lifecycle** – create posts with up to five categories, infinite-scroll feed with category filters, vote tracking per user, and comment threads with their own reactions.
- **Real-time chat** – `/ws` WebSocket endpoint keeps a synchronized contact list, routes direct messages, exposes typing indicators, and dispatches toast notifications when a new message arrives.
- **Rate limiting & validation** – per-IP limiter guards HTTP handlers, server-side validation (email, password strength, etc.) prevents invalid data, and HTML escaping protects the chat surface.
- **Auto migrations & seed data** – every start reads `schemas/*.sql`, creates tables, and loads seed categories/posts so development instances get data immediately.
- **Front-end router** – lightweight SPA (`web/src/index.js`) assembles shared layout pieces (navbar, sidebar, contacts), lazy-loads CSS per view, and coordinates fetch calls, toasts, and WebSocket events.

## Architecture Overview

### Backend layout

- `internal/app` – bootstraps logging, database connection (via `pkg/database`), repositories, services, controllers, and the HTTP server.
- `internal/controller` – HTTP handlers, middleware, router, WebSocket hub, and utilities like the IP-based rate limiter.
- `internal/service` – business logic per domain (users, posts, comments, categories, chat) layered on top of repositories.
- `internal/repository` – raw SQL access to SQLite. The `Repository` struct exposes composable interfaces for dependency injection.
- `pkg/smplJwt` – small self-contained HS256 JWT implementation used when issuing auth cookies.
- `pkg/utils` – registration validators and bcrypt helpers.
- `schemas` – ordered SQL migration files that create tables (`users`, `post`, `category`, votes, chats, etc.) and insert demo data.

### Frontend layout

- `web/index.html` – served by Go, injects `API_HOST_NAME` so `fetch`/WebSocket calls reach the same origin.
- `web/src/index.js` – SPA entry point. It matches routes, ensures auth state via `/api/is-valid`, bootstraps shared components, and toggles per-view styles.
- `web/src/views` – view classes (Home, Post, Create Post, Auth screens, Chat, Sidebar, Navbar, Contacts) that extend `AbstractView` to share helpers like dynamic CSS loading and category caching.
- `web/src/pkg` – browser utilities: `fetcher.js` standardizes API calls and redirects on 401/403, `Utils.js` handles toasts/local storage, and `WS.js` wraps the chat WebSocket connection.

### Authentication & sessions

1. `POST /api/signup` validates incoming payloads (`pkg/utils/validatorUser.go`) and stores credentials hashed via bcrypt.
2. `POST /api/signin` compares hashes, issues a JWT using `pkg/smplJwt`, and persists it in the `sessions` table before setting the HTTP-only cookie.
3. Every protected route uses `identify` middleware to pull the `token` cookie, ensure it exists in the DB, parse & validate it, and inject the user ID/token into the context for downstream handlers.

### Database & migrations

- Configure the driver (`sqlite3`), database file, and schema directory in `config/config.json`.
- On startup, `pkg/database.ConnectSqlte` opens/creates the `.db` file, pings it, then iterates over `schemas/*.sql` to run each migration in lexical order.
- Seed files (e.g., `10_insert_categories.sql`, `11_insert_posts.sql`, `12_insert_post_categories.sql`) give the SPA data to render immediately.

### Real-time chat & contacts

- `/ws` upgrades authenticated HTTP requests to a WebSocket handled in `internal/controller/messages.go`.
- Events: `new_user`, `msg`, `msg-error`, `typing`, `user-online`, `user-offline`, `error`.
- Messages persist through `MessagesService.CreateMessage`, which also updates the chat's `last_msg` timestamp used to sort contacts.
- The front-end listens for DOM-level custom events triggered by `WS.js` and surfaces toasts if a new DM arrives outside the chat view.

## API Overview

| Method | Path | Description |
| ------ | ---- | ----------- |
| `GET` | `/api/is-valid` | Checks whether the auth cookie is still valid. |
| `POST` | `/api/signup` | Register a new account. |
| `POST` | `/api/signin` | Authenticate and receive the JWT-backed cookie. |
| `POST` | `/api/signout` | Destroy the active session. |
| `GET` | `/api/posts/{category}?limit=&offset=` | Fetch posts within a category (including `General`, `my-posts`, etc.). |
| `GET` | `/api/post/{id}` | Retrieve a single post with comments and vote metadata. |
| `POST` | `/api/post/create` | Create a new post (title, text, up to 5 categories). |
| `POST` | `/api/post/vote` | Like/dislike a post; sending the same vote twice removes it. |
| `POST` | `/api/comment/create` | Add a comment to a post. |
| `POST` | `/api/comment/vote` | Like/dislike a comment. |
| `GET` | `/api/profile/posts/{userID}` | Posts authored by the given user. |
| `GET` | `/api/profile/liked-posts/{userID}` | Posts the user liked. |
| `GET` | `/api/profile/disliked-posts/{userID}` | Posts the user disliked. |
| `GET` | `/api/categories` | List of categories exposed in the sidebar. |
| `GET` | `/api/contacts` | Contacts list enriched with online state. |
| `GET` | `/api/chat/{userID}?limit=&offset=` | Ensures a chat exists with the target user and returns paginated messages. |
| `GET` | `/ws` (WebSocket) | Real-time direct messages, typing events, and presence signals. |

> The handlers apply an IP-based rate limiter (60 requests/min by default); the SPA displays a toast when `429` responses arrive.

## Getting Started

### Prerequisites

- Go 1.22 or newer (CGO must be enabled because `github.com/mattn/go-sqlite3` compiles C code).
- GCC/Clang toolchain and SQLite dev headers (`libsqlite3-dev`) when running outside Docker.
- Docker (optional) if you prefer containerized execution.

### Local setup

```bash
git clone <repo-url>
cd real_time_forum

# Install Go dependencies
go mod download

# Adjust host/port or database path when needed
cp config/config.json config/config.local.json  # optional backup/edit step

# Run the server (serves the SPA at http://localhost:8080)
go run .
```

The first run creates `forum.db`, applies every migration under `schemas`, and writes logs to `logfile.log`. Static assets and the SPA are served from `/web` so you do not need a separate Node build step.

### Configuration

`config/config.json`

```json
{
  "api": {
    "host": "localhost",
    "port": "8080"
  },
  "database": {
    "driver": "sqlite3",
    "fileName": "./forum.db",
    "schemeDir": "./schemas"
  }
}
```

- Update `port` if you want to run multiple instances.
- `fileName` can be an absolute path if you store the DB elsewhere.
- Point `schemeDir` to a different migrations folder if you need custom schemas.

### Docker workflow

Use the provided `Dockerfile` to bundle the backend and static assets.

```bash
# optional: the builder script prunes Docker cache and rebuilds the image
./builder.sh            # WARNING: runs `docker system prune -a`

# manual steps
docker build -t forum-img .
docker run -dp 8080:8080 --name forum forum-img

# attach to the running container for debugging
docker exec -it forum /bin/bash
```

The container exposes port `8080` and runs the compiled `Forum` binary created via `go build` inside the image.

## Front-end & Development Notes

- The SPA uses the History API, so hitting refresh on nested routes works because Go falls back to serving `index.html` and the client router rehydrates.
- `AbstractView` keeps category lists cached so multiple components do not flood `/api/categories`.
- Individual views add/remove scoped CSS files via `data-view-style` attributes to avoid leaking styles between screens.
- `fetcher.js` unwraps responses, redirects unauthorized users to `/sign-in`, and surfaces descriptive errors via `Utils.showError`.
- `WS.js` listens for DOM-level custom events (`send-msg`, `typing`) so views emit simple events instead of dealing with raw sockets.

## Directory Reference

```
.
├── internal
│   ├── app/               # Application bootstrap
│   ├── controller/        # HTTP handlers, router, middleware, WebSocket hub
│   ├── entity/            # Data transfer structs shared across layers
│   ├── repository/        # Persistence layer (SQLite queries)
│   ├── service/           # Business logic per domain
│   └── server/            # Thin HTTP server wrapper
├── pkg
│   ├── config/            # Loads config.json
│   ├── database/          # SQLite connection + migrations
│   ├── smplJwt/           # Custom JWT implementation
│   └── utils/             # Registration validators
├── schemas/               # SQL migrations & seed data
├── web/                   # Static assets, SPA entry, views, css
├── Dockerfile             # Container image definition
└── builder.sh             # Helper script for pruning/building Docker images
```

## Troubleshooting

- **CGO errors when building** – ensure a compiler toolchain is installed (e.g., `sudo apt install build-essential libsqlite3-dev`).
- **`429 Too Many Requests`** – the built-in `RateLimiter` throttles by IP; raise the limit in `internal/controller/rateLimiter.go` or reduce noisy clients.
- **Migrations not running** – confirm `config.database.schemeDir` points to the `schemas` folder relative to the working directory.
- **WebSocket disconnects** – the SPA listens for a custom `ws-closing` event; if you see frequent disconnects check server logs in `logfile.log` or inspect network proxies that might be closing idle connections.

## Next Steps

- Add automated tests around repositories/services (table-driven tests backed by an in-memory SQLite DB).
- Expose pagination metadata in REST responses to simplify infinite scroll on the client.
- Consider extracting the chat hub into its own package if you plan to expand into group conversations.

Enjoy building on top of the forum!

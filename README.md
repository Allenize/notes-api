# Notez

A premium-feeling notes web app with real authentication — built in Go, using only the standard library. Designed to run as a backend behind [Lime](https://github.com/Allenize/lime).

## Features

- **A real web UI** — sidebar navigation, dashboard, note cards, and a Markdown editor, all served same-origin at `/` (no CORS setup needed), styled with a warm cream/yellow/orange/brown palette
- **Dashboard** — quick stats, pinned/favorite/recent notes, and an activity timeline (all backed by real data, not placeholders)
- **Categories & tags** — organize notes and filter by either from the sidebar, with live counts
- **Favorites, pinning, archive, and trash** — with restore from trash, and a bulk "Empty Trash" action
- **Markdown editor** — write/preview toggle, supporting headings, bold/italic, checklists, code blocks, blockquotes, and tables — hand-written renderer, no external Markdown library, with HTML-escaping to prevent script injection from note content
- **Word count & reading time** — computed live as you type
- **JWT authentication** — hand-implemented HMAC-SHA256 signed tokens (no external JWT library)
- **Secure password hashing** — PBKDF2 with a random salt per user, implemented from stdlib crypto primitives (no external hashing library)
- **Access + refresh tokens** — short-lived access tokens (15 min) and long-lived refresh tokens (7 days), with a `/auth/refresh` endpoint
- **Per-user data isolation** — every note is scoped to its owner; users can't see or modify each other's data
- **Persistent storage** — a JSON file-backed store that survives restarts (see the caveat below for hosted deployments)
- **Pagination & filtering** — `limit`, `offset`, `done`, `category`, `tag`, `favorite`, `pinned`, `view`, and `q` (search) query params
- **Rate limiting** — a token-bucket limiter per client IP, implemented from scratch
- **Structured request logging** — every request gets a unique ID, logged with method, path, status, and duration
- **Graceful shutdown** — catches `SIGTERM`/`SIGINT` and finishes in-flight requests before exiting
- **Real test coverage** — unit tests for password hashing, JWT issuing/verification, and the store (including favorites/pin/archive, trash/restore, and a persistence-across-restart test)

Zero external dependencies — everything above, including the Markdown renderer, is built from Go's standard library and vanilla JavaScript.

## Not in this pass

The original design brief also called for calendar view, real-time sharing/collaboration, file attachments with drag-and-drop, and version history. Each of those needs real infrastructure this project doesn't have yet (e.g. attachments need cloud storage, since Render's free tier wipes local files on redeploy) — they're natural next steps, not implemented here.

## Endpoints

| Method | Path                | Auth required | Description                          |
|--------|---------------------|:--------------:|----------------------------------------|
| GET    | `/`                 | no              | Web UI (signup/login + notes dashboard) |
| GET    | `/health`           | no              | Health check                          |
| POST   | `/auth/signup`      | no              | Create an account, returns tokens     |
| POST   | `/auth/login`       | no              | Log in, returns tokens                |
| POST   | `/auth/refresh`     | no              | Exchange a refresh token for a new pair |
| GET    | `/tasks`            | yes             | List your notes (supports pagination/filtering) |
| POST   | `/tasks`            | yes             | Create a note                         |
| GET    | `/tasks/stats`      | yes             | Summary counts (total/favorites/pinned/archived/trashed/categories) |
| GET    | `/tasks/{id}`       | yes             | Get a single note                     |
| PUT    | `/tasks/{id}`       | yes             | Update a note (title/notes/category/tags/favorite/pinned/archived) |
| DELETE | `/tasks/{id}`       | yes             | Permanently delete a note             |
| POST   | `/tasks/{id}/trash` | yes             | Soft-delete a note (move to trash)    |
| POST   | `/tasks/{id}/restore` | yes           | Restore a note from trash             |

Authenticated requests need `Authorization: Bearer <access_token>`.

### Task list query params

- `view=all` (default) / `archived` / `trash` — which bucket to list
- `done=true` / `done=false` — filter by completion status
- `favorite=true` / `pinned=true` — filter by favorite/pinned status
- `category=Name` — filter by category (case-insensitive)
- `tag=name` — filter by tag (case-insensitive)
- `q=search+term` — search title and note body
- `limit=20` — page size (default 20, max 100)
- `offset=0` — page offset

## Environment variables

| Variable      | Description                                  | Default          |
|---------------|------------------------------------------------|--------------------|
| `PORT`        | Port to listen on                              | `9001`             |
| `JWT_SECRET`  | Secret used to sign tokens                     | random per run (⚠️ set this explicitly in production, or tokens break on every restart) |
| `DATA_FILE`   | Path to the persistent JSON data file          | `./data.json`      |

## Run locally

```bash
JWT_SECRET="something-long-and-random" go run ./cmd/notes-api
```

## Run tests

```bash
go test ./... -v
```

## Run with Docker

```bash
docker build -t notes-api .
docker run -p 9001:9001 -e JWT_SECRET="something-long-and-random" notes-api
```

## Deploying alongside Lime

1. Deploy this repo the same way as Lime (Render → New Web Service → select this repo → Free instance)
2. Set the `JWT_SECRET` environment variable to something long and random
3. Copy the resulting URL (e.g. `https://notes-api-xxxx.onrender.com`)
4. On your **Lime** service, set its `BACKENDS` environment variable to that URL
5. Redeploy Lime (or it'll pick it up automatically depending on your Render settings) — it will now health-check and route to this API

### ⚠️ A note on persistence when hosted

This store writes to a local JSON file, which works great for local development and testing. On Render's **free** tier, the filesystem is ephemeral — a redeploy or restart wipes it, so users and tasks won't survive those events in production. For anything beyond a demo, swap the store for a real database (Render offers a free PostgreSQL tier with its own time limits — check their current docs) or attach a persistent disk (a paid feature). The store is isolated behind a small interface-free API in `internal/store`, so swapping the backing implementation later is a contained change.

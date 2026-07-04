# notes-api

A notes/tasks REST API with real authentication — built in Go, using only the standard library. Designed to run as a backend behind [Lime](https://github.com/Allenize/lime).

## Features

- **JWT authentication** — hand-implemented HMAC-SHA256 signed tokens (no external JWT library)
- **Secure password hashing** — PBKDF2 with a random salt per user, implemented from stdlib crypto primitives (no external hashing library)
- **Access + refresh tokens** — short-lived access tokens (15 min) and long-lived refresh tokens (7 days), with a `/auth/refresh` endpoint
- **Per-user data isolation** — every task is scoped to its owner; users can't see or modify each other's data
- **Persistent storage** — a JSON file-backed store that survives restarts (see the caveat below for hosted deployments)
- **Pagination & filtering** — `limit`, `offset`, `done`, and `q` (search) query params on task listing
- **Rate limiting** — a token-bucket limiter per client IP, implemented from scratch
- **Structured request logging** — every request gets a unique ID, logged with method, path, status, and duration
- **Graceful shutdown** — catches `SIGTERM`/`SIGINT` and finishes in-flight requests before exiting
- **Real test coverage** — unit tests for password hashing, JWT issuing/verification, and the store (including a persistence-across-restart test)

Zero external dependencies — everything above is built from Go's standard library.

## Endpoints

| Method | Path             | Auth required | Description                          |
|--------|------------------|:--------------:|----------------------------------------|
| GET    | `/health`        | no              | Health check                          |
| POST   | `/auth/signup`   | no              | Create an account, returns tokens     |
| POST   | `/auth/login`    | no              | Log in, returns tokens                |
| POST   | `/auth/refresh`  | no              | Exchange a refresh token for a new pair |
| GET    | `/tasks`         | yes             | List your tasks (supports pagination/filtering) |
| POST   | `/tasks`         | yes             | Create a task                         |
| GET    | `/tasks/{id}`    | yes             | Get a single task                     |
| PUT    | `/tasks/{id}`    | yes             | Update a task                         |
| DELETE | `/tasks/{id}`    | yes             | Delete a task                         |

Authenticated requests need `Authorization: Bearer <access_token>`.

### Task list query params

- `done=true` / `done=false` — filter by completion status
- `q=search+term` — search title and notes
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

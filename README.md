# 📺 youtube-api

> A tiny, fast, dependency-light Go service that serves the latest **videos** and **playlists** from an allow-list of YouTube channels — with caching baked right into the binary. No Redis, no database, no fuss. 🪶

[![CI](https://github.com/ermos/youtube-api/actions/workflows/ci.yml/badge.svg)](https://github.com/ermos/youtube-api/actions/workflows/ci.yml)
[![Release](https://github.com/ermos/youtube-api/actions/workflows/release.yml/badge.svg)](https://github.com/ermos/youtube-api/actions/workflows/release.yml)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](./LICENSE)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev)

---

## ✨ Features

- 🚀 **Single static binary** — build it, ship it, run it. Nothing else to deploy.
- 🧠 **In-memory TTL cache** — responses are cached right inside the process. Zero external infrastructure.
- 🔐 **Channel allow-list** — only the channels you explicitly trust are served.
- 🌐 **CORS-ready** — drop it behind any frontend.
- 🩺 **Health endpoint** — `/health` for your load balancer.
- 🛡️ **Hardened CI** — tests, linting, and security scanning on every change.

## 🚦 Endpoints

| Method | Path             | Query params                          | Description                       |
| ------ | ---------------- | ------------------------------------- | --------------------------------- |
| `GET`  | `/api/videos`    | `channel_id` *(required)*, `max_results` | Latest videos of a channel     |
| `GET`  | `/api/playlists` | `channel_id` *(required)*, `max_results` | Latest playlists of a channel  |
| `GET`  | `/health`        | —                                     | Liveness probe (`200 OK`)         |

`max_results` defaults to `10` and is capped at `50`. Channels not in the allow-list get a `403`. 🙅

### Example

```bash
curl "http://localhost:8080/api/videos?channel_id=UC_x5XG1OV2P6uZZ5FSM9Ttw&max_results=5"
```

## ⚙️ Configuration

All config is via environment variables (a local `.env` is auto-loaded). See [`.env.example`](./.env.example).

| Variable                   | Required | Default          | Description                                  |
| -------------------------- | -------- | ---------------- | -------------------------------------------- |
| `PORT`                     | no       | `8080`           | HTTP listen port                             |
| `YOUTUBE_API_KEY`          | **yes**  | —                | YouTube Data API v3 key                      |
| `YOUTUBE_ALLOWED_CHANNELS` | **yes**  | —                | Comma-separated channel IDs to serve         |
| `CACHE_TTL_SECONDS`        | no       | `86400` (1 day)  | How long responses stay cached in memory     |

## 🏃 Run it

### From source

```bash
cp .env.example .env   # then fill in YOUTUBE_API_KEY & YOUTUBE_ALLOWED_CHANNELS
go run .
```

### With Docker 🐳

```bash
docker build -t youtube-api .
docker run --rm -p 8080:8080 --env-file .env -e PORT=8080 youtube-api
```

Pre-built images are published to GHCR on every release:

```bash
docker run --rm -p 8080:8080 --env-file .env ghcr.io/ermos/youtube-api:latest
```

## 🧪 Development

```bash
go test -race ./...   # run tests
go vet ./...          # vet
```

## 🤖 CI/CD

- **`ci.yml`** runs on every PR/branch: 🧪 tests (with `-race` + coverage), 🎯 `golangci-lint`, and 🔒 security scans (`govulncheck`, `gosec`, `Trivy`).
- **`release.yml`** runs on `main`: it reads your [Conventional Commits](https://www.conventionalcommits.org/) and automatically bumps **SemVer** (`feat:` → minor, `fix:` → patch, `BREAKING CHANGE` → major), tags the release, builds the Docker image, and pushes it to GHCR. 🏷️

> 💡 Commit with intent — `feat: add shorts endpoint` ships a new minor version all by itself.

## 📦 Project layout

```
main.go               # wiring: config → cache → youtube → router
internal/config       # env-based configuration
internal/cache        # in-memory TTL cache (the "Redis killer")
internal/youtube      # YouTube Data API client
internal/handlers     # HTTP handlers
```

## 📜 License

Licensed under the **GNU Affero General Public License v3.0** — see [`LICENSE`](./LICENSE).

This means you're free to use, study, share, and improve it. If you run a modified version as a network service, you must make your source available to its users. ❤️ Keep it open.

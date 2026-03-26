---
id: contribute-code
sidebar_position: 1
title: Contribute Code
description: Draft workflow for shipping code changes to Colibri.
---

Colibri welcomes pull requests, especially when they arrive with context. Use this checklist as a starting point and tailor it to your needs; the goal is to make reviewers productive and deployments safe.

## 1. Explore the project structure

- `cmd/api` holds the HTTP server bootstrap.
- `internal/server` wires middleware and routes.
- `internal/{posts,sources,categories}` contain handlers and helpers.
- `internal/database` is generated with `sqlc` and mirrors the Postgres schema.

Spend a few minutes reading the files around the functionality you plan to change. When in doubt, search for existing patterns and mimic them.

## 2. Set up your local environment (Docker Compose)

1. Duplicate `.env.example` to `.env`. The Go services read `DB_DRIVER=postgres` and the `*_FILE` variables that point to Docker secrets, so the default values are fine.
2. Create the secrets expected by `docker-compose.yml`:
   - `secrets/db-password.txt`: plain-text password for the local Postgres container.
   - `secrets/db-string.txt`: full connection string, e.g. `postgresql://postgres:<password>@postgres:5432/colibri_db?sslmode=disable` (note that the host is the Compose service name `postgres`).
   - `secrets/rabbitmq-url.txt`: e.g. `amqp://guest:guest@rabbitmq:5672/`.
   - `secrets/google-api-key.txt`: add a cheap google api key for summarization.
3. Bring up the stack: `docker compose up --build. The fetcher will run once every 4 hours by default, you can change that in the docker-compose.yml or simply restart your container to trigger it.

## 3. Follow the coding conventions

- Stick to Go 1.22+ features, avoid third-party helpers when the standard library suffices.
- Handlers return JSON and keep response shapes in `internal/{domain}/{domain}.go` structs.
- Prefer unit tests that work without external services.
- Keep middleware minimal. Rate limiting, security headers, and CORS are already configured in `internal/server/server.go`.

## 4. Validate your change

1. Run `go test ./...`.
2. Run `go fmt ./...`.
3. Run `staticcheck ./...`.
4. Run `gosec ./...`.
5. If you touched SQL, run `sqlc generate` and ensure no unexpected diffs remain.
6. Update `docs/api/openapi.yaml` when you add or modify endpoints, then regenerate docs with `npm --prefix docs run gen-api`.
7. Run `npm --prefix docs run build` to confirm the documentation site still compiles.

## 5. Prepare the pull request

- Summarize the change.
- Link related issues or discussions.
- Include screenshots or API samples when you change user-facing behavior.
- Keep commits focused. Squash if you did multiple commits.

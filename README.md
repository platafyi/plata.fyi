# plata.fyi

[![Deploy](https://github.com/platafyi/plata.fyi/actions/workflows/deploy.yml/badge.svg)](https://github.com/platafyi/plata.fyi/actions/workflows/deploy.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Anonymous salary sharing for Macedonian workers. No accounts, no passwords, no email — fill out a form, pass a Cloudflare 
Turnstile check, get a session token stored in your browser. All data is public and anonymous.

Every deployment automatically publishes a full CSV snapshot of all salary data as a [GitHub Release](../../releases). Free to download and use.

Inspired by [levels.fyi](https://levels.fyi).

---

## Stack

- **Backend** — Go 1.26+
- **Frontend** — Next.js 14 (App Router), TypeScript, Tailwind CSS
- **Database** — PostgreSQL 16
- **Bot protection** — Cloudflare Turnstile
- **Currency** — MKD only (Macedonian Denar)

---

## Local development

```bash
make dev-up       # start Postgres in Docker
make dev-migrate  # run all migrations

# in separate terminals:
cd backend && go run ./cmd/api
cd frontend && npm run dev
```

Copy `.env.example` to `.env` and fill in the required values.

Turnstile verification is skipped when `TURNSTILE_SECRET` is empty (dev mode).

---

## How anonymity works

Privacy is a core design constraint, not an afterthought.

### No email, no accounts

No email address is ever collected or stored. Authentication is replaced by a session token generated server-side 
after successful Turnstile verification. The token lives in `localStorage` for 30 days. 
If someone loses it (new device, cleared browser data), they start a fresh session 
old submissions are unrecoverable by design.

### Session tokens

On first submission, the backend generates a `randHex(32)` session token, stores it in the `tokens` table, 
and returns it to the browser. Subsequent requests authenticate with `Authorization: Bearer <token>`. 
Tokens expire after 30 days. Currently, there's no way to refresh a session. 

### Dates are truncated to month

`created_at` and `updated_at` on salary submissions are stored at month precision only 
(e.g. `2026-04-01`, not `2026-04-23 14:32:05`).

---

## Running tests

```bash
cd backend
go test -race ./internal/...
```

Integration tests against a real database (skipped automatically if no DB is available):

```bash
TEST_DB_URL=postgres://platafyi:platafyi@localhost:5433/platafyi?sslmode=disable \
  go test -v -run TestSearch ./internal/database/
```

---

## Migrations

Migrations are plain `.sql` files in `backend/migrations/`, numbered sequentially (`001_`, `002_`, …). 
The custom runner in `backend/cmd/migrate/` applies them in order and tracks applied files in a `schema_migrations` table, 
so each file runs exactly once.

Migrations are only forward (up). No backwards (down) migrations allowed. 

To add a migration:

1. Create `backend/migrations/NNN_description.sql` (increment the prefix).
2. Run `make dev-migrate` it skips already-applied files and applies new ones inside a transaction. 
3. If a migration fails, the transaction is rolled back and the run stops.

In production, migrations are applied automatically by the deploy pipeline before the new image rolls out. See [deploy.yml](.github/workflows/deploy.yml).

---

## Contributing

1. **Clone** the repo and create a branch from `main`.
2. **Set up locally.** Follow the [Local development](#local-development) steps above.
3. **Make your changes.** Keep PRs focused; one concern per PR.
4. **Run tests** before opening a PR:
   ```bash
   cd backend && go test -race ./internal/...
   cd frontend && npm run lint
   ```
5. **Open a pull request** against `main` with a clear description of what and why.

Please use [Conventional Commits](https://www.conventionalcommits.org/) for your commit messages, e.g. 
`feat: add company search`, `fix: handle missing token`, `docs: update contributing guide`.

A few things to keep in mind:

- Please include tests for backend changes. Existing tests should keep passing.
- Keep the API and database schema backwards compatible. Add new optional fields rather than renaming or removing existing ones.
- No personally identifiable information should ever be stored or logged. Anonymity is a core invariant.
- New database columns require a numbered migration file; never modify existing migration files.
- **A little copying is better than a little dependency** ([why](https://www.youtube.com/watch?v=PAAkCSZUG1c&t=9m28s)). 
If you do add one, run `go mod tidy && go mod vendor`.
- The project targets Go 1.26+ and Next.js 14.

---

## Deployment

Deployment is fully automated via [GitHub Actions](.github/workflows/deploy.yml). Every push to `main`:

1. **Builds** Docker images for the API and frontend and pushes them to GitHub Container Registry (`ghcr.io/platafyi/api`, `ghcr.io/platafyi/frontend`) 
tagged with the commit SHA.
2. **Backs up** the live database by exporting an anonymous CSV snapshot and publishing it as a GitHub Release.
3. **Deploys** by applying Kubernetes manifests via `kubectl apply -k k8s/` and waiting for the rollout to complete.

Manifests are in `k8s/` (managed with Kustomize). The deploy job patches the image tag to the current commit SHA before applying, 
so each deploy is pinned to an exact build.

Required GitHub secrets: `KUBECONFIG` (base64-encoded kubeconfig), `TURNSTILE_SECRET`, `DB_URL`.
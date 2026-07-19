# widget-srv Roadmap

Intentionally a **toy service for instruction/demos** — favor simplicity and teaching value
over production hardening. Working style: guided / pair-programming, one stage at a time.
Check items off as we go.

The **core path** (Stages 1–5) is the main learning arc. Everything under "Optional" is only
worth doing if a specific demo/lesson calls for it — not defaults.

## Done

- [x] **Stage 1 — DB + migrations**
  - golang-migrate with `//go:embed`ed SQL, runs on startup (idempotent, dirty-flag aware)
  - `product` + `inventory` tables; 3 seeded widgets
  - files: `internal/db/migrate.go`, `internal/db/migrations/*.sql`
- [x] **Stage 2 — pgxpool + minimal Gin server**
  - pool opened + startup ping (fail fast); `GET /healthz` (503 when DB down)
  - graceful shutdown on SIGINT/SIGTERM via `http.Server.Shutdown` (10s drain)
  - file: `main.go`

## In progress

- [ ] **Stage 3 — Containerize into compose project**
  - multi-stage `widget-srv/Dockerfile` (build in golang image, run on tiny base)
  - `widget-srv` service in `compose.yml`: `depends_on` postgres healthcheck,
    `DATABASE_URL` pointed at the `postgres` service name (not localhost)
  - optional: `.dockerignore`

## Upcoming

- [ ] **Stage 4 — sqlc + real endpoints**
  - `sqlc.yaml` (engine=postgresql, sql_package=pgx/v5), queries in `internal/db/query/widgets.sql`
  - generate type-safe code into `internal/db/gen/`
  - `GET /products` (list), `GET /products/:id` (with inventory count)
- [ ] **Stage 5 — update endpoints**
  - `PUT /products/:id` (name + attributes), `PUT /products/:id/inventory` (set count)
  - refactor routes/handlers into `internal/api/`

## Optional — only if a demo/lesson wants it

Nice for showing off a full stack or specific concepts, but not required for the core arc.

- [ ] **Flutter shop UI** — public browse/detail view (consumes `GET /products`, `GET /products/:id`)
- [ ] **Flutter admin panel** — manage products/inventory (consumes the `PUT` endpoints;
      would need create/delete endpoints added too)
  - If we build the UIs: CORS on widget-srv, per-env API base-URL, and (for admin) some auth.
- [ ] Dev tooling: `Makefile` (build / vet / lint / test / sqlc / run / compose up-down)
      and a `golangci-lint` config
- [ ] Round-trip migration test (up → down 1 → up) to catch bad `down` files

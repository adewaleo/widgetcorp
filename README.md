# WidgetCorp

A teaching monorepo for distributed-systems concepts, built around a fictional
online widget store. Each service is a small, focused teaching device; wired
together they demonstrate observability, CDC, durable workflows, caching,
queueing, search, and more.

Start here:

- **[docs/concepts.md](docs/concepts.md)** — what we demonstrate (and the article backlog)
- **[docs/architecture.md](docs/architecture.md)** — system design
- **[docs/repo-structure.md](docs/repo-structure.md)** — layout, builds, CI/CD
- **[docs/roadmap.md](docs/roadmap.md)** — roadmap (mirror of the GitHub Project)

## Prerequisites

- [Task](https://taskfile.dev) — the task runner (`brew install go-task/tap/go-task`
  or see the install docs)
- Docker + Docker Compose
- Language toolchains as you work on services: Go, Python 3.11+

## Getting started

```bash
task            # list everything you can run
task up         # start the core stack (once services exist)
task ps         # see what's running
task down       # tear down
```

Tasks are namespaced per service as they're built — e.g. `task catalog:build`,
`task catalog:test`. The root `Taskfile.yml` fans out to each service's own
Taskfile via `includes`.

## Repository layout

```
services/     one folder per deployable (+ _templates/ skeletons)
platform/     infra config (kafka, observability, temporal, lake, cache, queue, postgres)
deploy/       docker-compose files + profiles
contracts/    OpenAPI / proto / event schemas (contracts-first)
tools/        loadgen, seed data, scripts
docs/         concepts, architecture, roadmap, demo runbooks
```

See [docs/repo-structure.md](docs/repo-structure.md) for the full tree and the
per-service contract every service follows.

## Status

Early scaffolding (Phase 0). Follow progress in [docs/roadmap.md](docs/roadmap.md).

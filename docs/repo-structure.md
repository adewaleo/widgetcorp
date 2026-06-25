# Repository Structure, Builds & CI/CD

WidgetCorp is a **monorepo**. Polyglot services, shared infra config, contracts,
and tooling all live together so a reader can see the whole system at once.

## Directory layout

```
widgetcorp/
├── README.md                     # what this is, quickstart, demo index
├── LICENSE
├── Makefile                      # top-level entrypoint (up/down/seed/load/demo-*)
│
├── docs/
│   ├── concepts.md               # what we demo + article backlog
│   ├── repo-structure.md         # this file
│   ├── architecture.md           # system design
│   ├── concepts/                 # one explainer per concept
│   │   ├── observability.md
│   │   ├── cdc-debezium.md
│   │   ├── temporal.md
│   │   └── monolith-vs-micro.md
│   └── demos/                    # step-by-step runbooks
│       ├── observability-tour.md
│       ├── cdc-no-db-load.md
│       └── temporal-resilience.md
│
├── contracts/                    # contracts-first; the stable artifacts
│   ├── openapi/                  # *.yaml per service (today)
│   ├── proto/                    # *.proto (added when we demo gRPC)
│   └── events/                   # event schemas (outbox/Kafka payloads)
│
├── services/                     # one folder per deployable
│   ├── _templates/               # copy-me skeletons (go/, python/)
│   ├── catalog/      (Go)
│   ├── inventory/    (Go)
│   ├── order/        (Python/FastAPI)
│   ├── payment/      (Python)
│   ├── fulfillment/  (Temporal worker)
│   ├── consumers/    (search indexer, lake sink glue)
│   ├── bff/          (Python — backend-for-frontend)
│   └── web/          (simple UI for visualizing services)
│
├── platform/                     # infra config, not application code
│   ├── observability/            # otel-collector, data-prepper, dashboards
│   ├── kafka/                    # broker (KRaft), connect, debezium connectors
│   ├── queue/                    # rabbitmq or nats (work-queue demos)
│   ├── cache/                    # valkey
│   ├── temporal/                 # server + UI config
│   ├── lake/                     # minio, iceberg catalog (nessie), trino
│   └── postgres/                 # init scripts, schema-per-service
│
├── deploy/
│   ├── docker-compose.yml        # base (core services + postgres)
│   ├── compose.observability.yml
│   ├── compose.cdc.yml
│   ├── compose.lake.yml
│   ├── compose.temporal.yml
│   ├── compose.cache.yml
│   ├── compose.queue.yml
│   └── .env.example
│
├── tools/
│   ├── loadgen/                  # Locust scenarios (browse, order, burst)
│   ├── seed/                     # synthetic catalog/inventory data
│   └── scripts/                  # healthcheck waits, helpers
│
└── .github/
    └── workflows/                # CI/CD pipelines
```

## Per-service contract

Every service — regardless of language — ships the same surface so the polyglot
story stays consistent and adding a language later is mechanical:

| Item | Why |
|---|---|
| `Dockerfile` (multi-stage) | reproducible image, small runtime layer |
| `/healthz` + `/readyz` | compose healthchecks + ordered startup |
| OTel instrumentation | traces/logs/metrics out of the box |
| DB migrations | schema is versioned, not hand-applied |
| `Makefile` / task runner | uniform `build` / `test` / `lint` / `run` targets |
| tests | unit + a thin integration smoke test |
| `README.md` | what it does, endpoints, how to run alone |

The `_templates/go` and `_templates/python` skeletons encode all of the above.

## Build strategy

- **Each service builds independently** via its own multi-stage `Dockerfile`.
  No service depends on another's build.
- **Uniform Make targets** per service: `make build test lint image`.
- **Top-level Makefile** fans out: `make build-all`, `make test-all`, plus
  demo orchestration (`make up`, `make demo-cdc`, `make load`).
- **Compose profiles** assemble subsystems on demand:
  `docker compose --profile core up`, then add `--profile observability`, etc.
- **Pinned versions** for all infra images (Kafka, OpenSearch, Trino, Temporal…)
  so demos are reproducible.

## CI (GitHub Actions)

Triggered on PRs and pushes to `main`:

1. **Change detection** — only build/test services whose paths changed
   (path filters / matrix), keeping CI fast in a monorepo.
2. **Per-service matrix job**: lint → unit test → build image.
3. **Contract checks** — validate OpenAPI specs; (later) lint protobufs;
   verify code matches contracts.
4. **Integration smoke** — `compose --profile core up`, run health + a
   single end-to-end order, tear down.
5. **Security/quality** — dependency scan, Dockerfile lint (hadolint),
   secret scan.

## CD

- **Image publishing** — on merge to `main`, push tagged images to **GHCR**
  (`ghcr.io/<org>/widgetcorp-<service>:<sha>` + `:latest`).
- **Versioning** — tag images by git SHA; optional semver on release tags.
- Since this is a demo repo, "deploy" = images are pullable so anyone can run
  the stack without building locally. A future phase can add a k8s/Helm or
  Compose-on-a-host deploy target.

## Conventions

- **Trunk-based**, short-lived branches, PRs into `main`.
- **Conventional commits** (enables changelogs / semver later).
- One service = one folder = one image = one CI matrix entry.
- Infra config lives in `platform/`; how to run it lives in `deploy/`.

# Architecture

WidgetCorp is a fictional online widget store. The system is intentionally simple
per service — each one is a teaching device, not a product — but wired together so
that real distributed-systems concepts (tracing, CDC, durable workflows) emerge
naturally.

See [concepts.md](concepts.md) for what we demonstrate and
[repo-structure.md](repo-structure.md) for how the repo is laid out.

## Design principles

1. **Each service teaches one concept** where possible.
2. **Contracts-first** — OpenAPI is the stable artifact; transport (REST → gRPC →
   events) can evolve without rewrites.
3. **Evolvable by construction** — the outbox + CDC path exists from day one, so
   moving from REST to event-driven is additive.
4. **Simple to run** — `docker compose` with profiles; boot only what a demo needs.
5. **Open source only** — see the tooling table below.

## System overview

```
                          ┌──────────────┐
        shopper ─────────▶│   Web UI     │  (simple SPA: visualize services)
                          └──────┬───────┘
                                 │ REST/JSON
                          ┌──────▼───────┐
                          │     BFF      │  (Python — aggregates for the UI)
                          └──────┬───────┘
              ┌──────────────┬───┴──────────┬───────────────┐
              ▼              ▼               ▼               ▼
        ┌──────────┐  ┌───────────┐   ┌──────────┐   ┌──────────┐
        │ Catalog  │  │ Inventory │   │  Order   │──▶│ Payment  │
        │  (Go)    │  │   (Go)    │   │ (Python) │   │ (Python) │
        └────┬─────┘  └─────┬─────┘   └────┬─────┘   └────┬─────┘
             │ schema       │ schema       │ schema       │ schema
             └──────────────┴──────┬───────┴──────────────┘
                                   ▼              Order also starts the
                          PostgreSQL                fulfillment workflow:
                       (schema-per-service)       ┌─────────────────────┐
                                   │  WAL         │ Temporal (worker)   │
                                   ▼              │ reserve→charge→ship │
                        Debezium (CDC, reads WAL) └─────────────────────┘
                                   ▼
                                 Kafka
                          ┌────────┴─────────┐
                          ▼                  ▼
                  Search indexer       Iceberg sink
                  → OpenSearch         → MinIO/Iceberg → Trino (SQL)

   Observability (cross-cutting): every service → OTel Collector
                                  → Data Prepper → OpenSearch + Dashboards
```

## Services (initial)

| Service | Lang | Responsibility | Key endpoints (sketch) |
|---|---|---|---|
| **Catalog** | Go | Widget products | `GET /widgets`, `GET /widgets/{id}`, `POST /widgets` |
| **Inventory** | Go | Stock levels (the "hot table") | `GET /stock/{id}`, `POST /stock/{id}/reserve` |
| **Order** | Python | Order intake; orchestrates purchase | `POST /orders`, `GET /orders/{id}` |
| **Payment** | Python | Mock payments w/ injectable failures | `POST /charges`, `POST /refunds` |
| **BFF** | Python | Aggregates services for the UI | `GET /home`, `POST /checkout` |
| **Web UI** | (simple SPA) | Visualize + drive the services | — |
| **Fulfillment** | Go/Python | Temporal workflow worker | (Temporal, not HTTP) |
| **Consumers** | Python | Search indexer + lake sink glue | (Kafka, not HTTP) |

### BFF + Web UI

- The **BFF** is the single entry point for the UI: it fans out to Catalog,
  Inventory, Order, and Payment and returns UI-shaped responses. It keeps the UI
  dumb and gives us one clean place to originate traces.
- The **Web UI** is a deliberately simple single-page app whose job is to *make
  the system visible*: browse widgets, place an order, and watch state change.
  Where useful it links out to OpenSearch Dashboards / Temporal UI / Trino so a
  viewer can jump from an action to its trace, workflow, or lake row.
- Both are optional for backend-only demos (the load generator can stand in for a
  shopper), but the UI makes live demos far more compelling.

## Data & events

- **PostgreSQL**, one instance, **schema-per-service** to start (documented as
  splittable into a DB-per-service later). Keeps compose light while preserving
  logical isolation.
- **Transactional outbox**: each service writes domain events
  (`OrderPlaced`, `StockReserved`, …) to an outbox table in the same transaction
  as its state change. Debezium captures the outbox off the WAL → Kafka.
- **CDC** therefore drives both the search index and the lake export with **zero
  query load** on the operational tables — the headline demo.
- **Caching**: a **Valkey** cache fronts a read-heavy path (e.g. Catalog reads),
  invalidated by CDC events so it stays fresh without dual writes (see concepts
  #13). Added under the `cache` profile.
- **Queueing**: Kafka is the replayable event *log*; a **work queue**
  (RabbitMQ/NATS) handles task-style messaging like outbound notifications, to
  contrast queue vs. log vs. pub/sub (see concepts #14). Added under the
  `queue` profile.

## Communication (evolution path)

| Phase | Style | Where |
|---|---|---|
| Now | **REST/JSON** | UI→BFF→services, service→service |
| Later | **gRPC** | a service pair, to contrast latency/coupling/contracts |
| Later | **Event-driven** | consumers react to Kafka events instead of calling |
| Later | **Durable workflow** | Temporal owns the multi-step fulfillment saga |

## Cross-cutting infrastructure (all OSS)

| Concern | Tool | License |
|---|---|---|
| Telemetry pipeline | OpenTelemetry Collector | Apache-2.0 |
| Traces/logs/metrics store + UI | OpenSearch + Dashboards | Apache-2.0 |
| OTLP → OpenSearch | Data Prepper | Apache-2.0 |
| Event log / streaming | Apache Kafka (KRaft) | Apache-2.0 |
| Work queue / pub-sub | RabbitMQ or NATS | MPL-2.0 / Apache-2.0 |
| Cache | Valkey (Redis OSS fork) | BSD-3-Clause |
| CDC | Debezium on Kafka Connect | Apache-2.0 |
| Object store | MinIO | AGPL-3.0 |
| Table format | Apache Iceberg | Apache-2.0 |
| Iceberg catalog | Nessie (or Iceberg REST) | Apache-2.0 |
| Lake query | Trino | Apache-2.0 |
| Workflows | Temporal + UI | MIT |
| Load generation | Locust | MIT |
| Database | PostgreSQL | PostgreSQL |

Licensing note: MinIO is AGPL — a conscious choice for a demo repo; everything
else is Apache/MIT.

## Order placement (the canonical flow)

```
UI → BFF → Order.POST /orders
              ├─ start Temporal fulfillment workflow
              │     reserve → charge → ship → confirm
              │     (retries + compensation on failure)
              └─ write OrderPlaced to outbox  ──CDC──▶ Kafka
                                                        ├─ search indexer → OpenSearch
                                                        └─ Iceberg sink → lake
```

Every hop carries the same `trace_id`, so the whole flow is one trace in
OpenSearch Trace Analytics — and every log line links back to it.

## Deployment

- Each service is a Docker image; `docker compose` assembles the stack.
- **Profiles**: `core`, `observability`, `cdc`, `lake`, `temporal`, `cache`,
  `queue` — boot only what a given demo needs (see
  [repo-structure.md](repo-structure.md)).

## Roadmap (high level)

- **Phase 0** — skeleton: monorepo, Makefile, base compose, Catalog end-to-end.
- **Phase 1** — core path: Inventory, Order, Payment, BFF + UI; place an order.
- **Phase 2** — observability: OTel → Data Prepper → OpenSearch; trace tour.
- **Phase 3** — CDC: outbox + Debezium → search index + Iceberg lake; no-load demo.
- **Phase 4** — Temporal: fulfillment workflow with retries/compensation.
- **Phase 5** — contrasts: gRPC, event-driven, monolith-vs-microservices.

# Widget Corp roadmap

> **Generated file — do not edit by hand.** Mirror of the GitHub Project;
> regenerate with `tools/scripts/dump-roadmap.sh`. Source of truth: [Project](https://github.com/users/adewaleo/projects/1).

**0/31 done.**

## P0 (0/6)

- [ ] **[P0] Base docker-compose + .env.example + profiles** — _Todo_
  - Base compose with core services + Postgres; wire compose profiles (core/observability/cdc/lake/temporal/cache/queue). Concept #10.
- [ ] **[P0] CI: path-filtered build/lint/test matrix** — _Todo_
  - GitHub Actions monorepo CI (change detection, per-service matrix, contract checks, integration smoke). Concept #10.
- [ ] **[P0] Catalog service end-to-end (Go)** — _Todo_
  - Catalog API in Go: Postgres schema, migrations, Dockerfile, health checks; first service running end-to-end. Concept #1.
- [ ] **[P0] Go service template** — _Todo_
  - services/_templates/go: Dockerfile, /healthz+/readyz, OTel init, migrations, Makefile, tests. Concept #1.
- [ ] **[P0] Monorepo skeleton + top-level Makefile** — _In Progress_
  - Set up monorepo layout per docs/repo-structure.md and the Makefile entrypoint (up/down/seed/load/demo-*).
- [ ] **[P0] Python service template** — _Todo_
  - services/_templates/python: Dockerfile, /healthz+/readyz, OTel init, migrations, Makefile, tests. Concept #1.

## P1 (0/7)

- [ ] **[P1] BFF (Python)** — _Todo_
  - Backend-for-frontend aggregating Catalog/Inventory/Order/Payment for the UI. Concept #1.
- [ ] **[P1] End-to-end place-an-order flow** — _Todo_
  - Wire UI->BFF->Order->Inventory/Payment so an order completes across services.
- [ ] **[P1] Inventory service (Go)** — _Todo_
  - Stock levels service (the CDC hot table). Concept #1.
- [ ] **[P1] OpenAPI contracts for core services** — _Todo_
  - contracts/openapi specs as the stable artifacts (transport swappable). Concepts #1, #8.
- [ ] **[P1] Order service (Python/FastAPI)** — _Todo_
  - Order intake; orchestrates the purchase; trace origin. Concept #1.
- [ ] **[P1] Payment mock (Python)** — _Todo_
  - Mock payments with injectable latency/failures for resilience demos. Concept #12.
- [ ] **[P1] Simple Web UI** — _Todo_
  - Simple SPA to visualize and drive the services; links out to OpenSearch/Temporal/Trino.

## P2 (0/4)

- [ ] **[P2] OTel instrumentation across services** — _Todo_
  - Add OpenTelemetry traces/logs/metrics to every service; single OTLP endpoint. Concept #2.
- [ ] **[P2] Observability pipeline: Collector + Data Prepper + OpenSearch** — _Todo_
  - OTel Collector -> Data Prepper -> OpenSearch + Dashboards. Concept #2.
- [ ] **[P2] RED metrics + log-trace correlation** — _Todo_
  - RED metrics per service; trace_id baked into every log line. Concept #2.
- [ ] **[P2] Trace tour demo + runbook** — _Todo_
  - Follow one order across services in Trace Analytics; docs/demos runbook. Concept #2.

## P3 (0/5)

- [ ] **[P3] CDC vs naive poller no-DB-load demo + runbook** — _Todo_
  - Headline demo: poller vs CDC under shopper load, on one dashboard. Concept #5.
- [ ] **[P3] Iceberg sink -> MinIO + Nessie + Trino** — _Todo_
  - Land CDC changes as Iceberg tables on MinIO; catalog via Nessie; query with Trino. Concepts #5, #6.
- [ ] **[P3] Kafka + Kafka Connect + Debezium** — _Todo_
  - Kafka (KRaft) + Connect + Debezium capturing changes off the Postgres WAL. Concept #4.
- [ ] **[P3] Search indexer -> OpenSearch (CDC read model)** — _Todo_
  - Build product search as a CDC-driven read model / CQRS projection. Concept #15.
- [ ] **[P3] Transactional outbox in services** — _Todo_
  - Write domain events (OrderPlaced/StockReserved) to outbox tables in-transaction. Concept #3.

## P4 (0/3)

- [ ] **[P4] Fulfillment workflow + compensation** — _Todo_
  - reserve->charge->ship->confirm with retries, timers, and saga compensation. Concept #7.
- [ ] **[P4] Resilience demo: kill worker mid-flight + runbook** — _Todo_
  - Kill the worker and watch the workflow resume; contrast with a hand-rolled state machine. Concept #7.
- [ ] **[P4] Temporal server + UI** — _Todo_
  - Stand up Temporal server + UI under the temporal compose profile. Concept #7.

## P5 (0/6)

- [ ] **[P5] Caching with Valkey + CDC invalidation** — _Todo_
  - Cache a read-heavy path; invalidate via CDC events (cache-aside/read-through). Concept #13.
- [ ] **[P5] Event-driven variant** — _Todo_
  - Consumers react to Kafka events instead of calling REST. Concept #8.
- [ ] **[P5] Locust load generation: steady + burst** — _Todo_
  - Synthetic traffic scenarios incl. burst mode for the CDC load contrast. Concept #11.
- [ ] **[P5] Monolith vs microservices comparison** — _Todo_
  - Same feature set as one deployable vs split services; compare deploy/blast-radius/traces. Concept #9.
- [ ] **[P5] Queueing: alternatives to / alongside Kafka** — _Todo_
  - Compare Kafka vs Redpanda/NATS/RabbitMQ; queue vs log vs pub-sub. Concept #14.
- [ ] **[P5] gRPC variant for a service pair** — _Todo_
  - Serve the same contract over gRPC to contrast latency/coupling/contracts. Concept #8.

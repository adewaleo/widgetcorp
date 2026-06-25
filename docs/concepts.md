# Concepts to Demonstrate

WidgetCorp is a teaching repo built around a fictional online widget store. Each
concept below is something we want to **demo runnably** and **write an article
about**. The domain (catalog, inventory, orders, payments, fulfillment) is just a
believable backdrop that generates the events these concepts need.

Guiding rule: one concept per article, each backed by a `docs/demos/<name>.md`
runbook and the minimum set of services needed to show it.

Legend: 🟢 early / foundational · 🟡 mid · 🔵 later / advanced

---

## 1. APIs in different languages 🟢
- The same service "shape" (health checks, migrations, OTel, Dockerfile, tests)
  implemented in **Go** and **Python**, with **Kotlin** planned later.
- What stays the same vs. what differs across languages (idioms, tooling, build).
- Contracts-first development: OpenAPI as the stable artifact, transport swappable.

## 2. Observability: logs, traces, metrics 🟢
- Instrument every service with **OpenTelemetry**; one OTLP endpoint.
- Pipeline to **OpenSearch** via **Data Prepper**; view in **OpenSearch Dashboards**.
- Follow a single order across services in **Trace Analytics**.
- Correlate logs ↔ traces via `trace_id` baked into every log line.
- RED metrics (Rate, Errors, Duration) per service.
- Article angle: "what each signal type is for, and when each one saves you."

## 3. Transactional outbox pattern 🟡
- Reliably publish domain events (`OrderPlaced`, `StockReserved`) without
  dual-write bugs.
- Sets up the event-driven path even while services still call REST.

## 4. Change Data Capture (CDC) with Debezium 🟡
- Capture changes off the Postgres **WAL** — no query load on the source table.
- Stream changes to **Kafka**.
- Article angle: "CDC vs. polling, and why reading the log beats querying the table."

## 5. CDC → datalake export without overloading the DB 🟡 ⭐
- The headline demo. Naive poller (`SELECT ... WHERE updated_at > ?`) vs. CDC.
- Same dashboard: DB CPU / query latency / replication lag under shopper load.
- Land changes in **Apache Iceberg** on **MinIO**; query with **Trino**.
- Article angle: "decoupling analytics from your operational database."

## 6. Datalake / lakehouse basics 🔵
- Object storage (MinIO/S3), table format (Iceberg), query engine (Trino).
- Optional: **Nessie** for "git-for-data" branching and time-travel.

## 7. Workflow orchestration with Temporal 🟡 ⭐
- Order fulfillment as a durable workflow: reserve → charge → ship → confirm.
- Retries, timers, and **compensation/saga** (release stock, refund) on failure.
- Resilience demo: kill the worker mid-flight, watch it resume.
- Article angle: "Temporal vs. a hand-rolled status-column state machine."

## 8. Inter-service communication tradeoffs 🔵
- REST/JSON (baseline) vs. **gRPC** (same contract, different transport).
- vs. **event-driven** (react to Kafka events instead of calling).
- When to reach for each; latency, coupling, and debuggability tradeoffs.

## 9. Monolith vs. microservices 🔵
- Same feature set as one deployable vs. split services.
- Compare deploy granularity, failure blast radius, and trace shape.

## 10. Packaging & orchestration 🟢
- Each service as a **Docker image**; **docker compose** with **profiles**
  (`core`, `observability`, `cdc`, `lake`, `temporal`, `cache`, `queue`) to
  compose demos.

## 11. Synthetic traffic generation 🟢
- **Locust** scenarios: steady browse + order rate, plus a "burst" mode to drive
  the CDC load-contrast demo.

## 12. Resilience & failure injection 🔵
- Injectable latency/failures in the Payment mock.
- Retries, timeouts, circuit breaking, and how they show up in traces.

## 13. Caching technologies & patterns 🟡
- Add a cache in front of a read-heavy path (e.g. Catalog reads).
- Patterns: cache-aside, read-through, write-through, TTL/eviction.
- Cache invalidation driven by **CDC events** (tie-in with #4) so the cache
  stays fresh without dual writes.
- Compare options: **Valkey** (OSS Redis fork) vs. **Memcached**.
- Article angle: "caching patterns and the two hard problems (naming, invalidation)."

## 14. Queueing & messaging technologies 🟡
- Contrast the messaging models we use:
  - **Log / streaming** — Kafka (replayable, multi-consumer; already in the stack).
  - **Work/task queue** — RabbitMQ or NATS (competing consumers, acks, DLQs).
  - **Lightweight pub/sub** — NATS or Valkey streams.
- Patterns: at-least-once delivery, dead-letter queues, backpressure, idempotency.
- Where each fits: e.g. a queue for outbound notifications/emails vs. the Kafka
  event log for CDC.
- Article angle: "queue vs. log vs. pub/sub — picking the right messaging model."

---

## Candidate article series (rough order)

1. Building the same API in Go and Python
2. Observability from zero with OpenTelemetry + OpenSearch
3. The transactional outbox pattern
4. CDC with Debezium: stop polling your database
5. Exporting to a datalake without hurting your DB ⭐
6. Durable workflows with Temporal ⭐
7. REST vs. gRPC vs. events: choosing a transport
8. Caching patterns with Valkey (and invalidation via CDC)
9. Queue vs. log vs. pub/sub: picking a messaging model
10. Monolith vs. microservices, measured not preached

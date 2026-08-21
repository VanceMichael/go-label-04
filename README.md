# AirBridge Cargo Operations

AirBridge coordinates international air-cargo handoffs: a shipper creates a cargo movement, an operations coordinator assigns it to a flight leg, customs and security checks advance through explicit states, and the ground team confirms loading and departure. The service keeps tenant data isolated, records audit events and publishes durable outbox events for downstream airline systems.

## Running

Start PostgreSQL with `docker compose up -d postgres`, set `DATABASE_URL`, then run `go run ./cmd/server`. Migrations are applied at startup. The API exposes `/healthz`, `/readyz`, `/v1/auth/login`, `/v1/shipments`, `/v1/flight-legs`, `/v1/customs-cases`, and `/v1/operations/summary`.

Roles are `shipper`, `coordinator`, and `ground_agent`. Sessions are opaque, server-side, expire by TTL, and are revoked on logout or user deactivation.

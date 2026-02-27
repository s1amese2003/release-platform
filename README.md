# Release Platform (Go)

A lightweight artifact admission and approval release platform, implemented with Go + Gin + Asynq + PostgreSQL.

## Implemented Scope (MVP)

- Artifact upload (`zip/jar/war`) with SHA-256 hash
- Secure extraction with Zip Slip protection
- Latest SQL location:
  - `BOOT-INF/classes/upgrade/YYYYMMDD/upgrade.sql`
  - Auto-pick max date folder
- SQL risk scanning (HIGH/MEDIUM/LOW)
- `bootstrap-dev.yml` config diff against active online snapshot
- Policy decision (`PASS/WARN/BLOCK`)
- Approval actions and deployment gating
- Deployment simulation + active config snapshot update
- Async scan worker (Redis + Asynq)
- Audit log records
- Built-in lightweight web console (`/`) for release workflow

## Architecture

- HTTP API: Gin
- DB: PostgreSQL (GORM auto-migration)
- Queue: Redis + Asynq
- Storage: local filesystem (`./data`)

Processes:

- API server: `cmd/server`
- Worker: `cmd/worker`

## Quick Start

### 1. Prepare env

Copy env file:

```powershell
Copy-Item .env.example .env
```

### 2. Start infra (async mode)

```powershell
docker compose up -d
```

### 3. Install dependencies and run (async mode)

```powershell
go mod tidy
go run ./cmd/server
go run ./cmd/worker
```

If you want to run in separate terminals:

- Terminal A: `go run ./cmd/server`
- Terminal B: `go run ./cmd/worker`

### 4. Open web console

After server starts:

- Console: `http://127.0.0.1:8080/`
- Health: `http://127.0.0.1:8080/healthz`

### Local lightweight mode (no PostgreSQL/Redis)

If you only want quick local verification on one process:

```powershell
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_DSN='./data/release_platform.db'
$env:SCAN_MODE='sync'
go run ./cmd/server
```

In this mode, scanning runs synchronously and worker process is not required.

## API

### Create release ticket

```http
POST /api/releases
Content-Type: application/json

{
  "application": "order-service",
  "owner": "team-order",
  "environment": "dev",
  "version": "2026.02.26-01",
  "uploader": "alice"
}
```

### Upload artifact

```http
POST /api/releases/{id}/artifact
Content-Type: multipart/form-data
Form field: artifact=<file>
```

### Trigger scan

```http
POST /api/releases/{id}/scan
X-User: alice
```

### Query report

```http
GET /api/releases/{id}/report
```

### Approve release

```http
POST /api/releases/{id}/approve
X-Role: owner
Content-Type: application/json

{
  "approver": "bob",
  "action": "approve",
  "comment": "validated",
  "level": 2
}
```

### Deploy release

```http
POST /api/releases/{id}/deploy
X-Role: ops
Content-Type: application/json

{
  "operator": "ops-user"
}
```

### Policy APIs

```http
GET /api/policies
GET /api/policies?environment=dev
PUT /api/policies/{id}
```

## Status Model

- `scan_status`: `PENDING|RUNNING|DONE|FAILED`
- `approval_status`: `PENDING|APPROVED|REJECTED|NOT_REQUIRED`
- `admission_result`: `PASS|WARN|BLOCK`
- `deploy_status`: `NOT_TRIGGERED|RUNNING|SUCCESS|FAILED`

## SQL Risk Rules (Current)

- HIGH:
  - `DROP DATABASE`
  - `DROP TABLE`
  - `TRUNCATE TABLE`
  - `DELETE` without `WHERE`
  - `UPDATE` without `WHERE`
  - `ALTER TABLE ... DROP COLUMN`
- MEDIUM:
  - `RENAME TABLE`
  - index rebuild (`CREATE/DROP INDEX`)
- LOW:
  - create table without `IF NOT EXISTS`
  - drop table without `IF EXISTS`
  - missing rollback note in SQL file

## Config Diff Rules

- Flatten YAML to `key=value`
- Compare against active snapshot for same `application + environment`
- Diff type: `ADDED|DELETED|MODIFIED`
- Sensitive key masking: password/secret/token/etc.
- Critical keys include datasource/nacos and sensitive keys

## Project Structure

```text
cmd/
  server/main.go
  worker/main.go
internal/
  api/
  config/
  db/
  model/
  queue/
  service/
  util/
web/
  index.html
  assets/
```

## Limitations (Current MVP)

- SQL parser is rule-based (not full AST parser)
- Deployment is simulated callback point
- RBAC is lightweight header-based middleware (can be replaced by Casbin)
- Local filesystem storage only (can be switched to S3/MinIO)

## Next Recommended Enhancements

- Integrate Casbin RBAC and SSO
- Add AST SQL parser per database dialect
- Add CI pre-check API and webhook callback
- Add frontend (Vue3 + Element Plus) for report visualization
- Add export report (PDF/Excel)

# Release Platform（Go）

一个轻量的制品准入与审批发布平台，基于 Go + Gin + Asynq + PostgreSQL 实现。

## 已实现能力（MVP）

- 上传制品（`zip/jar/war`）并计算 SHA-256
- 安全解压（支持 Zip Slip 防护）
- 自动定位最新升级 SQL：
  - `BOOT-INF/classes/upgrade/YYYYMMDD/upgrade.sql`
  - 多日期目录时自动选最大日期
- SQL 风险扫描（`HIGH/MEDIUM/LOW`）
- `bootstrap-dev.yml` 与线上激活配置快照对比
- 策略引擎输出准入结论（`PASS/WARN/BLOCK`）
- 审批动作与部署门禁
- 部署模拟回调 + 发布后写入新基线快照
- 异步扫描 Worker（Redis + Asynq）
- 审计日志记录
- 内置轻量 Web 控制台（`/`）

## 技术架构

- HTTP API：Gin
- 数据库：PostgreSQL（GORM 自动迁移）
- 队列：Redis + Asynq
- 存储：本地文件系统（`./data`）

进程划分：

- API 服务：`cmd/server`
- Worker 服务：`cmd/worker`

## 快速启动

### 1. 准备环境变量

复制模板：

```powershell
Copy-Item .env.example .env
```

### 2. 启动依赖（异步模式）

```powershell
docker compose up -d
```

### 3. 安装依赖并运行（异步模式）

```powershell
go mod tidy
go run ./cmd/server
go run ./cmd/worker
```

如需分两个终端：

- 终端 A：`go run ./cmd/server`
- 终端 B：`go run ./cmd/worker`

### 4. 打开控制台

服务启动后访问：

- 控制台：`http://127.0.0.1:8080/`
- 健康检查：`http://127.0.0.1:8080/healthz`

### 本机轻量模式（不依赖 PostgreSQL/Redis）

如果你只想本机快速验证，可使用单进程模式：

```powershell
$env:DATABASE_DRIVER='sqlite'
$env:DATABASE_DSN='./data/release_platform.db'
$env:SCAN_MODE='sync'
go run ./cmd/server
```

说明：该模式下扫描为同步执行，不需要启动 Worker。

## 接口说明

### 创建发布单

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

### 上传制品

```http
POST /api/releases/{id}/artifact
Content-Type: multipart/form-data
Form field: artifact=<file>
```

### 触发扫描

```http
POST /api/releases/{id}/scan
X-User: alice
```

### 查询报告

```http
GET /api/releases/{id}/report
```

### 审批发布单

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

### 触发部署

```http
POST /api/releases/{id}/deploy
X-Role: ops
Content-Type: application/json

{
  "operator": "ops-user"
}
```

### 策略接口

```http
GET /api/policies
GET /api/policies?environment=dev
PUT /api/policies/{id}
```

## 状态字段

- `scan_status`：`PENDING|RUNNING|DONE|FAILED`
- `approval_status`：`PENDING|APPROVED|REJECTED|NOT_REQUIRED`
- `admission_result`：`PASS|WARN|BLOCK`
- `deploy_status`：`NOT_TRIGGERED|RUNNING|SUCCESS|FAILED`

## SQL 风险规则（当前）

- HIGH：
  - `DROP DATABASE`
  - `DROP TABLE`
  - `TRUNCATE TABLE`
  - 无 `WHERE` 的 `DELETE`
  - 无 `WHERE` 的 `UPDATE`
  - `ALTER TABLE ... DROP COLUMN`
- MEDIUM：
  - `RENAME TABLE`
  - 索引重建（`CREATE/DROP INDEX`）
- LOW：
  - 建表未使用 `IF NOT EXISTS`
  - 删表未使用 `IF EXISTS`
  - SQL 文件缺少回滚说明

## 配置差异规则

- YAML 扁平化为 `key=value`
- 与同 `application + environment` 的激活快照对比
- 差异类型：`ADDED|DELETED|MODIFIED`
- 敏感项自动脱敏（password/secret/token 等）
- 关键项包括 datasource/nacos 及敏感配置

## 项目目录

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

## 当前限制（MVP）

- SQL 解析为规则匹配，尚未引入完整 AST 解析
- 部署为模拟回调点，未接真实部署系统
- RBAC 为轻量 Header 校验，可替换为 Casbin
- 当前仅支持本地文件存储，可扩展到 S3/MinIO

## 后续建议

- 接入 Casbin + SSO
- 按数据库方言引入 AST SQL 解析
- 增加 CI 预检查 API 与 Webhook
- 增加报表导出（PDF/Excel）
- 完善白名单、过期策略与二级审批策略中心

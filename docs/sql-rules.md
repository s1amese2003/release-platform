# SQL 审计规则说明

## 风险等级

- `P0`：禁止，命中后自动驳回
- `P1`：警告，需要 DBA/开发负责人确认
- `P2`：提示，建议补充说明后继续

## 已实现规则

### P0

- `DROP TABLE` / `DROP DATABASE`
- `TRUNCATE TABLE`
- `DELETE` 无 `WHERE`
- `UPDATE` 无 `WHERE`
- 操作系统库：`mysql` / `information_schema` / `performance_schema`

### P1

- `ALTER TABLE DROP COLUMN`
- `ALTER TABLE MODIFY/CHANGE COLUMN`
- 大表加索引但缺少在线算法提示
- `GRANT` / `REVOKE`

### P2

- `ALTER TABLE ADD COLUMN` 无默认值
- `INSERT INTO db.table` 跨库写入
- `CREATE TABLE IF NOT EXISTS`

## 扩展规则方式

规则位于：`release-server/src/main/java/com/huatai/release/engine/sql/rule/`

新增规则步骤：

1. 实现 `SqlRule` 接口
2. 在 `SqlRiskAnalyzer` 中注册规则
3. 补充对应的测试样例
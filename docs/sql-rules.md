# SQL 瀹¤瑙勫垯璇存槑

## 椋庨櫓绛夌骇

- `P0`锛氱姝紝鍛戒腑鍚庤嚜鍔ㄩ┏鍥?- `P1`锛氳鍛婏紝闇€ DBA/寮€鍙戣礋璐ｄ汉纭
- `P2`锛氭彁绀猴紝寤鸿琛ュ厖璇存槑鍚庣户缁?
## 宸插疄鐜拌鍒?
### P0

- `DROP TABLE` / `DROP DATABASE`
- `TRUNCATE TABLE`
- `DELETE` 鏃?`WHERE`
- `UPDATE` 鏃?`WHERE`
- 鎿嶄綔绯荤粺搴擄紙`mysql`/`information_schema`/`performance_schema`锛?- 鐢熶骇鐜 `spring.profiles.active != prod`锛堥厤缃鍒欎骇鐢?P0锛?
### P1

- `ALTER TABLE DROP COLUMN`
- `ALTER TABLE MODIFY/CHANGE COLUMN`
- 鍔犵储寮曟湭鏄惧紡 `ALGORITHM=INPLACE`
- `GRANT` / `REVOKE`

### P2

- `ALTER TABLE ADD COLUMN` 鏃犻粯璁ゅ€?- `INSERT INTO db.table` 璺ㄥ簱鍐欏叆
- `CREATE TABLE IF NOT EXISTS`

## 鎵╁睍鏂瑰紡

瑙勫垯浣嶄簬 `release-server/src/main/java/com/huatai/release/engine/sql/rule/`銆?鏂板瑙勫垯姝ラ锛?
1. 瀹炵幇 `SqlRule` 鎺ュ彛銆?2. 鍦?`SqlRiskAnalyzer` 鐨?`rules` 鍒楄〃涓敞鍐屻€?3. 涓鸿鍒欒ˉ鍏呮祴璇曟牱渚嬶紙寤鸿鎸?P0/P1/P2 鍒嗙粍锛夈€?

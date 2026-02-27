# 配置基线说明

## 基线模型

基线按 `env + appName` 维度存储，值采用扁平化 `key=value` 结构。

环境通常包括：

- `dev`
- `sit`
- `uat`
- `prod`

当前实现默认使用内存存储（`BaselineMemoryStore`），后续可替换为 MySQL `env_baseline` 表。

## 比对维度

- 数据库：`spring.datasource.*`
- Redis：`spring.redis.*`
- Nacos：`nacos.*`
- 环境标识：`spring.profiles.active`
- 端口：`server.port`
- 调试/文档开关：`swagger.enabled`、`springdoc.api-docs.enabled`、`debug.*`

## 生产环境规则

- `spring.profiles.active` 必须为 `prod`（P0）
- 以下开关在生产环境不得为 `true`（P1）
- `swagger.enabled`
- `springdoc.api-docs.enabled`
- `debug.*`

## 前端基线页面

前端“环境基线”页面支持：

1. 按环境和应用加载基线
2. 使用 `key=value` 文本编辑
3. 一键保存回后端
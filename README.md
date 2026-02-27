# 项目上线平台

## 项目简介

本项目用于支撑 ITSM 上线发布管控，包含以下能力：

- 上传 JAR/ZIP 并自动解析升级内容
- SQL 审计（P0/P1/P2 风险分级）
- 配置比对（与环境基线逐项对比）
- 依赖扫描（版本变化与高风险提示）
- 审批流动作（APPROVE / REJECT / DEPLOY）

## 目录结构

- `release-server`：Spring Boot 3 后端
- `release-web`：Vue3 + Element Plus 前端
- `docs`：规则与基线文档

## 后端启动

```powershell
cd E:\release-platform\release-server
mvn --% spring-boot:run -DskipTests -Dspring-boot.run.arguments=--server.port=8081
```

## 前端启动

```powershell
cd E:\release-platform\release-web
npm install
npm run dev -- --host 0.0.0.0 --port 5173
```

## 访问地址

- 前端：`http://127.0.0.1:5173`
- 后端健康检查：`http://127.0.0.1:8081/api/system/health`

## 前端代理说明

前端默认将 `/api` 代理到 `http://127.0.0.1:8081`。

如果你本地后端端口不是 `8081`，可在启动前端前设置：

```powershell
$env:VITE_API_TARGET='http://127.0.0.1:8082'
npm run dev -- --host 0.0.0.0 --port 5173
```
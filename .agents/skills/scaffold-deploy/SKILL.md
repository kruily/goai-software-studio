---
name: scaffold-deploy
description: 按技术选型生成 deploy/ 下的部署物料。提供 nginx 配置模板、Dockerfile 模板和 docker-compose 模板。由 tech-lead / devops 在 bootstrap 或微服务拆分后调用。
---

# scaffold-deploy

你按当前项目的技术选型（module 数量、单体/微服务、数据库/队列/存储），
生成 `deploy/` 下的部署骨架。不预建实际文件到 deploy/，物料由你按需创建。

## 使用时机

- bootstrap 阶段 3 结束后，需要基础 deploy 骨架。
- 微服务拆分出新服务后，需要新增部署物料。
- 用户说"帮我搞一下部署配置""加 nginx 转发"。

## 前置

- 已确认架构选型（单体 vs 微服务、数据库、队列）。
- 已存在 `deploy/README.md`（约束说明，已有）。

## 执行流程

### 1. 读当前项目状态

- 从 `rules/tech-selection.md` 读：架构、数据库、队列。
- 从 `backend/` 下模块目录确认有几个模块。

### 2. 生成 nginx 配置（微服务 → references 中的模板）

单体：不需要 nginx，或简单转发到 api 端口。

微服务：为每个 `{module}/api` 生成转发规则。按 `references/nginx-template.md` 模板创建，nginx 配置块内容参考该文件：

每个 `{module}/api` 对应一个 `location` 块：
```
upstream {module}-api {
    server 127.0.0.1:{port};
}

server {
    listen 80;
    server_name api.example.com;

    location /api/v1/{module}/ {
        proxy_pass http://{module}-api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

常见端口约定：user=8881, order=8882, payment=8883, gateway/单体=8888。

### 3. 生成 Dockerfile（按 references/ 中的模板）

每个模块一个 Dockerfile。`scripts/gen-dockerfile.sh`（按需生成）。

模板内容（多阶段构建）：
```dockerfile
FROM golang:{go_version}-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG module=user
ARG type=api
RUN CGO_ENABLED=0 go build -o /server ./{module}/{type}

FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /server /server
COPY --from=builder /app/{module}/{type}/etc /etc
ENTRYPOINT ["/server", "-f", "/etc/{module}.yaml"]
```

### 4. 生成 docker-compose

按选型生成 `deploy/compose/docker-compose.yml`：
- 单体：backend api + 数据库 + Redis（如需要）。
- 微服务：各模块 api + rpc + mq + 数据库 + Redis + MinIO（如需要）。

### 5. 验证

- 生成的 nginx 配置语法：`nginx -t -c deploy/nginx/backend.conf`。
- 生成的 compose 格式：`docker compose -f deploy/compose/docker-compose.yml config`。

## 完成后

报告：生成了哪些 deploy 文件、nginx 转发规则、Dockerfile 的多阶段构建细节。

## 禁止

- 不要生成本地开发不需要的配置（如生产 SSL 证书、K8s → 除非显式要求）。
- 不覆盖用户已手动修改过的 deploy/ 文件。

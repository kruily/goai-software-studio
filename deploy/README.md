# deploy

运维部署物料集中目录，由 `devops` agent 维护。应用代码里不放部署逻辑。

```
deploy/
├── docker/     # 各服务 Dockerfile（多阶段构建，产物最小化）
├── compose/    # docker-compose.{dev,prod}.yml，编排 backend + 数据库/Redis/存储等基建
├── nginx/      # 反向代理转发：微服务下多个 api 服务由 nginx 转发（不用 go-zero gateway）
├── k8s/        # K8s manifests / helm（如需要）
├── ci/         # CI 流水线（构建、lint、发布）
└── env/        # 分环境 .env.example（不含真实密钥）
```

## 约定

- **不提交生产密钥/证书**；compose、`.env.example` 中凭证均为开发默认值。
- nginx 转发规则与后端各 api 服务路由前缀（`/api/v1/...`）保持一致。
- compose 编排的基建服务由 `rules/tech-selection.md` 的选型决定（数据库/队列/存储）。

## 由 bootstrap / devops 生成

具体 Dockerfile、compose、nginx conf 在 bootstrap 阶段或由 devops 按选型生成，此处仅为骨架。

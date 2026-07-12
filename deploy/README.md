# deploy

运维部署物料集中目录，由 `devops` agent 维护。应用代码里不放部署逻辑。

```
deploy/
├── docker/     # 各服务 Dockerfile（多阶段构建，产物最小化）
├── nginx/      # 反向代理转发：微服务下多个 api 服务由 nginx 转发（不用 go-zero gateway）
├── ci/         # CI 流水线（构建、lint、发布）
```

## 约定

- **不提交生产密钥/证书**；所有凭证均为开发默认值。
- nginx 转发规则与后端各 api 服务路由前缀（`/api/v1/...`）保持一致。
- 具体编排由 `rules/tech-selection.md` 的选型决定（数据库/队列/存储）。

## 由 bootstrap / devops 生成

具体 Dockerfile、compose、nginx conf 在 bootstrap 阶段或由 devops 按选型生成，此处仅为骨架。

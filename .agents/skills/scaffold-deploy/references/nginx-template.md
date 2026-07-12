# nginx 反向代理配置模板

## 微服务场景

多个 api 服务需要统一的入口转发。每个 `{module}/api` 一个 upstream。

```nginx
# deploy/nginx/backend.conf
# 由 scaffold-deploy 技能按模块列表生成

upstream user-api {
    server 127.0.0.1:8881;
}
# 新增模块在此追加：
# upstream {module}-api {
#     server 127.0.0.1:{port};
# }

server {
    listen 80;
    server_name api.example.com;

    # 跨域配置（如需）
    add_header Access-Control-Allow-Origin $http_origin;
    add_header Access-Control-Allow-Methods "GET, POST, OPTIONS";
    add_header Access-Control-Allow-Headers "Content-Type, Authorization";

    location /api/v1/user/ {
        proxy_pass http://user-api;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 新增模块在此追加：
    # location /api/v1/{module}/ {
    #     proxy_pass http://{module}-api;
    #     proxy_set_header Host $host;
    #     proxy_set_header X-Real-IP $remote_addr;
    # }

    # 健康检查端点
    location /health {
        return 200 "OK";
    }
}
```

## 单体场景

单体只有一个 api 端口，无需 nginx。
开发时直接访问 `http://localhost:8888/api/v1/{domain}/{action}`。
生产环境如需 nginx，按上述模板单个 upstream 即可。

## 端口约定

| 模块 | api | rpc | mq |
|------|-----|-----|----|
| user | 8881 | 9891 | — |
| order | 8882 | 9892 | — |
| payment | 8883 | 9893 | — |
| gateway/单体 | 8888 | — | — |

## 添加新模块

1. 在 `upstream` 块加一个 `{module}-api`。
2. 在 `location` 块加一个 `/api/v1/{module}/`。
3. `nginx -t` 验证。
4. `nginx -s reload`。

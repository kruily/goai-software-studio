# 环境变量示例（开发默认值，勿放生产密钥）。
# 复制到项目根 .env 使用；.env 已被 .gitignore 忽略。
# 由 scaffold-deploy 技能生成。

# --- 服务 ---
APP_ENV=dev
APP_HOST=0.0.0.0
APP_PORT=8888

# --- 数据库（按 tech-selection.md 选型填写；示例为 PostgreSQL）---
DB_DRIVER=postgres
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=app
DB_PARAMS=sslmode=disable

# --- Redis（缓存/锁/Asynq 后端）---
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

# --- 对象存储（示例为 MinIO）---
STORAGE_DRIVER=minio
STORAGE_ENDPOINT=127.0.0.1:9000
STORAGE_ACCESS_KEY_ID=minioadmin
STORAGE_ACCESS_SECRET=minioadmin
STORAGE_BUCKET=app
STORAGE_USE_SSL=false

# --- JWT ---
JWT_ACCESS_SECRET=dev-only-change-me
JWT_ACCESS_EXPIRE=86400

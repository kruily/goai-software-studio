# CI 流水线模板

项目克隆后，在 `.github/workflows/` 或 `.gitlab-ci.yml` 使用。

## GitHub Actions（推荐）

```yaml
# .github/workflows/ci.yml
# 由用户手动复制到仓库根。bootstrap 后按选型调整。
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  lint:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6

  unit-test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Run tests
        run: go test -race -count=1 -coverprofile=coverage.out ./...

  build:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: backend
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.26'
      - name: Build
        run: go build ./...
```

## GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - lint
  - test
  - build

lint:
  stage: lint
  image: golang:1.26
  script:
    - cd backend && go vet ./...

test:
  stage: test
  image: golang:1.26
  script:
    - cd backend && go test -race -count=1 ./...

build:
  stage: build
  image: golang:1.26
  script:
    - cd backend && go build ./...
```

## 集成测试 Job（可选，需 Docker）

```yaml
integration-test:
  runs-on: ubuntu-latest
  services:
    postgres:
      image: postgres:16
      env:
        POSTGRES_USER: test
        POSTGRES_PASSWORD: test
        POSTGRES_DB: test
    redis:
      image: redis:7
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: '1.26'
    - name: Integration tests
      run: cd backend && go test -tags=integration -count=1 ./...
      env:
        DB_HOST: postgres
        REDIS_ADDR: redis:6379
```

## 使用方式

- 复制对应文件到仓库根。
- 按项目名调整 `working-directory: backend`（若项目名不同）。
- 按选型调整 `go-version`（与 `go.mod` 一致）。

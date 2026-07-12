---
name: load-test
description: 对后端 API 执行负载/压力测试（k6），验证高并发下的延迟与错误率。从 OpenAPI spec 或 .api/.proto 定义生成 k6 脚本，按指定阈值（p95<500ms、错误率<1%）评估系统稳定性。用于上线前/大促前的容量评估。
---

# load-test

你用 k6 对后端 API 执行负载与压力测试。测试场景、阈值和虚拟用户数来自 spec 的非功能约束
（预估并发量、响应时间要求）。

## 前置

- k6 已安装：`brew install k6` 或 `docker run grafana/k6`。
- 后端服务运行在测试环境（或本地）。

## 使用时机

- 上线前需要评估容量。
- 大流量活动前排障。
- 非功能约束（`PROJECT.md`）要求验证性能指标。
- 用户说"做一下压测""看看能扛多少并发"。

## 执行流程

### 1. 读非功能约束与技术选型

- 从 `PROJECT.md` 读性能要求：预估并发、p95 响应时间、错误率上限。
- 从 `rules/tech-selection.md` 读部署规模。

### 2. 生成 k6 脚本

放在 `backend/loadtest/` 目录下：

```
backend/loadtest/
├── smoke.js        # 冒烟测试：1-2 VUs，验证基本功能
├── load.js         # 负载测试：目标并发持续 N 分钟
├── stress.js       # 压力测试：逐步增加并发直到系统瓶颈
└── spike.js        # 尖峰测试：短时间内爆发式并发
```

k6 脚本示例（`load.js`）：
```javascript
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '1m', target: 50 },   // 1 分钟内升到 50 并发
    { duration: '3m', target: 50 },   // 持续 3 分钟
    { duration: '1m', target: 0 },    // 1 分钟内降到 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% 请求在 500ms 内
    http_req_failed: ['rate<0.01'],   // 错误率低于 1%
  },
};

export default function () {
  const res = http.post('http://localhost:8888/api/v1/front/user/getProfile');
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(1);
}
```

### 3. 执行与收集

```bash
k6 run backend/loadtest/load.js
```

- 收集输出：请求数、p50/p95/p99 延迟、错误率、每秒请求数。
- 若阈值溢出，标记性能风险并报告给 tech-lead。

### 4. 报告

- 测试参数：VU 数、持续时间、阈值。
- 结果摘要：通过/未达标的阈值、最大并发、瓶颈推测（CPU/DB/网络）。
- 附 k6 HTML 报告（`k6 run --out json=result.json`）。

## 完成

汇报负载测试结果、阈值达成情况、性能建议。

## 禁止

- 不在生产环境执行负载测试（必须用隔离的测试环境）。
- 不设无意义的高阈值（应与 PROJECT.md 的非功能约束一致）。
- 不忽略预热阶段（冷启动的延迟不反映真实性能）。

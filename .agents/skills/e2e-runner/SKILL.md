---
name: e2e-runner
description: 运行端到端浏览器测试（Playwright），覆盖前端关键用户流程。从 spec 的用户故事生成 E2E 脚本，截图验证 UI 状态，报告失败。用于部署后/PR 合并前的全链路验证。
---

# e2e-runner

你用 Playwright 运行浏览器端到端测试，覆盖前端关键用户流程。测试场景源自 OpenSpec 的 spec
和 PROJECT.md 的核心功能清单。

## 前置

- 前端项目已初始化，可本地启动（`npm run dev` 或其他）。
- Playwright 已安装：`npm init playwright@latest` 或 `npx playwright install`。
- 后端服务本地运行或可用测试环境。

## 使用时机

- 前端功能开发完成（frontend-dev 交付后），需要验证完整用户流程。
- PR 合并前需要全链路回归。
- 用户说"跑一下 E2E""看看登录流程是否正常"。

## 执行流程

### 1. 从 spec 提取用户故事

- 读 `openspec/changes/<name>/specs/` 的用户故事/场景。
- 每个场景对应一个 Playwright test case：
  - `GIVEN 用户已登录 WHEN 访问个人资料 THEN 显示昵称和头像`

### 2. 生成 Playwright 测试

```
{front-end}/e2e/
├── login.spec.ts       # 登录/注册流程
├── profile.spec.ts     # 用户资料
└── {feature}.spec.ts   # 按 feature 分文件
```

Playwright 测试骨架：
```typescript
import { test, expect } from '@playwright/test'

test('用户登录后能看到个人资料', async ({ page }) => {
  await page.goto('/login')
  await page.fill('[data-testid="email"]', 'user@example.com')
  await page.fill('[data-testid="password"]', 'password123')
  await page.click('[data-testid="login-btn"]')
  await expect(page.locator('[data-testid="user-name"]')).toBeVisible()
})
```

### 3. 数据隔离与回滚

- 测试数据通过 API 创建，测试完成后调用清理接口删除。
- 使用独立的测试数据库或事务回滚。
- 使用 `data-testid` 属性定位元素，禁止依赖 CSS class（防选择器脆弱）。

### 4. 执行与截图

```bash
cd frontend  # 或 admin-web/mobile
npx playwright test --reporter=list
```

- 失败时自动截图和录屏，附到 bug 报告中。
- 每次 E2E 运行前运行后端集测，确认后端 API 正常（否则 E2E 失败不一定是前端问题）。

### 5. 报告

- 输出：通过/失败数、截图 URL、失败堆栈。
- flaky 判断：一次失败 + 一次重试通过 = flaky（标记不报 bug）；两次失败 = 真实 bug。

## 完成

报告 E2E 结果、截图、失败详情。

## 禁止

- 不依赖 CSS class 定位元素（用 `data-testid`）。
- 不测试非关键路径（E2E 预算有限，只覆盖核心用户流程）。
- 不把生产凭证/密钥传给 E2E runner。

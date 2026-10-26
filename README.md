# 布迹

布迹是一款面向手作爱好者的本地优先布料手帐。当前仓库正在从原生 HTML/CSS/JavaScript 原型迁移到 uni-app + Go API + PostgreSQL 架构，并使用 Vercel Services 部署。

## 仓库结构

```text
.
├── app.js / index.html / styles.css  # 现有交互原型，迁移期间保持不变
├── apps/uniapp                       # uni-app Vue 3 + TypeScript 客户端
├── services/api                      # Go HTTP API，Vercel container runtime
├── database/migrations               # PostgreSQL 数据结构
└── deploy/cloudbase                  # 暂缓使用的 CloudBase 部署参考
```

## 本地开发

### 客户端

```bash
cd apps/uniapp
npm install
npm run dev:h5
```

微信小程序：

```bash
npm run dev:mp-weixin
```

然后使用微信开发者工具打开 `apps/uniapp/dist/dev/mp-weixin`。

### API

```bash
cd services/api
go test ./...
go run ./cmd/api
```

API 会读取仓库根目录 `.env`，连接 PostgreSQL 并自动执行尚未应用的迁移。默认监听 `http://localhost:8080`，健康检查为 `GET /healthz`。

## 环境变量

复制根目录 `.env.example`，不要提交真实密钥。支持 `DATABASE_URL`、Supabase 提供的 `POSTGRES_URL` 或 `POSTGRES_URL_NON_POOLING`。飞书 App Secret、微信 App Secret 和数据库凭证必须只存在于服务端或部署平台环境变量中。

## 当前迁移边界

- 根目录原型仍是现有完整交互和视觉基准。
- uni-app 已建立相同视觉语言的工作台与五个页面入口。
- Go API、PostgreSQL schema 和 Vercel 多服务部署已经建立。
- 登录、图片上传、正式数据库连接和飞书同步将在后续迭代接入。

## Vercel Services

根目录 `vercel.json` 将 H5 和 Go API 部署在同一个 Vercel 项目：

- `/healthz`、`/v1/*` → Go API 容器
- 其他请求 → uni-app H5

在 Vercel 项目设置中把 Framework Preset 设为 `Services`。生产环境同域请求无需设置 `VITE_API_BASE_URL`；本地开发仍默认访问 `http://localhost:8080`。


# 布迹

布迹是一款面向手作爱好者的本地优先布料手帐。当前仓库正在从原生 HTML/CSS/JavaScript 原型迁移到 uni-app + CloudBase 架构。

## 仓库结构

```text
.
├── app.js / index.html / styles.css  # 现有交互原型，迁移期间保持不变
├── apps/uniapp                       # uni-app Vue 3 + TypeScript 客户端
├── services/api                      # Go HTTP API，面向 CloudBase 云托管
├── database/migrations               # PostgreSQL 数据结构
└── deploy/cloudbase                  # CloudBase 部署说明
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

默认监听 `http://localhost:8080`，健康检查为 `GET /healthz`。

## 环境变量

复制根目录 `.env.example`，不要提交真实密钥。飞书 App Secret、微信 App Secret 和数据库凭证必须只存在于服务端或 CloudBase 环境变量中。

## 当前迁移边界

- 根目录原型仍是现有完整交互和视觉基准。
- uni-app 已建立相同视觉语言的工作台与五个页面入口。
- Go API、PostgreSQL schema 和 CloudBase 容器部署骨架已经建立。
- 登录、图片上传、正式数据库连接和飞书同步将在后续迭代接入。


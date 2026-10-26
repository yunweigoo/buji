# CloudBase 部署准备

## 目标资源

1. 一个 CloudBase 环境，建议开发环境与生产环境分开。
2. CloudBase PostgreSQL。
3. CloudBase 云存储，默认私有读。
4. CloudBase 云托管服务 `buji-api`。
5. 后续用于微信登录和定时同步的云函数。

## 云托管服务

构建目录选择 `services/api`，Dockerfile 位于该目录。服务监听平台注入的 `PORT`，默认值为 `8080`。

至少配置：

```text
APP_ENV=production
DATABASE_URL=...
CLOUDBASE_ENV_ID=...
CLOUDBASE_STORAGE_BUCKET=...
```

接入飞书后再添加：

```text
FEISHU_APP_ID=...
FEISHU_APP_SECRET=...
FEISHU_BASE_APP_TOKEN=...
```

不要将真实密钥写入 Git、Dockerfile 或客户端 `manifest.json`。

## 健康检查

容器部署后配置：

```text
GET /healthz
```

成功响应示例：

```json
{
  "data": {
    "status": "ok",
    "service": "buji-api"
  },
  "requestId": "req-..."
}
```

## 暂未自动创建资源的原因

资源创建需要明确 CloudBase 环境 ID、地域、计费方案和微信小程序 AppID。当前代码已经具备部署边界；在这些参数确认前，不自动创建可能产生费用的云资源。


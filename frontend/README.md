# ExchangeApp 前端

Vue 3、TypeScript、Vue Router 和 Vite。支持汇率记录与换算、登录注册、新增汇率、文章发布与阅读，以及点赞。

## 本地运行

需要 Node.js 20.19+ 或 22.12+。先配置 MySQL 和 Redis，再在同级的 `backend` 目录运行 `go run .` 启动接口。然后在 `frontend` 目录执行：

```sh
npm ci
npm run dev
```

打开 Vite 打印的本地地址，默认端口为 5173。Vite 将 `/api` 请求代理到 `http://127.0.0.1:3000`。如后端地址不同，把 `.env.example` 复制为 `.env.local` 并修改 `API_PROXY_TARGET`。

## 构建

```sh
npm run build
npm run preview
```

`build` 会先检查 TypeScript，再输出到 `dist/`。部署时优先由同源网关转发 `/api/*` 到后端。也可在构建前设置 `VITE_API_BASE_URL` 为包含 `/api` 的后端地址，并确保后端允许前端来源。静态服务器需将前端页面路径回退到 `index.html`。

## 接口行为

- 汇率列表可匿名访问。新增汇率和文章功能需要登录，登录状态保存在当前浏览器标签页。
- 换算优先使用该方向最新的有效汇率。直接汇率缺失时，使用反向记录的倒数。相同币种按 1:1 换算。
- 当前点赞接口每次成功请求都会加一。前端重新获取服务端计数，并在请求未完成时禁用按钮。
- 发布文章遇到网络错误或服务端错误时，界面会提示先查看发布结果。数据库写入成功后，缓存清理仍有失败的可能。

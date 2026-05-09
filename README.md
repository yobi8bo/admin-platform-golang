# admin-platform

企业级后台管理平台骨架。

## 技术栈

- 后端：Go、Gin、GORM、PostgreSQL、Redis、RustFS、Viper、Zap
- 前端：Vue 3、JavaScript、Vite、Vue Router、Pinia、Axios、Ant Design Vue、SCSS、vue-i18n

## 本地启动

启动依赖：

```bash
docker compose up -d
```

启动后端：

```bash
cd backend
go mod tidy
go run ./cmd/server
```

启动前端：

```bash
cd frontend
npm install
npm run dev
```

默认访问地址：

- 前端：http://localhost:5173
- 后端：http://localhost:8085/healthz
- RustFS Console：http://localhost:9001

默认账号：

- 用户名：`admin`
- 密码：`Admin@123`

## 权限模型

- 菜单权限：后端 `sys_menus` 返回菜单树，前端动态生成路由。
- 按钮权限：前端通过 `v-permission` 控制显示。
- API 权限：后端中间件根据角色菜单权限强制校验。
- 数据权限：角色模型已预留 `data_scope`，业务查询可按部门/本人范围扩展。

## 目录

```text
backend/
  cmd/server
  internal/bootstrap
  internal/middleware
  internal/modules
  internal/pkg
  migrations
frontend/
  src/api
  src/router
  src/stores
  src/layouts
  src/views
```

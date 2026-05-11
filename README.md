# admin-platform

企业级中后台管理平台，面向权限、组织、文件、审计等通用管理场景。项目采用前后端分离架构。

## 项目定位

本项目不是单页演示模板，而是一套可继续扩展的企业后台基础工程：

- 统一登录、JWT 鉴权、刷新令牌与退出登录。
- 用户、角色、菜单、部门等系统管理能力。
- 菜单权限、按钮权限、API 权限、数据权限预留。
- 文件上传、临时下载链接、删除权限控制。
- 登录日志、操作日志与审计追踪。
- 企业运营后台 UI，采用固定左侧导航、顶部信息栏、紧凑筛选区、专业表格和指标工作台。

## 技术栈

### 后端

- Go
- Gin
- GORM
- PostgreSQL
- Redis
- RustFS / S3 兼容对象存储
- Viper
- Zap
- JWT

### 前端

- Vue 3
- JavaScript
- Vite
- Vue Router
- Pinia
- Axios
- Ant Design Vue
- SCSS
- vue-i18n

## 功能模块

| 模块 | 说明 |
| --- | --- |
| 认证中心 | 登录、登出、当前用户信息、权限拉取、JWT 鉴权 |
| 工作台 | 运营指标、访问趋势、模块状态、待处理事项、权限变更概览 |
| 用户管理 | 用户筛选、分页表格、角色授权、账号状态、批量操作预留 |
| 角色管理 | 角色编码、数据范围、菜单权限配置 |
| 菜单管理 | 菜单树、按钮权限、前端导航来源 |
| 部门管理 | 组织树与部门数据维护 |
| 文件管理 | 文件上传、下载链接、删除权限 |
| 审计日志 | 登录日志、操作日志、API 操作记录 |
| 个人中心 | 个人资料、头像上传与账号信息维护 |

## 本地启动

### 1. 准备环境

建议版本：

- Go 1.23+
- Node.js 20+
- Docker / Docker Compose

复制环境变量文件：

```bash
cp .env.example .env
```

### 2. 启动基础依赖

```bash
docker compose up -d
```

依赖服务：

- PostgreSQL：`localhost:5432`
- Redis：`localhost:8379`
- RustFS API：`localhost:9000`
- RustFS Console：`http://localhost:9001`

数据库初始化 SQL 位于 `backend/migrations`，PostgreSQL 容器首次启动时会自动执行。

### 3. 启动后端

```bash
cd backend
go mod tidy
go run ./cmd/server
```

默认配置文件：`backend/configs/config.yaml`

后端地址：

- Health Check：`http://localhost:8085/healthz`
- API 前缀：`http://localhost:8085/api`

### 4. 启动前端

```bash
cd frontend
npm install
npm run dev
```

默认访问地址：

- 前端：`http://localhost:5173`

Vite 已配置 `/api` 和 `/healthz` 代理到 `http://localhost:8085`。

## 默认账号

```text
用户名：admin
密码：Admin@123
```

## 目录结构

```text
.
├── backend
│   ├── cmd/server              # 服务入口
│   ├── configs                 # 后端配置
│   ├── internal
│   │   ├── bootstrap           # 依赖初始化、路由注册
│   │   ├── config              # 配置加载
│   │   ├── middleware          # 鉴权、权限、审计、日志中间件
│   │   ├── modules
│   │   │   ├── audit           # 审计日志
│   │   │   ├── auth            # 认证授权
│   │   │   ├── file            # 文件服务
│   │   │   └── system          # 用户、角色、菜单、部门
│   │   └── pkg                 # 通用工具包
│   └── migrations              # 数据库初始化与增量脚本
├── frontend
│   ├── src
│   │   ├── api                 # Axios 请求层与业务 API
│   │   ├── assets/styles       # 全局企业后台样式
│   │   ├── components          # 通用业务组件
│   │   ├── hooks               # 组合式业务逻辑
│   │   ├── layouts             # 应用布局
│   │   ├── locales             # 国际化资源
│   │   ├── router              # 路由与权限守卫
│   │   ├── stores              # Pinia 状态管理
│   │   ├── utils               # 权限指令等工具
│   │   └── views               # 业务页面
│   └── vite.config.js
├── docker-compose.yml
└── .env.example
```

## 常用命令

后端：

```bash
cd backend
go test ./...
go run ./cmd/server
```

前端：

```bash
cd frontend
npm install
npm run dev
npm run build
```

依赖服务：

```bash
docker compose up -d
docker compose ps
docker compose logs -f
```

## 配置说明

后端默认读取 `backend/configs/config.yaml`，主要配置项：

| 配置 | 说明 |
| --- | --- |
| `server.addr` | 后端监听地址，默认 `:8085` |
| `server.allowedOrigins` | CORS 允许来源 |
| `database.dsn` | PostgreSQL 连接串 |
| `redis.addr` | Redis 地址 |
| `rustfs.endpoint` | RustFS / S3 地址 |
| `rustfs.bucket` | 文件桶名称 |
| `jwt.secret` | JWT 签名密钥 |
| `jwt.accessTTLMinutes` | Access Token 有效期 |
| `jwt.refreshTTLDays` | Refresh Token 有效期 |

生产环境必须修改 `JWT_SECRET`、数据库密码、RustFS 密钥，并按实际域名配置 CORS。


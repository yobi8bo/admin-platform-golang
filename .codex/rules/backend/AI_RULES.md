# Backend AI Coding Rules

本文档用于约束 AI 在 `backend` 目录生成或修改 Go 后端代码时的边界、风格和架构。所有后端改动必须优先遵循现有代码，其次遵循本文档。

## 1. 项目技术栈与边界

- 后端是 Go 1.22 项目，模块名为 `admin-platform/backend`。
- Web 框架使用 Gin，ORM 使用 GORM，数据库为 PostgreSQL。
- 缓存/会话状态使用 Redis，对象存储使用 RustFS，日志使用 zap。
- 不要引入新的 Web 框架、ORM、配置框架、日志框架或依赖注入框架，除非用户明确要求。
- 所有业务代码必须位于 `backend/internal` 下；入口只放在 `backend/cmd/server`。
- 不要修改前端、部署、根目录文件，除非需求明确涉及。

## 2. 目录与架构约定

现有结构是主要边界：

```text
cmd/server/                 程序入口
internal/bootstrap/         应用初始化、依赖装配、路由注册
internal/config/            配置加载与配置结构
internal/middleware/        Gin 中间件
internal/modules/<module>/  业务模块
internal/pkg/               可复用基础工具
migrations/                 SQL 迁移
configs/                    配置文件
```

新增业务能力优先放在 `internal/modules/<module>` 中。一个模块至少包含：

- `model.go`：GORM 模型、表名、模块内 DTO 类型。
- `handler.go`：Gin handler、路由注册、请求解析、响应。
- `*_test.go`：核心行为、边界条件或 bugfix 的测试。

只有跨模块复用且稳定的能力才能放入 `internal/pkg`。不要为了单个模块提前抽象公共包。

## 3. 路由与 Handler 规则

- 每个模块暴露 `NewHandler(...) *Handler` 和 `Register(...)` 方法。
- 路由注册集中在模块自己的 `Register` 方法中，再由 `internal/bootstrap/app.go` 装配。
- 公开接口通过模块自己的 `RegisterPublic` 或在 bootstrap 中明确注册；默认接口应走私有路由。
- 私有接口必须经过 `middleware.Auth(...)`。
- 需要权限控制的接口必须使用 `require("module:resource:action")`。
- 不要在 handler 中直接 `c.JSON`，统一使用 `response.OK`、`response.Created`、`response.Fail`。
- Handler 方法负责 HTTP 层工作：绑定参数、读取上下文、调用数据库/服务、返回响应。不要在 handler 中堆放大量无关工具逻辑，复杂逻辑拆成同文件私有函数或模块内 service/helper。

示例风格：

```go
func (h *Handler) Register(rg *gin.RouterGroup, require func(string) gin.HandlerFunc) {
	group := rg.Group("/example")
	group.GET("", require("example:list"), h.List)
	group.POST("", require("example:create"), h.Create)
}
```

## 4. 请求、响应与错误

- 请求结构体使用小写私有类型，如 `createUserReq`、`updateProfileReq`。
- JSON 字段使用 lower camelCase，例如 `roleIds`、`pageSize`、`createdAt`。
- 必填字段使用 Gin binding tag，例如 `binding:"required"`。
- 成功响应统一：
  - 查询/更新：`response.OK(c, data)`
  - 创建：`response.Created(c, data)`
- 分页响应统一使用 `response.Page[T]`，字段为 `list`、`total`、`page`、`pageSize`。
- 错误响应统一使用 `response.Fail(c, httpStatus, errs.CodeXxx, message)`。
- 不要把内部敏感信息、密码哈希、token secret、数据库连接串返回给前端。
- 对用户可理解的业务错误，返回明确中文消息；对认证/系统通用错误可沿用已有英文消息。

常用错误码：

- 参数错误：`http.StatusBadRequest` + `errs.CodeBadRequest`
- 未登录/令牌错误：`http.StatusUnauthorized` + `errs.CodeUnauthorized`
- 无权限：`http.StatusForbidden` + `errs.CodeForbidden`
- 资源不存在：`http.StatusNotFound` + `errs.CodeNotFound`
- 系统错误：`http.StatusInternalServerError` + `errs.CodeInternal`

## 5. 数据库与模型

- 模型使用 GORM struct tag，表名通过 `TableName()` 显式声明。
- 系统业务表沿用 `sys_` 前缀。
- 通用主键使用 `uint`，时间字段优先使用 `timex.DateTime`。
- 带创建/更新/删除审计字段的模型优先复用 `system.BaseModel`，除非模块已有自己的模型约定。
- 列名、关联表名、权限数据必须与迁移 SQL 保持一致。
- 查询必须处理 `.Error`，不要忽略数据库错误。
- 多表或多步骤写入必须使用 `h.db.Transaction(...)`，尤其是创建主表并写关联表时。
- 更新时显式控制允许更新的字段，避免直接把请求结构体整包写入数据库。
- 删除默认使用 GORM 删除；是否软删除取决于模型是否包含 `gorm.DeletedAt`。
- 新增表、字段、初始权限、初始菜单时必须新增 migration，不要只依赖启动时自动修补。

## 6. 权限、认证与审计

- 用户身份从 `contextx.UserID(c)` 获取，不要从请求体或 query 中信任 userId。
- 用户角色从 `contextx.RoleIDs(c)` 获取。
- 新增敏感接口必须配置权限字符串，并补充对应菜单/权限迁移。
- 修改数据的接口必须走私有路由，使 `middleware.OperationAudit` 能记录操作日志。
- 登录、登出、密码、token 相关逻辑必须沿用 `internal/modules/auth` 与 `internal/pkg/jwt` 的模式。
- 密码只能使用 `internal/pkg/crypto` 中的哈希与校验函数，不要明文保存或比较。

## 7. 配置、日志与上下文

- 配置字段集中在 `internal/config`，不要在业务代码中硬编码可变配置。
- 日志使用 zap，通过已有 logger 初始化链路接入，不要使用 `fmt.Println` 或标准库 `log` 记录业务日志。
- 外部调用、Redis、对象存储应尽量使用 `c.Request.Context()` 或调用方传入的 `context.Context`。
- 不要在 handler 中创建新的全局客户端。依赖应从 `bootstrap.New` 初始化后注入 Handler。

## 8. 文件与对象存储

- 文件上传、下载 URL、删除逻辑优先复用 `internal/modules/file` 的模式。
- 对上传文件名必须使用 `filepath.Base` 等方式处理，避免路径穿越。
- 头像或图片类接口必须校验 `Content-Type`。
- 预签名 URL 必须设置合理过期时间，不要返回长期公开链接。

## 9. 代码风格

- 所有 Go 代码必须通过 `gofmt`。
- import 使用标准 Go 分组：标准库、项目内包、第三方包。
- 命名遵循 Go 习惯：导出类型用 PascalCase，私有函数/类型用 camelCase。
- 小函数优先，避免一个 handler 承担过多不相关逻辑。
- 不要做无关重构、无关格式化或大范围移动文件。
- 不要引入泛型、反射、复杂抽象来解决简单 CRUD。
- 不要吞掉关键错误；如果确实可以忽略，必须能从上下文看出无副作用，例如审计日志写入失败。

## 10. 注释规则

AI 生成或修改后端代码时，必须补充能帮助维护者理解业务意图的注释。注释应解释“为什么这样做”或“这里承担什么边界”，不要复述代码字面行为。

- 新增导出类型、导出函数、导出常量必须添加符合 Go 规范的文档注释，以标识符名称开头。
- 新增 Handler、Service、Repository、模块初始化方法时，应在类型或关键方法上说明模块职责、权限边界或调用场景。
- 涉及权限校验、认证、审计、事务、幂等、并发、数据一致性、文件安全、外部服务调用的代码，必须在关键逻辑前添加简短注释说明约束和原因。
- 新增请求/响应 DTO 或 GORM 模型时，如字段业务含义不直观、存在枚举值、状态流转、单位、默认值或安全限制，必须添加字段注释。
- 新增 SQL migration 时，文件头或关键语句附近应说明迁移目的；新增初始权限、菜单、种子数据时必须说明其对应的接口或业务能力。
- 对复杂查询、跨表关联、批量更新、软删除/硬删除选择、预签名 URL、上传文件校验等容易误改的逻辑，必须添加维护性注释。
- 修改既有无注释代码时，如果改动触及核心业务、权限、数据模型或非显而易见逻辑，应顺手补齐相关注释，但不要为无关代码做大范围补注释。
- 注释语言优先使用中文，除非文件中已有统一英文注释风格。
- 禁止添加无价值注释，例如“获取参数”“调用方法”“返回结果”这类仅复述下一行代码的注释。
- 禁止用注释掩盖复杂代码；如果逻辑难以通过简短注释说清，应优先拆分函数或命名清晰的局部变量。

## 11. 测试与验证

- 修改公共逻辑、中间件、权限、鉴权、响应格式、分页、审计、文件安全逻辑时必须补测试。
- 新增模块至少覆盖：
  - 参数非法时返回 400。
  - 无权限或未登录时返回对应错误。
  - 正常创建/查询/更新/删除路径。
  - 关键边界条件，例如空列表、无关联数据、重复数据。
- 优先使用 Go 标准测试和项目已有测试风格。
- 提交前至少运行：

```bash
go test ./...
```

如只改了文档或配置，可以说明未运行测试的原因。

## 12. AI 生成代码的禁止事项

AI 不得：

- 绕过 `response` 直接拼响应格式。
- 绕过 `errs` 自定义一套错误码。
- 绕过 `middleware.Auth` 或 `RequirePermission` 暴露敏感接口。
- 在业务代码中硬编码 JWT secret、数据库 DSN、Redis 密码、对象存储密钥。
- 新增未经用户确认的大型依赖或框架。
- 用字符串拼接用户输入构造 SQL；必须使用 GORM 参数绑定或安全 API。
- 在未确认需求时改变已有 API 路径、JSON 字段名、响应结构或数据库表结构。
- 删除或重写现有迁移文件；新增变更应追加新的 migration。
- 返回密码字段、secret、refresh token 存储 key 等敏感信息。
- 为了“更整洁”做跨模块大重构。

## 13. AI 输出要求

AI 在生成后端代码前应先说明：

- 会修改哪些文件。
- 是否新增接口、权限、迁移或配置项。
- 是否需要测试。
- 是否会新增或补充注释，以及注释覆盖的关键业务点。

AI 完成后必须说明：

- 实际修改了哪些文件。
- 新增或变更了哪些接口/权限/数据表。
- 运行了哪些验证命令；如果没运行，说明原因。
- 是否存在兼容性或迁移注意事项。
- 新增或补充了哪些关键注释；如果没有补充注释，说明原因。

## 14. 推荐的 AI 指令模板

后续让 AI 写后端代码时，可以附上这段：

```text
请严格遵守 backend/AI_RULES.md。
本项目后端使用 Go + Gin + GORM + PostgreSQL。
新增代码应符合现有 internal/modules、internal/pkg、response、errs、middleware、migration 的约定。
生成或修改后端代码时必须添加有维护价值的注释，尤其是导出标识符、权限、事务、模型字段和复杂业务逻辑。
不要做无关重构，不要改动既有 API 行为，除非我明确要求。
完成后请运行 gofmt 和 go test ./...，并说明修改文件、接口变化、迁移变化和测试结果。
```

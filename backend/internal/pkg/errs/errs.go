package errs

const (
	// CodeUnauthorized 表示未登录、令牌缺失或令牌无效。
	CodeUnauthorized = 40100
	// CodeForbidden 表示已登录但缺少访问当前资源的权限。
	CodeForbidden = 40300
	// CodeBadRequest 表示请求参数或业务前置条件不合法。
	CodeBadRequest = 40000
	// CodeNotFound 表示请求的业务资源不存在。
	CodeNotFound = 40400
	// CodeInternal 表示服务端内部错误。
	CodeInternal = 50000
)

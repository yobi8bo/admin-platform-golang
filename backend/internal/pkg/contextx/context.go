package contextx

import "github.com/gin-gonic/gin"

const (
	// UserIDKey 是认证中间件写入 Gin 上下文的用户 ID 键。
	UserIDKey = "userId"
	// RoleIDsKey 是认证中间件写入 Gin 上下文的角色 ID 列表键。
	RoleIDsKey = "roleIds"
)

// UserID 从 Gin 上下文读取当前用户 ID，未登录或类型不匹配时返回 0。
func UserID(c *gin.Context) uint {
	v, _ := c.Get(UserIDKey)
	id, _ := v.(uint)
	return id
}

// RoleIDs 从 Gin 上下文读取当前用户角色 ID 列表，缺失时返回 nil。
func RoleIDs(c *gin.Context) []uint {
	v, _ := c.Get(RoleIDsKey)
	ids, _ := v.([]uint)
	return ids
}

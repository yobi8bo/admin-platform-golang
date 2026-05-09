package contextx

import "github.com/gin-gonic/gin"

const (
	UserIDKey  = "userId"
	RoleIDsKey = "roleIds"
)

func UserID(c *gin.Context) uint {
	v, _ := c.Get(UserIDKey)
	id, _ := v.(uint)
	return id
}

func RoleIDs(c *gin.Context) []uint {
	v, _ := c.Get(RoleIDsKey)
	ids, _ := v.([]uint)
	return ids
}

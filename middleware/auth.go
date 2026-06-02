package middleware

import (
	"GinChat/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JwtAuth(ws ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		isWs := false
		if len(ws) > 0 {
			isWs = ws[0]
		}
		var auth string
		if isWs {
			auth = c.Param("token")
		} else {
			auth = c.GetHeader("Authorization")
		}
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			utils.Fail(c, 401, "请先登录")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		// 解析 token
		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "登录已过期或无效")
			c.Abort()
			return
		}

		// 把用户ID存到 gin context 里
		c.Set("userId", claims.UserId)
		c.Next()
	}
}

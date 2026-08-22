// 创建时间：2026/8/21 下午6:06
package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"learn_layer_demo/utils"
	"net/http"
)

// JWTAuthMiddleware JWT鉴权中间件
func JWTAuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.Fail(c, http.StatusUnauthorized, 40100, "请求头缺少Authorization")
		c.Abort()
		return
	}

	var tokenStr string
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
	if err != nil || tokenStr == "" {
		utils.Fail(c, http.StatusUnauthorized, 40101, "Authorization格式错误，应为Bearer token")
		c.Abort()
		return
	}

	userId, err := utils.ParseToken(tokenStr)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, 40102, "token非法或已过期")
		c.Abort()
		return
	}
	// 将userId存入上下文
	c.Set("userId", userId)
	// 放行
	c.Next()
}

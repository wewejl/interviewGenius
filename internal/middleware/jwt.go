package middleware

import (
	"interviewGenius/internal/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWT 认证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string
		// 从 Authorization 头获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果没有token
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := util.ParseToken(token)
		if err != nil {
			zap.L().Error("解析令牌失败", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "无效的认证令牌",
			})
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// CheckPermission 权限检查中间件
func CheckPermission(method, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
				"msg":  "未认证",
			})
			c.Abort()
			return
		}

		// 这里应该有权限检查逻辑
		// 暂不实现详细的权限检查
		hasPermission := true

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "没有权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuth 管理员权限验证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code": http.StatusForbidden,
				"msg":  "需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

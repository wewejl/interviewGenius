package v1

import (
	"interviewGenius/internal/model"
	"interviewGenius/internal/pkg/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// VerifyPermissionRequest 验证权限请求
type VerifyPermissionRequest struct {
	Method     string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	Path       string `json:"path" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}

// VerifyPermission 验证用户是否有权限
func VerifyPermission(c *gin.Context) {
	var req VerifyPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取当前用户ID
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 检查用户是否有权限
	hasPermission, err := model.CheckUserPermission(userID, req.Method, req.Path)
	if err != nil {
		zap.L().Error("检查权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "检查权限失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "验证权限成功",
		"data": gin.H{
			"has_permission": hasPermission,
		},
	})
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 验证刷新令牌
	token := strings.TrimPrefix(req.RefreshToken, "Bearer ")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "无效的刷新令牌",
			"data": nil,
		})
		return
	}

	// 获取用户ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 检查用户是否存在
	user, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "用户不存在",
			"data": nil,
		})
		return
	}

	// 生成新的访问令牌
	newToken, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		zap.L().Error("生成令牌失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成令牌失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "刷新令牌成功",
		"data": gin.H{
			"token": newToken,
		},
	})
}

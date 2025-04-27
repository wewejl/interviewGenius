package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"interviewGenius/internal/model"
	"interviewGenius/pkg/util"
	"net/http"
)

// AuthVerifyRequest 权限验证请求
type AuthVerifyRequest struct {
	Method string `json:"method" binding:"required"`
	Path   string `json:"path" binding:"required"`
}

// VerifyPermission 验证当前用户权限
// @Summary 验证当前用户权限
// @Description 验证当前用户是否有特定权限
// @Tags 权限验证
// @Accept json
// @Produce json
// @Param data body AuthVerifyRequest true "权限验证信息"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/verify [post]
func VerifyPermission(c *gin.Context) {
	var req AuthVerifyRequest
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

	// 检查权限
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
		"msg":  "验证成功",
		"data": gin.H{
			"has_permission": hasPermission,
		},
	})
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 刷新JWT令牌
// @Tags 权限验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func RefreshToken(c *gin.Context) {
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

	// 获取用户信息
	user, err := model.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "用户不存在",
			"data": nil,
		})
		return
	}

	// 生成新的JWT令牌
	token, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		zap.L().Error("生成Token失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成Token失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "刷新令牌成功",
		"data": gin.H{
			"token": token,
		},
	})
}

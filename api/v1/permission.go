package v1

import (
	"interviewGenius/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Method      string `json:"method" binding:"required,oneof=GET POST PUT PATCH DELETE"`
	PathPattern string `json:"path_pattern" binding:"required"`
}

// UpdatePermissionRequest 更新权限请求
type UpdatePermissionRequest struct {
	Method      string `json:"method" binding:"oneof=GET POST PUT PATCH DELETE"`
	PathPattern string `json:"path_pattern"`
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param data body CreatePermissionRequest true "权限信息"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/permissions [post]
func CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 创建权限
	permission := &model.Permission{
		Method:      req.Method,
		PathPattern: req.PathPattern,
	}

	if err := model.CreatePermission(permission); err != nil {
		zap.L().Error("创建权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "创建权限失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "创建成功",
		"data": gin.H{
			"id":           permission.ID,
			"method":       permission.Method,
			"path_pattern": permission.PathPattern,
		},
	})
}

// GetPermissionList 获取权限列表
// @Summary 获取权限列表
// @Description 获取所有权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/permissions [get]
func GetPermissionList(c *gin.Context) {
	permissions, err := model.GetPermissionList()
	if err != nil {
		zap.L().Error("获取权限列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取权限列表失败",
			"data": nil,
		})
		return
	}

	// 组装权限数据
	permissionData := make([]map[string]interface{}, 0, len(permissions))
	for _, permission := range permissions {
		permissionData = append(permissionData, map[string]interface{}{
			"id":           permission.ID,
			"method":       permission.Method,
			"path_pattern": permission.PathPattern,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"data": permissionData,
		},
	})
}

// GetPermissionDetail 获取权限详情
// @Summary 获取权限详情
// @Description 获取指定权限的详细信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/permissions/{id} [get]
func GetPermissionDetail(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的权限ID",
			"data": nil,
		})
		return
	}

	// 获取权限
	permission, err := model.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "权限不存在",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"id":           permission.ID,
			"method":       permission.Method,
			"path_pattern": permission.PathPattern,
		},
	})
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新指定权限的信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Param data body UpdatePermissionRequest true "权限信息"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/permissions/{id} [put]
func UpdatePermission(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的权限ID",
			"data": nil,
		})
		return
	}

	// 获取请求参数
	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取权限
	permission, err := model.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "权限不存在",
			"data": nil,
		})
		return
	}

	// 更新权限信息
	if req.Method != "" {
		permission.Method = req.Method
	}
	if req.PathPattern != "" {
		permission.PathPattern = req.PathPattern
	}

	// 更新权限
	if err := model.UpdatePermission(permission); err != nil {
		zap.L().Error("更新权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "更新失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "更新成功",
		"data": gin.H{
			"id":           permission.ID,
			"method":       permission.Method,
			"path_pattern": permission.PathPattern,
		},
	})
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除指定权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path string true "权限ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/permissions/{id} [delete]
func DeletePermission(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的权限ID",
			"data": nil,
		})
		return
	}

	// 删除权限
	if err := model.DeletePermission(id); err != nil {
		zap.L().Error("删除权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "删除失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "删除成功",
		"data": gin.H{
			"success": true,
		},
	})
}

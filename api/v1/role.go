package v1

import (
	"interviewGenius/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName string `json:"role_name" binding:"required,min=2,max=32"`
	IsSuper  bool   `json:"is_super"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName string `json:"role_name"`
	IsSuper  bool   `json:"is_super"`
}

// RolePermissionRequest 角色权限请求
type RolePermissionRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param data body CreateRoleRequest true "角色信息"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles [post]
func CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 检查角色名是否已存在
	_, err := model.GetRoleByName(req.RoleName)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "角色名已存在",
			"data": nil,
		})
		return
	}

	// 创建角色
	role := &model.Role{
		RoleName: req.RoleName,
		IsSuper:  req.IsSuper,
	}

	if err := model.CreateRole(role); err != nil {
		zap.L().Error("创建角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "创建角色失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "创建成功",
		"data": gin.H{
			"id":        role.ID,
			"role_name": role.RoleName,
			"is_super":  role.IsSuper,
		},
	})
}

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles [get]
func GetRoleList(c *gin.Context) {
	roles, err := model.GetRoleList()
	if err != nil {
		zap.L().Error("获取角色列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取角色列表失败",
			"data": nil,
		})
		return
	}

	// 组装角色数据
	roleData := make([]map[string]interface{}, 0, len(roles))
	for _, role := range roles {
		roleData = append(roleData, map[string]interface{}{
			"id":        role.ID,
			"role_name": role.RoleName,
			"is_super":  role.IsSuper,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"data": roleData,
		},
	})
}

// GetRoleDetail 获取角色详情
// @Summary 获取角色详情
// @Description 获取角色详情及权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [get]
func GetRoleDetail(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 获取角色及其权限
	role, err := model.GetRoleWithPermissions(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "角色不存在",
			"data": nil,
		})
		return
	}

	// 组装权限数据
	permissionData := make([]map[string]interface{}, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
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
			"id":          role.ID,
			"role_name":   role.RoleName,
			"is_super":    role.IsSuper,
			"permissions": permissionData,
		},
	})
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body UpdateRoleRequest true "角色信息"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [put]
func UpdateRole(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 获取请求参数
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取角色
	role, err := model.GetRoleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "角色不存在",
			"data": nil,
		})
		return
	}

	// 如果更新角色名，检查名称是否已存在
	if req.RoleName != "" && req.RoleName != role.RoleName {
		existingRole, err := model.GetRoleByName(req.RoleName)
		if err == nil && existingRole.ID != role.ID {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "角色名已存在",
				"data": nil,
			})
			return
		}
		role.RoleName = req.RoleName
	}

	// 更新超级管理员标志
	role.IsSuper = req.IsSuper

	// 更新角色
	if err := model.UpdateRole(role); err != nil {
		zap.L().Error("更新角色失败", zap.Error(err))
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
			"id":        role.ID,
			"role_name": role.RoleName,
			"is_super":  role.IsSuper,
		},
	})
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除指定角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles/{id} [delete]
func DeleteRole(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 删除角色
	if err := model.DeleteRole(uint(id)); err != nil {
		zap.L().Error("删除角色失败", zap.Error(err))
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

// AddRolePermissions 为角色添加权限
// @Summary 为角色分配权限
// @Description 为指定角色分配权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param data body RolePermissionRequest true "权限ID列表"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles/{id}/permissions [post]
func AddRolePermissions(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 获取请求参数
	var req RolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 为角色添加权限
	if err := model.AddPermissionsToRole(uint(id), req.PermissionIDs); err != nil {
		zap.L().Error("为角色添加权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "添加权限失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "添加权限成功",
		"data": gin.H{
			"success": true,
		},
	})
}

// GetRolePermissions 获取角色的权限
// @Summary 获取角色的权限
// @Description 获取指定角色的所有权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/v1/roles/{id}/permissions [get]
func GetRolePermissions(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 获取角色的权限
	permissions, err := model.GetRolePermissions(uint(id))
	if err != nil {
		zap.L().Error("获取角色权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取权限失败",
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

// RemoveRolePermission 移除角色的权限
// @Summary 从角色中移除权限
// @Description 从指定角色中移除指定权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param permissionId path string true "权限ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "成功移除权限"
// @Failure 400 {object} map[string]interface{} "参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /api/v1/roles/{id}/permissions/{permissionId} [delete]
func RemoveRolePermission(c *gin.Context) {
	// 从URL参数中获取roleId和permissionId
	roleIdStr := c.Param("id")
	permissionIdStr := c.Param("permissionId")

	// 转换roleId为uint类型
	roleId, err := strconv.ParseUint(roleIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "角色ID格式错误",
			"data": nil,
		})
		return
	}

	// 转换permissionId为uuid类型
	permissionId, err := uuid.Parse(permissionIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "权限ID格式错误",
			"data": nil,
		})
		return
	}

	// 调用模型层移除权限
	err = model.RemovePermissionFromRole(uint(roleId), permissionId)
	if err != nil {
		zap.L().Error("移除角色权限失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "成功移除权限",
		"data": gin.H{
			"success": true,
		},
	})
}

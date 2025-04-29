package v1

import (
	"interviewGenius/internal/dto"
	"interviewGenius/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService *service.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{
		roleService: service.NewRoleService(),
	}
}

// CreateRole 创建角色
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req dto.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	role, err := c.roleService.CreateRole(&req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, role)
}

// GetRoleList 获取角色列表
func (c *RoleController) GetRoleList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	roles, total, err := c.roleService.GetRoleList(page, limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"roles": roles,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetRoleDetail 获取角色详情
func (c *RoleController) GetRoleDetail(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	role, err := c.roleService.GetRoleDetail(uint(id))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, role)
}

// UpdateRole 更新角色
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := c.roleService.UpdateRole(uint(id), &req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
}

// DeleteRole 删除角色
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	if err := c.roleService.DeleteRole(uint(id)); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
}

// AddRolePermissions 添加角色权限
func (c *RoleController) AddRolePermissions(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	var req dto.RolePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := c.roleService.AddRolePermissions(uint(roleID), req.PermissionIDs); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
}

// GetRolePermissions 获取角色权限
func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	permissions, err := c.roleService.GetRolePermissions(uint(roleID))
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, permissions)
}

// RemoveRolePermission 移除角色权限
func (c *RoleController) RemoveRolePermission(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "无效的角色ID"})
		return
	}

	permissionID := ctx.Param("permissionId")
	if permissionID == "" {
		ctx.JSON(400, gin.H{"error": "无效的权限ID"})
		return
	}

	if err := c.roleService.RemoveRolePermission(uint(roleID), permissionID); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(200)
}

package dto

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	RoleName    string `json:"role_name" binding:"required,min=2,max=32"`
	Name        string `json:"name" binding:"required,min=2,max=32"`
	Description string `json:"description" binding:"max=255"`
	IsSuper     bool   `json:"is_super"`
}

// UpdateRoleRequest 更新角色请求
type UpdateRoleRequest struct {
	RoleName    string `json:"role_name" binding:"required,min=2,max=32"`
	Name        string `json:"name" binding:"required,min=2,max=32"`
	Description string `json:"description" binding:"max=255"`
	IsSuper     bool   `json:"is_super"`
}

// RoleResponse 角色响应
type RoleResponse struct {
	ID          uint             `json:"id"`
	RoleName    string           `json:"role_name"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	IsSuper     bool             `json:"is_super"`
	Permissions []PermissionInfo `json:"permissions,omitempty"`
}

// PermissionInfo 权限信息
type PermissionInfo struct {
	ID          string `json:"id"`
	Method      string `json:"method"`
	PathPattern string `json:"path_pattern"`
	Description string `json:"description,omitempty"`
}

// RolePermissionRequest 角色权限请求
type RolePermissionRequest struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"`
}

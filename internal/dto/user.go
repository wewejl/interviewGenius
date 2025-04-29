package dto

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6,max=30"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserRoleRequest 用户角色请求
type UserRoleRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID       string     `json:"id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Roles    []RoleInfo `json:"roles,omitempty"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	ID       uint   `json:"id"`
	RoleName string `json:"role_name"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

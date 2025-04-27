package model

import (
	"github.com/google/uuid"
)

// RolePermission 角色权限关联模型
type RolePermission struct {
	Model
	RoleID       uint      `gorm:"not null" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:char(36);not null" json:"permission_id"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permission"
}

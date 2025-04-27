package model

import (
	"github.com/google/uuid"
)

// UserRole 用户角色关联模型
type UserRole struct {
	UserID uuid.UUID `gorm:"type:char(36);primaryKey;column:user_id" json:"user_id"`
	RoleID uint      `gorm:"primaryKey;column:role_id" json:"role_id"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_role"
}

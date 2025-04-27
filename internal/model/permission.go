package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission 权限模型
type Permission struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Method      string    `json:"method" gorm:"size:8;not null"`
	PathPattern string    `json:"path_pattern" gorm:"size:128;not null"`
	Description string    `json:"description" gorm:"size:128"`
	Roles       []*Role   `json:"roles,omitempty" gorm:"many2many:role_permission;"`
}

// BeforeCreate 创建前生成UUID
func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// CreatePermission 创建权限
func CreatePermission(permission *Permission) error {
	return DB.Create(permission).Error
}

// GetPermissionByID 根据ID获取权限
func GetPermissionByID(id uuid.UUID) (*Permission, error) {
	var permission Permission
	if err := DB.Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionList 获取权限列表
func GetPermissionList() ([]*Permission, error) {
	var permissions []*Permission
	if err := DB.Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}

// UpdatePermission 更新权限
func UpdatePermission(permission *Permission) error {
	return DB.Model(&Permission{}).Where("id = ?", permission.ID).Updates(map[string]interface{}{
		"method":       permission.Method,
		"path_pattern": permission.PathPattern,
		"description":  permission.Description,
	}).Error
}

// DeletePermission 删除权限
func DeletePermission(id uuid.UUID) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除权限与角色的关联
		if err := tx.Exec("DELETE FROM role_permission WHERE permission_id = ?", id).Error; err != nil {
			return err
		}

		// 删除权限
		if err := tx.Where("id = ?", id).Delete(&Permission{}).Error; err != nil {
			return err
		}

		return nil
	})
}

// CheckUserPermission 检查用户是否有特定权限
func CheckUserPermission(userID uuid.UUID, method, path string) (bool, error) {
	var count int64

	// 复杂SQL查询，检查用户是否有匹配的权限
	// 1. 获取用户的所有角色
	// 2. 检查这些角色是否有匹配的权限
	query := `
		SELECT COUNT(*) 
		FROM permissions p
		JOIN role_permission rp ON p.id = rp.permission_id
		JOIN user_role ur ON rp.role_id = ur.role_id
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = ?
		AND (
			r.is_super = true 
			OR (
				p.method = ? 
				AND ? LIKE REPLACE(REPLACE(p.path_pattern, '*', '%'), '?', '_')
			)
		)
	`

	if err := DB.Raw(query, userID, method, path).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

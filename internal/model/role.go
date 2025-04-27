package model

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	Model
	RoleName    string        `gorm:"size:50;not null;unique" json:"role_name"`
	Name        string        `gorm:"size:50;not null" json:"name"`
	Description string        `gorm:"size:255" json:"description"`
	IsSuper     bool          `gorm:"default:false" json:"is_super"`
	Users       []*User       `gorm:"many2many:user_role;" json:"users,omitempty"`
	Permissions []*Permission `gorm:"many2many:role_permission;" json:"permissions,omitempty"`
}

// CreateRole 创建角色
func CreateRole(role *Role) error {
	return DB.Create(role).Error
}

// GetRoleByID 根据ID获取角色
func GetRoleByID(id uint) (*Role, error) {
	var role Role
	if err := DB.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleWithPermissions 获取角色及其权限
func GetRoleWithPermissions(id uint) (*Role, error) {
	var role Role
	if err := DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName 根据名称获取角色
func GetRoleByName(roleName string) (*Role, error) {
	var role Role
	if err := DB.Where("role_name = ?", roleName).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// GetRoleList 获取角色列表
func GetRoleList() ([]*Role, error) {
	var roles []*Role
	if err := DB.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// UpdateRole 更新角色
func UpdateRole(role *Role) error {
	return DB.Save(role).Error
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// 删除角色与用户的关联
		if err := tx.Exec("DELETE FROM user_role WHERE role_id = ?", id).Error; err != nil {
			return err
		}

		// 删除角色与权限的关联
		if err := tx.Exec("DELETE FROM role_permission WHERE role_id = ?", id).Error; err != nil {
			return err
		}

		// 删除角色
		if err := tx.Delete(&Role{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// AddPermissionsToRole 为角色添加权限
func AddPermissionsToRole(roleID uint, permissionIDs []string) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		var permissions []*Permission
		for _, idStr := range permissionIDs {
			var permission Permission
			if err := tx.Where("id = ?", idStr).First(&permission).Error; err != nil {
				return err
			}
			permissions = append(permissions, &permission)
		}

		if err := tx.Model(&role).Association("Permissions").Append(permissions); err != nil {
			return err
		}

		return nil
	})
}

// GetRolePermissions 获取角色的权限
func GetRolePermissions(roleID uint) ([]*Permission, error) {
	var role Role
	if err := DB.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	return role.Permissions, nil
}

// RemovePermissionFromRole 从角色中移除权限
func RemovePermissionFromRole(roleID uint, permissionID uuid.UUID) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var role Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		var permission Permission
		if err := tx.Where("id = ?", permissionID).First(&permission).Error; err != nil {
			return err
		}

		if err := tx.Model(&role).Association("Permissions").Delete(&permission); err != nil {
			return err
		}

		return nil
	})
}

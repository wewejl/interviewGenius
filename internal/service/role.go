package service

import (
	"errors"
	"github.com/google/uuid"
	"interviewGenius/internal/dto"
	"interviewGenius/internal/model"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// CreateRole 创建角色
func (s *RoleService) CreateRole(req *dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	// 检查角色名是否已存在
	_, err := model.GetRoleByName(req.RoleName)
	if err == nil {
		return nil, errors.New("角色名已存在")
	}

	role := &model.Role{
		RoleName:    req.RoleName,
		Name:        req.Name,
		Description: req.Description,
		IsSuper:     req.IsSuper,
	}

	if err := model.CreateRole(role); err != nil {
		return nil, err
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Name:        role.Name,
		Description: role.Description,
		IsSuper:     role.IsSuper,
	}, nil
}

// GetRoleList 获取角色列表
func (s *RoleService) GetRoleList(page, limit int) ([]*dto.RoleResponse, int64, error) {
	roles, total, err := model.GetRoleList(page, limit)
	if err != nil {
		return nil, 0, err
	}

	roleList := make([]*dto.RoleResponse, 0)
	for _, role := range roles {
		permissions := make([]dto.PermissionInfo, 0)
		for _, perm := range role.Permissions {
			permissions = append(permissions, dto.PermissionInfo{
				ID:          perm.ID.String(),
				Method:      perm.Method,
				PathPattern: perm.PathPattern,
				Description: perm.Description,
			})
		}

		roleList = append(roleList, &dto.RoleResponse{
			ID:          role.ID,
			RoleName:    role.RoleName,
			Name:        role.Name,
			Description: role.Description,
			IsSuper:     role.IsSuper,
			Permissions: permissions,
		})
	}

	return roleList, total, nil
}

// GetRoleDetail 获取角色详情
func (s *RoleService) GetRoleDetail(id uint) (*dto.RoleResponse, error) {
	role, err := model.GetRoleWithPermissions(id)
	if err != nil {
		return nil, err
	}

	permissions := make([]dto.PermissionInfo, 0)
	for _, perm := range role.Permissions {
		permissions = append(permissions, dto.PermissionInfo{
			ID:          perm.ID.String(),
			Method:      perm.Method,
			PathPattern: perm.PathPattern,
			Description: perm.Description,
		})
	}

	return &dto.RoleResponse{
		ID:          role.ID,
		RoleName:    role.RoleName,
		Name:        role.Name,
		Description: role.Description,
		IsSuper:     role.IsSuper,
		Permissions: permissions,
	}, nil
}

// UpdateRole 更新角色
func (s *RoleService) UpdateRole(id uint, req *dto.UpdateRoleRequest) error {
	role, err := model.GetRoleByID(id)
	if err != nil {
		return err
	}

	// 检查角色名是否已存在
	existingRole, err := model.GetRoleByName(req.RoleName)
	if err == nil && existingRole.ID != id {
		return errors.New("角色名已存在")
	}

	role.RoleName = req.RoleName
	role.Name = req.Name
	role.Description = req.Description
	role.IsSuper = req.IsSuper

	return model.UpdateRole(role)
}

// DeleteRole 删除角色
func (s *RoleService) DeleteRole(id uint) error {
	return model.DeleteRole(id)
}

// AddRolePermissions 添加角色权限
func (s *RoleService) AddRolePermissions(roleID uint, permissionIDs []string) error {
	return model.AddPermissionsToRole(roleID, permissionIDs)
}

// GetRolePermissions 获取角色权限
func (s *RoleService) GetRolePermissions(roleID uint) ([]dto.PermissionInfo, error) {
	permissions, err := model.GetRolePermissions(roleID)
	if err != nil {
		return nil, err
	}

	permissionList := make([]dto.PermissionInfo, 0)
	for _, perm := range permissions {
		permissionList = append(permissionList, dto.PermissionInfo{
			ID:          perm.ID.String(),
			Method:      perm.Method,
			PathPattern: perm.PathPattern,
		})
	}

	return permissionList, nil
}

// RemoveRolePermission 移除角色权限
func (s *RoleService) RemoveRolePermission(roleID uint, permissionID string) error {
	permID, err := uuid.Parse(permissionID)
	if err != nil {
		return errors.New("无效的权限ID")
	}
	return model.RemovePermissionFromRole(roleID, permID)
}

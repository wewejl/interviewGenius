package service

import (
	"errors"
	"interviewGenius/internal/dto"
	"interviewGenius/internal/model"
	"interviewGenius/internal/pkg/util"

	"github.com/google/uuid"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// Register 注册用户
func (s *UserService) Register(req *dto.RegisterRequest) (*dto.TokenResponse, error) {
	// 检查用户名是否已存在
	_, err := model.GetUserByUsername(req.Username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := model.CreateUser(user); err != nil {
		return nil, err
	}

	// 生成Token
	token, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Token:    token,
	}, nil
}

// Login 用户登录
func (s *UserService) Login(req *dto.LoginRequest) (*dto.TokenResponse, error) {
	// 获取用户
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成Token
	token, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		return nil, err
	}

	return &dto.TokenResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Token:    token,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(id uuid.UUID, currentUserID uuid.UUID) (*dto.UserResponse, error) {
	// 获取用户信息
	user, err := model.GetUserByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 检查权限
	if currentUserID != id {
		// TODO: 实现详细的权限检查逻辑
	}

	// 组装角色信息
	roles := make([]dto.RoleInfo, 0)
	for _, role := range user.Roles {
		roles = append(roles, dto.RoleInfo{
			ID:       role.ID,
			RoleName: role.RoleName,
		})
	}

	return &dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
	}, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uuid.UUID, req *dto.UpdateUserRequest) error {
	user, err := model.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.OldPassword != "" && req.NewPassword != "" {
		if !user.CheckPassword(req.OldPassword) {
			return errors.New("原密码错误")
		}
		user.Password = req.NewPassword
	}

	return model.UpdateUser(user)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return model.DeleteUser(id)
}

// AddUserRoles 添加用户角色
func (s *UserService) AddUserRoles(userID uuid.UUID, roleIDs []uint) error {
	return model.AddRolesToUser(userID, roleIDs)
}

// GetUserRoles 获取用户角色
func (s *UserService) GetUserRoles(userID uuid.UUID) ([]dto.RoleInfo, error) {
	roles, err := model.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	roleInfos := make([]dto.RoleInfo, 0)
	for _, role := range roles {
		roleInfos = append(roleInfos, dto.RoleInfo{
			ID:       role.ID,
			RoleName: role.RoleName,
		})
	}
	return roleInfos, nil
}

// RemoveUserRole 移除用户角色
func (s *UserService) RemoveUserRole(userID uuid.UUID, roleID uint) error {
	return model.RemoveRoleFromUser(userID, roleID)
}

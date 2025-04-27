package model

import (
	"go.uber.org/zap"
)

// API路由权限定义
var apiRoutes = []Permission{
	// 用户管理
	{Method: "POST", PathPattern: "/api/v1/user", Description: "创建用户"},
	{Method: "PUT", PathPattern: "/api/v1/user/:id", Description: "更新用户"},
	{Method: "DELETE", PathPattern: "/api/v1/user/:id", Description: "删除用户"},
	{Method: "GET", PathPattern: "/api/v1/user/:id", Description: "获取用户详情"},
	{Method: "GET", PathPattern: "/api/v1/users", Description: "获取用户列表"},

	// 角色管理
	{Method: "POST", PathPattern: "/api/v1/role", Description: "创建角色"},
	{Method: "PUT", PathPattern: "/api/v1/role/:id", Description: "更新角色"},
	{Method: "DELETE", PathPattern: "/api/v1/role/:id", Description: "删除角色"},
	{Method: "GET", PathPattern: "/api/v1/role/:id", Description: "获取角色详情"},
	{Method: "GET", PathPattern: "/api/v1/roles", Description: "获取角色列表"},
	{Method: "POST", PathPattern: "/api/v1/role/:id/permissions", Description: "分配角色权限"},

	// 权限管理
	{Method: "GET", PathPattern: "/api/v1/permissions", Description: "获取权限列表"},

	// 权限验证
	{Method: "POST", PathPattern: "/api/v1/auth/login", Description: "用户登录"},
	{Method: "POST", PathPattern: "/api/v1/auth/logout", Description: "用户登出"},
	{Method: "GET", PathPattern: "/api/v1/auth/info", Description: "获取当前用户信息"},
}

// 普通用户可访问的路由
var regularUserRoutes = []string{
	"/api/v1/auth/login",
	"/api/v1/auth/logout",
	"/api/v1/auth/info",
}

// 会员卡类型定义
var memberCardTypes = []MemberCard{
	{Name: "日卡", DurationDays: 1, Price: 998, Description: "24小时内无限次使用服务"},
	{Name: "周卡", DurationDays: 7, Price: 2998, Description: "7天内无限次使用服务"},
	{Name: "双周卡", DurationDays: 14, Price: 4998, Description: "14天内无限次使用服务"},
	{Name: "月卡", DurationDays: 30, Price: 8998, Description: "30天内无限次使用服务"},
	{Name: "双月卡", DurationDays: 60, Price: 15998, Description: "60天内无限次使用服务"},
}

// InitData 初始化系统数据
func InitData() error {
	// 检查权限表是否为空
	var permissionCount int64
	if err := DB.Model(&Permission{}).Count(&permissionCount).Error; err != nil {
		zap.L().Error("检查权限表失败", zap.Error(err))
		return err
	}

	// 检查角色表是否为空
	var roleCount int64
	if err := DB.Model(&Role{}).Count(&roleCount).Error; err != nil {
		zap.L().Error("检查角色表失败", zap.Error(err))
		return err
	}

	// 检查会员卡表是否为空
	var memberCardCount int64
	if err := DB.Model(&MemberCard{}).Count(&memberCardCount).Error; err != nil {
		zap.L().Error("检查会员卡表失败", zap.Error(err))
		return err
	}

	// 如果权限和角色表都为空，则初始化系统数据
	if permissionCount == 0 && roleCount == 0 {
		// 开始事务
		tx := DB.Begin()

		// 创建所有API路由权限
		for _, route := range apiRoutes {
			if err := tx.Create(&route).Error; err != nil {
				tx.Rollback()
				zap.L().Error("创建API路由权限失败", zap.Error(err))
				return err
			}
		}

		// 创建管理员角色
		adminRole := Role{
			RoleName:    "admin",
			Name:        "管理员",
			Description: "系统管理员，拥有所有权限",
			IsSuper:     true,
		}
		if err := tx.Create(&adminRole).Error; err != nil {
			tx.Rollback()
			zap.L().Error("创建管理员角色失败", zap.Error(err))
			return err
		}

		// 为管理员角色分配所有权限
		var permissions []Permission
		if err := tx.Find(&permissions).Error; err != nil {
			tx.Rollback()
			zap.L().Error("查询所有权限失败", zap.Error(err))
			return err
		}

		// 创建管理员角色与权限的关联
		for _, permission := range permissions {
			rolePermission := RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: permission.ID,
			}
			if err := tx.Create(&rolePermission).Error; err != nil {
				tx.Rollback()
				zap.L().Error("创建角色权限关联失败", zap.Error(err))
				return err
			}
		}

		// 创建普通用户角色
		regularRole := Role{
			RoleName:    "regular_user",
			Name:        "普通用户",
			Description: "普通用户，拥有基本权限",
		}
		if err := tx.Create(&regularRole).Error; err != nil {
			tx.Rollback()
			zap.L().Error("创建普通用户角色失败", zap.Error(err))
			return err
		}

		// 为普通用户角色分配基本权限
		for _, routePath := range regularUserRoutes {
			var permission Permission
			if err := tx.Where("path_pattern = ?", routePath).First(&permission).Error; err != nil {
				tx.Rollback()
				zap.L().Error("查询基本权限失败", zap.String("path", routePath), zap.Error(err))
				return err
			}

			rolePermission := RolePermission{
				RoleID:       regularRole.ID,
				PermissionID: permission.ID,
			}
			if err := tx.Create(&rolePermission).Error; err != nil {
				tx.Rollback()
				zap.L().Error("创建普通用户角色权限关联失败", zap.Error(err))
				return err
			}
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			zap.L().Error("提交事务失败", zap.Error(err))
			return err
		}

		zap.L().Info("系统初始化数据成功",
			zap.Int("permissions_count", len(apiRoutes)),
			zap.String("admin_role", adminRole.Name),
			zap.String("regular_role", regularRole.Name))
	}

	// 如果会员卡表为空，则初始化会员卡数据
	if memberCardCount == 0 {
		// 开始事务
		tx := DB.Begin()

		// 创建所有会员卡类型
		for _, card := range memberCardTypes {
			if err := tx.Create(&card).Error; err != nil {
				tx.Rollback()
				zap.L().Error("创建会员卡类型失败", zap.Error(err))
				return err
			}
		}

		// 提交事务
		if err := tx.Commit().Error; err != nil {
			zap.L().Error("提交事务失败", zap.Error(err))
			return err
		}

		zap.L().Info("会员卡数据初始化成功", zap.Int("card_types_count", len(memberCardTypes)))
	}

	return nil
}

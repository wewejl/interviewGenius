package v1

import (
	"interviewGenius/internal/model"
	"interviewGenius/internal/pkg/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

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

// Register 注册用户
// @Summary 用户注册
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body RegisterRequest true "用户注册信息"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /api/v1/users/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 检查用户名是否已存在
	_, err := model.GetUserByUsername(req.Username)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "用户名已存在",
			"data": nil,
		})
		return
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	if err := model.CreateUser(user); err != nil {
		zap.L().Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "创建用户失败",
			"data": nil,
		})
		return
	}

	// 生成Token
	token, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		zap.L().Error("生成Token失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成Token失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"token":    token,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "用户登录信息"
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Router /api/v1/users/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取用户
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		return
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "用户名或密码错误",
			"data": nil,
		})
		return
	}

	// 生成Token
	token, err := util.GenerateToken(user.ID.String(), user.Username)
	if err != nil {
		zap.L().Error("生成Token失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "生成Token失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"token":    token,
		},
	})
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /api/v1/users/{id} [get]
func GetUserInfo(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	// 获取用户信息
	user, err := model.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "用户不存在",
			"data": nil,
		})
		return
	}

	// 检查是否是获取自己的信息或有权限
	currentID, err := uuid.Parse(currentUserID.(string))
	hasPermission := false

	if currentID == id { // 是自己
		hasPermission = true
	} else {
		// 这里应该有权限检查逻辑
		// 暂不实现详细的权限检查
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "没有权限访问",
			"data": nil,
		})
		return
	}

	// 组装角色信息
	roleData := make([]map[string]interface{}, 0)
	for _, role := range user.Roles {
		roleData = append(roleData, map[string]interface{}{
			"id":        role.ID,
			"role_name": role.RoleName,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"roles":    roleData,
		},
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /api/v1/users [get]
func GetUserList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 获取用户列表
	users, total, err := model.GetUserList(page, limit)
	if err != nil {
		zap.L().Error("获取用户列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取用户列表失败",
			"data": nil,
		})
		return
	}

	// 组装用户数据
	userData := make([]map[string]interface{}, 0)
	for _, user := range users {
		userData = append(userData, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"total": total,
			"data":  userData,
		},
	})
}

// UpdateUser 更新用户信息
// @Summary 更新用户信息
// @Description 更新用户的信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param data body UpdateUserRequest true "更新用户信息"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Router /api/v1/users/{id} [put]
func UpdateUser(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 获取请求参数
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	// 获取用户信息
	user, err := model.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "用户不存在",
			"data": nil,
		})
		return
	}

	// 检查是否是更新自己的信息或有权限
	currentID, err := uuid.Parse(currentUserID.(string))
	hasPermission := false

	if currentID == id { // 是自己
		hasPermission = true
	} else {
		// 这里应该有权限检查逻辑
		// 暂不实现详细的权限检查
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, gin.H{
			"code": http.StatusForbidden,
			"msg":  "没有权限更新",
			"data": nil,
		})
		return
	}

	// 如果要修改密码，需要验证旧密码
	if req.NewPassword != "" {
		// 验证旧密码
		if req.OldPassword == "" || !user.CheckPassword(req.OldPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusBadRequest,
				"msg":  "旧密码错误",
				"data": nil,
			})
			return
		}
		user.Password = req.NewPassword
	}

	// 更新用户信息
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := model.UpdateUser(user); err != nil {
		zap.L().Error("更新用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "更新失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "更新成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除指定的用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 401 {object} util.Response
// @Failure 403 {object} util.Response
// @Router /api/v1/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 删除用户
	if err := model.DeleteUser(id); err != nil {
		zap.L().Error("删除用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "删除失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "删除成功",
		"data": gin.H{
			"success": true,
		},
	})
}

// AddUserRoles 为用户分配角色
// @Summary 为用户分配角色
// @Description 为指定用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param data body UserRoleRequest true "角色ID列表"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Router /api/v1/users/{id}/roles [post]
func AddUserRoles(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 获取请求参数
	var req UserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	// 为用户添加角色
	if err := model.AddRolesToUser(id, req.RoleIDs); err != nil {
		zap.L().Error("为用户添加角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "添加角色失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "添加角色成功",
		"data": gin.H{
			"success": true,
		},
	})
}

// GetUserRoles 获取用户的角色
// @Summary 获取用户的角色
// @Description 获取指定用户的所有角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Router /api/v1/users/{id}/roles [get]
func GetUserRoles(c *gin.Context) {
	// 获取路径参数
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	// 获取用户的角色
	roles, err := model.GetUserRoles(id)
	if err != nil {
		zap.L().Error("获取用户角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取角色失败",
			"data": nil,
		})
		return
	}

	// 组装角色数据
	roleData := make([]map[string]interface{}, 0)
	for _, role := range roles {
		roleData = append(roleData, map[string]interface{}{
			"id":        role.ID,
			"role_name": role.RoleName,
			"is_super":  role.IsSuper,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"data": roleData,
		},
	})
}

// RemoveUserRole 移除用户的角色
// @Summary 移除用户的角色
// @Description 移除用户的特定角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param roleId path int true "角色ID"
// @Security BearerAuth
// @Success 200 {object} util.Response
// @Failure 400 {object} util.Response
// @Failure 401 {object} util.Response
// @Router /api/v1/users/{id}/roles/{roleId} [delete]
func RemoveUserRole(c *gin.Context) {
	// 获取路径参数
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	// 移除用户的角色
	if err := model.RemoveRoleFromUser(userID, uint(roleID)); err != nil {
		zap.L().Error("移除用户角色失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "移除角色失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "移除角色成功",
		"data": gin.H{
			"success": true,
		},
	})
}

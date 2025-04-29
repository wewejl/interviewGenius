package v1

import (
	"interviewGenius/internal/dto"
	"interviewGenius/internal/model"
	"interviewGenius/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// Register 注册用户
// @Summary 用户注册
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.RegisterRequest true "用户注册信息"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} dto.Response
// @Router /api/v1/users/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	resp, err := c.userService.Register(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
		"data": resp,
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.LoginRequest true "用户登录信息"
// @Success 200 {object} dto.TokenResponse
// @Failure 400 {object} dto.Response
// @Router /api/v1/users/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	resp, err := c.userService.Login(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "登录成功",
		"data": resp,
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
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.Response
// @Failure 403 {object} dto.Response
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	currentUserID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "未认证",
			"data": nil,
		})
		return
	}

	currentID, err := uuid.Parse(currentUserID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "系统错误",
			"data": nil,
		})
		return
	}

	resp, err := c.userService.GetUserInfo(id, currentID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": resp,
	})
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	if err := c.userService.UpdateUser(id, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "更新成功",
		"data": nil,
	})
}

// DeleteUser 删除用户
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	if err := c.userService.DeleteUser(id); err != nil {
		zap.L().Error("删除用户失败", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "删除用户失败",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "删除成功",
		"data": nil,
	})
}

// AddUserRoles 为用户添加角色
func (c *UserController) AddUserRoles(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	var req dto.UserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的请求参数",
			"data": nil,
		})
		return
	}

	if err := c.userService.AddUserRoles(id, req.RoleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "添加角色失败",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "添加成功",
		"data": nil,
	})
}

// GetUserRoles 获取用户角色
func (c *UserController) GetUserRoles(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	roles, err := c.userService.GetUserRoles(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取角色失败",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": roles,
	})
}

// RemoveUserRole 移除用户角色
func (c *UserController) RemoveUserRole(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的用户ID",
			"data": nil,
		})
		return
	}

	roleIDStr := ctx.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "无效的角色ID",
			"data": nil,
		})
		return
	}

	if err := c.userService.RemoveUserRole(userID, uint(roleID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "移除角色失败",
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "移除成功",
		"data": nil,
	})
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 分页获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Router /api/v1/users [get]
func (c *UserController) GetUserList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	users, total, err := model.GetUserList(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取用户列表失败",
			"data": nil,
		})
		return
	}

	// 转换为响应格式
	userList := make([]dto.UserResponse, 0)
	for _, user := range users {
		roles := make([]dto.RoleInfo, 0)
		for _, role := range user.Roles {
			roles = append(roles, dto.RoleInfo{
				ID:       role.ID,
				RoleName: role.RoleName,
			})
		}

		userList = append(userList, dto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Roles:    roles,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "获取成功",
		"data": gin.H{
			"list":  userList,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

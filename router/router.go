package router

import (
	"interviewGenius/api/v1"
	"interviewGenius/docs"
	"interviewGenius/internal/middleware"
	"interviewGenius/pkg/setting"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	// 设置运行模式
	gin.SetMode(setting.ServerSetting.RunMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger 文档
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 直接提供swagger.json
	r.GET("/swagger.json", func(c *gin.Context) {
		c.String(http.StatusOK, docs.SwaggerInfo.SwaggerTemplate)
	})

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 用户管理
		users := apiV1.Group("/users")
		{
			// 无需认证的接口
			users.POST("/register", v1.Register)
			users.POST("/login", v1.Login)

			// 需要认证的接口
			users.Use(middleware.JWT())
			{
				users.GET("", v1.GetUserList)
				users.GET("/:id", v1.GetUserInfo)
				users.PUT("/:id", v1.UpdateUser)
				users.DELETE("/:id", v1.DeleteUser)

				// 用户角色管理
				users.POST("/:id/roles", v1.AddUserRoles)
				users.GET("/:id/roles", v1.GetUserRoles)
				users.DELETE("/:id/roles/:roleId", v1.RemoveUserRole)
			}
		}

		// 角色管理
		roles := apiV1.Group("/roles")
		roles.Use(middleware.JWT())
		{
			roles.POST("", v1.CreateRole)
			roles.GET("", v1.GetRoleList)
			roles.GET("/:id", v1.GetRoleDetail)
			roles.PUT("/:id", v1.UpdateRole)
			roles.DELETE("/:id", v1.DeleteRole)

			// 角色权限管理
			roles.POST("/:id/permissions", v1.AddRolePermissions)
			roles.GET("/:id/permissions", v1.GetRolePermissions)
			roles.DELETE("/:id/permissions/:permissionId", v1.RemoveRolePermission)
		}

		// 权限管理
		permissions := apiV1.Group("/permissions")
		permissions.Use(middleware.JWT())
		{
			permissions.POST("", v1.CreatePermission)
			permissions.GET("", v1.GetPermissionList)
			permissions.GET("/:id", v1.GetPermissionDetail)
			permissions.PUT("/:id", v1.UpdatePermission)
			permissions.DELETE("/:id", v1.DeletePermission)
		}

		// 权限验证
		auth := apiV1.Group("/auth")
		auth.Use(middleware.JWT())
		{
			auth.POST("/verify", v1.VerifyPermission)
			auth.POST("/refresh", v1.RefreshToken)
		}

		// 会员相关接口
		member := apiV1.Group("/member")
		{
			// 无需认证的接口
			member.GET("/cards", v1.GetMemberCards)

			// 需要认证的接口
			member.Use(middleware.JWT())
			{
				member.GET("/info", v1.GetMemberInfo)
				member.POST("/order", v1.CreateOrder)
				member.POST("/order/:id/pay", v1.PayOrder)
				member.GET("/orders", v1.GetOrders)
				member.GET("/check", v1.CheckServiceAccess)
			}
		}
	}

	return r
}

package router

import (
	"interviewGenius/docs"
	v1 "interviewGenius/internal/controller/v1"
	"interviewGenius/internal/middleware"
	"interviewGenius/internal/pkg/setting"
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

	// 初始化控制器
	userController := v1.NewUserController()
	roleController := v1.NewRoleController()
	paymentController := v1.NewPaymentController()

	// API v1
	apiV1 := r.Group("/api/v1")
	{
		// 用户管理
		users := apiV1.Group("/users")
		{
			// 无需认证的接口
			users.POST("/register", userController.Register)
			users.POST("/login", userController.Login)

			// 需要认证的接口
			users.Use(middleware.JWT())
			{
				users.GET("", userController.GetUserList)
				users.GET("/:id", userController.GetUserInfo)
				users.PUT("/:id", userController.UpdateUser)
				users.DELETE("/:id", userController.DeleteUser)

				// 用户角色管理
				users.POST("/:id/roles", userController.AddUserRoles)
				users.GET("/:id/roles", userController.GetUserRoles)
				users.DELETE("/:id/roles/:roleId", userController.RemoveUserRole)
			}
		}

		// 角色管理
		roles := apiV1.Group("/roles")
		roles.Use(middleware.JWT())
		{
			roles.POST("", roleController.CreateRole)
			roles.GET("", roleController.GetRoleList)
			roles.GET("/:id", roleController.GetRoleDetail)
			roles.PUT("/:id", roleController.UpdateRole)
			roles.DELETE("/:id", roleController.DeleteRole)

			// 角色权限管理
			roles.POST("/:id/permissions", roleController.AddRolePermissions)
			roles.GET("/:id/permissions", roleController.GetRolePermissions)
			roles.DELETE("/:id/permissions/:permissionId", roleController.RemoveRolePermission)
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

		// 支付相关路由
		paymentGroup := apiV1.Group("/payment")
		{
			paymentGroup.POST("/create", paymentController.CreatePayment)
			paymentGroup.POST("/notify", paymentController.HandlePaymentNotify)
			paymentGroup.GET("/return", paymentController.HandlePaymentReturn)
		}
	}

	return r
}

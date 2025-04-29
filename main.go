package main

import (
	"fmt"
	"interviewGenius/config"
	"interviewGenius/internal/model"
	"interviewGenius/internal/pkg/setting"
	"interviewGenius/internal/router"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// @title           InterviewGenius API
// @version         1.0
// @description     面试系统后端API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 初始化配置
	if err := setting.Init(); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		return
	}

	// 初始化日志
	if err := config.InitLogger(); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		return
	}
	defer zap.L().Sync()

	// 初始化数据库
	if err := model.Init(); err != nil {
		zap.L().Error("初始化数据库失败", zap.Error(err))
		return
	}

	// 初始化路由
	r := router.InitRouter()

	// 配置服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.AppSetting.Port),
		Handler:        r,
		ReadTimeout:    time.Duration(setting.ServerSetting.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(setting.ServerSetting.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务
	zap.L().Info(fmt.Sprintf("服务启动成功，监听端口: %d", setting.AppSetting.Port))
	if err := server.ListenAndServe(); err != nil {
		zap.L().Fatal("服务启动失败", zap.Error(err))
	}
}

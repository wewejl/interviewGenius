package model

import (
	"fmt"
	"interviewGenius/internal/pkg/setting"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go.uber.org/zap"
)

var DB *gorm.DB

// Model 基础模型
type Model struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// Init 初始化数据库连接
func Init() error {
	var err error

	// 数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Port,
		setting.DatabaseSetting.Name)

	// 日志配置
	logLevel := logger.Info
	if setting.ServerSetting.RunMode == "release" {
		logLevel = logger.Error
	}

	// 打开数据库连接
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.DatabaseSetting.TablePrefix,
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		zap.L().Error("连接数据库失败", zap.Error(err))
		return err
	}

	// 迁移数据库表
	if err = DB.AutoMigrate(&User{}, &Role{}, &Permission{}, &RolePermission{}, &MemberCard{}, &Order{}); err != nil {
		zap.L().Error("迁移数据库表失败", zap.Error(err))
		return err
	}

	// 确保User和Role的多对多关系表被创建
	if err = DB.SetupJoinTable(&User{}, "Roles", &UserRole{}); err != nil {
		zap.L().Error("设置User-Role关联表失败", zap.Error(err))
		return err
	}

	// 初始化系统数据（权限、角色）
	if err = InitData(); err != nil {
		zap.L().Error("初始化系统数据失败", zap.Error(err))
		return err
	}

	zap.L().Info("数据库连接和初始化数据成功")
	return nil
}

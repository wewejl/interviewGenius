package setting

import (
	"github.com/spf13/viper"
	"time"
)

type App struct {
	JwtSecret string
	Port      int
}

type Server struct {
	RunMode      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Port        int
	Name        string
	TablePrefix string
}

var (
	AppSetting      = &App{}
	ServerSetting   = &Server{}
	DatabaseSetting = &Database{}
)

// Init 初始化配置
func Init() error {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 加载App配置
	if err := viper.UnmarshalKey("app", AppSetting); err != nil {
		return err
	}

	// 加载Server配置
	if err := viper.UnmarshalKey("server", ServerSetting); err != nil {
		return err
	}

	// 加载Database配置
	if err := viper.UnmarshalKey("database", DatabaseSetting); err != nil {
		return err
	}

	return nil
}

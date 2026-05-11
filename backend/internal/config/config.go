package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 聚合应用启动所需的全部配置，字段名需要与配置文件和环境变量映射保持一致。
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RustFS   RustFSConfig   `mapstructure:"rustfs"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 描述 HTTP 服务与跨域策略。
type ServerConfig struct {
	Env            string   `mapstructure:"env"`
	Addr           string   `mapstructure:"addr"`
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
}

// DatabaseConfig 保存 PostgreSQL 连接配置，DSN 属于敏感信息，不应返回给前端。
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

// RedisConfig 保存 Redis 连接配置，当前用于刷新令牌等会话状态。
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// RustFSConfig 保存对象存储配置，密钥只能在服务端初始化客户端时使用。
type RustFSConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"useSSL"`
}

// JWTConfig 控制访问令牌和刷新令牌的签发密钥与有效期。
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessTTLMinutes int    `mapstructure:"accessTTLMinutes"`
	RefreshTTLDays   int    `mapstructure:"refreshTTLDays"`
}

// AccessTTL 返回访问令牌有效期，单位由配置中的分钟转换为 time.Duration。
func (c JWTConfig) AccessTTL() time.Duration {
	return time.Duration(c.AccessTTLMinutes) * time.Minute
}

// RefreshTTL 返回刷新令牌有效期，单位由配置中的天转换为 time.Duration。
func (c JWTConfig) RefreshTTL() time.Duration {
	return time.Duration(c.RefreshTTLDays) * 24 * time.Hour
}

// LogConfig 控制 zap 日志级别。
type LogConfig struct {
	Level string `mapstructure:"level"`
}

// Load 从配置文件加载配置，并允许 ADMIN_PLATFORM_ 前缀环境变量覆盖同名配置项。
func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetEnvPrefix("ADMIN_PLATFORM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

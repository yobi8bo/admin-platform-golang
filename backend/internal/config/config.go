package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RustFS   RustFSConfig   `mapstructure:"rustfs"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Env            string   `mapstructure:"env"`
	Addr           string   `mapstructure:"addr"`
	AllowedOrigins []string `mapstructure:"allowedOrigins"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RustFSConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"accessKey"`
	SecretKey string `mapstructure:"secretKey"`
	Bucket    string `mapstructure:"bucket"`
	UseSSL    bool   `mapstructure:"useSSL"`
}

type JWTConfig struct {
	Secret           string `mapstructure:"secret"`
	AccessTTLMinutes int    `mapstructure:"accessTTLMinutes"`
	RefreshTTLDays   int    `mapstructure:"refreshTTLDays"`
}

func (c JWTConfig) AccessTTL() time.Duration {
	return time.Duration(c.AccessTTLMinutes) * time.Minute
}

func (c JWTConfig) RefreshTTL() time.Duration {
	return time.Duration(c.RefreshTTLDays) * 24 * time.Hour
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

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

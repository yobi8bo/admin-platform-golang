package logger

import (
	"admin-platform/backend/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New 根据配置创建 zap logger，debug 级别使用开发格式便于本地排查。
func New(cfg config.LogConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	if cfg.Level == "debug" {
		zapCfg = zap.NewDevelopmentConfig()
		zapCfg.Level = zap.NewAtomicLevelAt(level)
	}
	return zapCfg.Build()
}

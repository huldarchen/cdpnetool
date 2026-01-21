package db

import (
	"context"
	"time"

	"cdpnetool/internal/logger"

	glog "gorm.io/gorm/logger"
)

// Logger 自定义GORM logger实现，对接项目统一日志系统
type Logger struct {
	internalLogger logger.Logger
	LogLevel       glog.LogLevel
}

// NewLogger 创建新的 Logger 实例
func NewLogger(l logger.Logger) *Logger {
	return &Logger{
		internalLogger: l,
		LogLevel:       glog.Info,
	}
}

// LogMode 实现 logger.Interface 接口
func (l *Logger) LogMode(level glog.LogLevel) glog.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 打印 info 级别日志
func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= glog.Info {
		l.internalLogger.Info(msg, data...)
	}
}

// Warn 打印 warn 级别日志
func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= glog.Warn {
		l.internalLogger.Warn(msg, data...)
	}
}

// Error 打印 error 级别日志
func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	if l.LogLevel >= glog.Error {
		l.internalLogger.Error(msg, data...)
	}
}

// Trace 打印 SQL 执行详情（核心审计点）
func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= glog.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []any{
		"sql", sql,
		"rows", rows,
		"timeMs", float64(elapsed.Nanoseconds()) / 1e6,
	}

	switch {
	case err != nil && l.LogLevel >= glog.Error:
		l.internalLogger.Error("SQL执行错误", append(fields, "error", err)...)
	case elapsed > time.Second && l.LogLevel >= glog.Warn:
		l.internalLogger.Warn("慢SQL查询", append(fields, "threshold", "1s")...)
	case l.LogLevel == glog.Info:
		l.internalLogger.Debug("SQL执行", fields...)
	}
}

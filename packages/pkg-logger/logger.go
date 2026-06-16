package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 常量定义，确保全局一致的 Context Key
type contextKey string

const (
	TraceIDKey  contextKey = "trace_id"
	TenantIDKey contextKey = "tenant_id"
	UserIDKey   contextKey = "user_id"
)

var globalLogger *zap.Logger

// Init 初始化全局日志记录器
// 必须在微服务启动时调用，例如：logger.Init("service-gateway", "info")
func Init(serviceName string, level string) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 符合世界级监控采集标准的时间戳
	encoderConfig.MessageKey = "message"

	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		logLevel = zapcore.InfoLevel
	}

	// 统一输出为 JSON 格式，方便 ELK 等日志中心收集分析
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout), // 容器化部署最佳实践：输出到标准输出，外层抓取
		logLevel,
	)

	// 全局注入 service 名称，区分不同微服务
	globalLogger = zap.New(core).With(
		zap.String("module", serviceName),
	)
}

// Ctx 携带 Context 的日志记录器，强制实施 TraceID 等追踪规范
// 业务代码中应该强制使用此方法：logger.Ctx(ctx).Info("订单已生成", zap.String("order_no", "123"))
func Ctx(ctx context.Context) *zap.Logger {
	if globalLogger == nil {
		Init("unknown-service", "info") // 兜底保护
	}

	logger := globalLogger

	// 从 Context 中提取并自动注入 TraceID
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		logger = logger.With(zap.String("trace_id", traceID))
	}
	// 提取 TenantID (SaaS 隔离基因)
	if tenantID, ok := ctx.Value(TenantIDKey).(string); ok && tenantID != "" {
		logger = logger.With(zap.String("tenant_id", tenantID))
	}
	// 提取 UserID
	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		logger = logger.With(zap.String("user_id", userID))
	}

	return logger
}

// 基础兼容方法
func Info(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Info(msg, fields...)
	}
}
func Error(msg string, fields ...zap.Field) {
	if globalLogger != nil {
		globalLogger.Error(msg, fields...)
	}
}

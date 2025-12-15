// internal/logger/logger.go
package logger

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	once     sync.Once
	instance *logrus.Logger
)

// Init 初始化日志（整个程序只调用一次，通常在 main.go 最前面）
func Init(env string) {
	once.Do(func() {
		instance = logrus.New()

		// 基础设置
		instance.SetOutput(os.Stdout)
		instance.SetReportCaller(true) // 显示文件名和行号（可选，生产可关）

		if env == "production" || env == "prod" {
			// 生产环境：JSON 格式，无颜色
			instance.SetFormatter(&logrus.JSONFormatter{
				PrettyPrint: false,
			})
			instance.SetLevel(logrus.WarnLevel) // 生产只看 Warn 以上
		} else {
			// 开发环境：超级好看的彩色文本
			instance.SetFormatter(&logrus.TextFormatter{
				ForceColors:            true,
				FullTimestamp:          true,
				TimestampFormat:        "2006-01-02 15:04:05",
				DisableLevelTruncation: true,
				PadLevelText:           true,
				QuoteEmptyFields:       true,
			})
			instance.SetLevel(logrus.DebugLevel) // 开发看全部
		}

		// 可选：添加 prefixed 让字段也带颜色（超漂亮）
		// instance.AddHook(prefixed.NewHook())
	})
}

func Println(args ...interface{}) {
	getLogger().Println(args...)
}
func Debug(args ...interface{}) {
	getLogger().Debug(args...)
}
func Info(args ...interface{}) {
	getLogger().Info(args...)
}
func Warn(args ...interface{}) {
	getLogger().Warn(args...)
}
func Error(args ...interface{}) {
	getLogger().Error(args...)
}
func Fatal(args ...interface{}) {
	getLogger().Fatal(args...)
}

// 带格式的
func Debugf(format string, args ...interface{}) {
	getLogger().Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	getLogger().Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	getLogger().Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	getLogger().Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	getLogger().Fatalf(format, args...)
}

// 带字段的（最推荐）
func WithField(key string, value interface{}) *logrus.Entry {
	return getLogger().WithField(key, value)
}
func WithFields(fields logrus.Fields) *logrus.Entry {
	return getLogger().WithFields(fields)
}

// 内部获取实例
func getLogger() *logrus.Logger {
	if instance == nil {
		// 没初始化就用默认开发模式
		Init("development")
	}
	return instance
}

// 暴露原始 logger（高级用法）
func Get() *logrus.Logger {
	return getLogger()
}

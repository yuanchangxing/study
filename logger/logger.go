// internal/logger/logger.go
package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	once     sync.Once
	instance *logrus.Logger

	// 当前 flags 配置
	flags           int
	timestampFormat = "2006-01-02 15:04:05"
)

// 模仿官方 log 包的 flags
const (
	Ldate         = 1 << iota // 日期: 2006/01/02
	Ltime                     // 时间: 15:04:05
	Lmicroseconds             // 微秒: 15:04:05.000000
	Llongfile                 // 完整文件路径
	Lshortfile                // 短文件名（覆盖 Llongfile）
	Lfuncname                 // 函数名: package/FuncName()
	LstdFlags     = Ldate | Ltime
)

// 级别颜色（只用于 [LEVEL ] 部分）
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[0;36m" // DEBUG
	colorGreen  = "\033[0;32m" // INFO（非粗体，更柔和）
	colorYellow = "\033[0;33m" // WARN
	colorRed    = "\033[0;31m" // ERROR/FATAL（全红，非粗体）
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b bytes.Buffer

	// 1. [LEVEL ] 带颜色
	levelStr := strings.ToUpper(entry.Level.String())
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = colorCyan
	case logrus.InfoLevel:
		levelColor = colorGreen
	case logrus.WarnLevel:
		levelColor = colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = colorRed
	default:
		levelColor = colorReset
	}

	b.WriteString(levelColor)
	b.WriteString(fmt.Sprintf("[%5s]", levelStr))
	b.WriteString(colorReset) // 颜色在这里结束，后续全部无色

	// 2. [时间]
	var ts string
	if flags&Lmicroseconds != 0 {
		ts = entry.Time.Format("2006-01-02 15:04:05.000000")
	} else if flags&(Ldate|Ltime) != 0 {
		if flags&Ldate != 0 && flags&Ltime != 0 {
			ts = entry.Time.Format(timestampFormat)
		} else if flags&Ldate != 0 {
			ts = entry.Time.Format("2006-01-02")
		} else if flags&Ltime != 0 {
			ts = entry.Time.Format("15:04:05")
		}
	}
	if ts != "" {
		b.WriteString(fmt.Sprintf(" [%s]", ts))
	}

	// 3. 调用者信息（无色）
	file := "???"
	line := 0
	funcName := "???"

	for i := 5; ; i++ {
		pc, fpath, l, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(fpath, "/internal/logger/") || strings.HasSuffix(fpath, "logger.go") {
			continue
		}
		file = fpath
		line = l
		if flags&Lfuncname != 0 && pc != 0 {
			full := runtime.FuncForPC(pc).Name()
			if idx := strings.LastIndex(full, "/"); idx != -1 {
				full = full[idx+1:]
			}
			if idx := strings.LastIndex(full, "."); idx != -1 {
				full = full[:idx] + "/" + full[idx+1:]
			}
			funcName = full
		}
		break
	}

	// 文件路径处理
	displayFile := file
	if flags&Lshortfile != 0 {
		displayFile = filepath.Base(file)
	} else if flags&Llongfile != 0 {
		if pwd := os.Getenv("PWD"); pwd != "" {
			if rel, err := filepath.Rel(pwd, file); err == nil && !strings.HasPrefix(rel, "..") {
				displayFile = filepath.ToSlash(rel)
			}
		}
	}

	b.WriteString(fmt.Sprintf(" %s:%d", displayFile, line))

	if flags&Lfuncname != 0 {
		b.WriteString(fmt.Sprintf(" %s()", funcName))
	}

	// 4. 消息内容（无色）
	b.WriteString(" " + entry.Message)

	// 换行
	b.WriteString("\n")

	return b.Bytes(), nil
}

func Init(env string, flagList ...int) {
	once.Do(func() {
		instance = logrus.New()
		instance.SetOutput(os.Stdout)
		instance.SetReportCaller(true)
		instance.SetFormatter(&CustomFormatter{})

		// 计算 flags
		flags = 0
		for _, f := range flagList {
			flags |= f
		}
		if flags == 0 {
			flags = LstdFlags
		}

		if flags&Lmicroseconds != 0 {
			timestampFormat = "2006-01-02 15:04:05.000000"
		}

		if env == "production" || env == "prod" {
			instance.SetLevel(logrus.WarnLevel)
		} else {
			instance.SetLevel(logrus.DebugLevel)
		}
	})
}

// 封装函数不变
func Println(args ...interface{}) { getLogger().Println(args...) }
func Debug(args ...interface{})   { getLogger().Debug(args...) }
func Info(args ...interface{})    { getLogger().Info(args...) }
func Warn(args ...interface{})    { getLogger().Warn(args...) }
func Error(args ...interface{})   { getLogger().Error(args...) }
func Fatal(args ...interface{})   { getLogger().Fatal(args...) }

func Debugf(format string, args ...interface{}) { getLogger().Debugf(format, args...) }
func Infof(format string, args ...interface{})  { getLogger().Infof(format, args...) }
func Warnf(format string, args ...interface{})  { getLogger().Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { getLogger().Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { getLogger().Fatalf(format, args...) }

func WithField(key string, value interface{}) *logrus.Entry {
	return getLogger().WithField(key, value)
}
func WithFields(fields logrus.Fields) *logrus.Entry {
	return getLogger().WithFields(fields)
}

func getLogger() *logrus.Logger {
	if instance == nil {
		Init("development", LstdFlags)
	}
	return instance
}

func Get() *logrus.Logger {
	return getLogger()
}

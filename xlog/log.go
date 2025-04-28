package xlog

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var Logger *logrus.Logger

func SwitchColor(open bool) {
	f := Logger.Formatter.(*CustomTextFormatter)
	ff := func(colo *color.Color) {
		colo.EnableColor()
	}
	if !open {
		ff = func(colo *color.Color) {
			colo.DisableColor()
		}
	}

	ff(f.ColorError)
	ff(f.ColorWarning)
	ff(f.ColorDebug)
	ff(f.ColorFatal)
	ff(f.ColorInfo)
}

func init() {
	Logger = logrus.New()
	var cus = &CustomTextFormatter{
		ForceColors:  true,
		ColorInfo:    color.New(color.FgWhite),
		ColorWarning: color.New(color.FgYellow),
		ColorError:   color.New(color.FgRed),
		ColorDebug:   color.New(color.FgBlue),
		ColorFatal:   color.New(color.FgRed),
	}
	Logger.Formatter = cus
	Logger.Level = logrus.DebugLevel

	Logger.ReportCaller = true

}

type CustomTextFormatter struct {
	logrus.TextFormatter
	ForceColors  bool
	ColorWarning *color.Color
	ColorError   *color.Color
	ColorDebug   *color.Color
	ColorInfo    *color.Color
	ColorFatal   *color.Color
	flag         int32
}

func (f *CustomTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.ForceColors {
		var s = formatMessageWithTimestamp(entry)
		switch entry.Level {
		case logrus.DebugLevel:
			f.ColorDebug.Println(s)
		case logrus.InfoLevel:
			f.ColorInfo.Println(s)
		case logrus.WarnLevel:
			f.ColorWarning.Println(s)
		case logrus.ErrorLevel:
			f.ColorError.Println(s)
		case logrus.FatalLevel:
			f.ColorFatal.Println(s)
		default:
			f.ColorInfo.Println(s)
		}
		return nil, nil
	} else {
		// 否则，返回默认格式化输出
		return f.formatDefault(entry)
	}
}

// 格式化消息并添加时间戳
func formatMessageWithTimestamp(entry *logrus.Entry) string {
	timestamp := entry.Time.Format(time.DateTime)
	var line string
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		line = fmt.Sprintf(" %s:%-3d", fName, entry.Caller.Line)
	}
	return "[" + levelString(entry.Level) + "]" + timestamp + line + ": " + entry.Message
}

func levelString(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "D"
	case logrus.InfoLevel:
		return "I"
	case logrus.WarnLevel:
		return "W"
	case logrus.ErrorLevel:
		return "E"
	default:
		return level.String()
	}
}

// 格式化默认输出的方法
func (f *CustomTextFormatter) formatDefault(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(time.DateTime)
	level := entry.Level.String()
	msg := entry.Message

	return []byte(fmt.Sprintf("[%s] %s %s\n", level, timestamp, msg)), nil
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}
func Warnf(farmat string, args ...interface{}) {
	Logger.Warnf(farmat, args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

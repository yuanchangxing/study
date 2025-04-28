package xlog

import (
	"testing"
)

func TestCustomTextFormatter_PrintColored(t *testing.T) {
	SwitchColor(true)
	Logger.Infof("This is a %s test", "test")
	Logger.Debug("This is a debug test", "test")
	Logger.Errorf("This is a %s test", "test")
	Logger.Warnf("This is a %s test", "test")
	Infof("This is a %s test", "test")
	Debugf("This is a %s test", "test")
	Warnf("This is a %s test", "test")
	Errorf("This is a %s test", "test")
}

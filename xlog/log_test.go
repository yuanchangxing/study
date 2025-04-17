package xlog

import (
	"testing"
)

func TestCustomTextFormatter_PrintColored(t *testing.T) {
	SwitchColor(true)
	logger.Infof("This is a %s test", "test")
	logger.Debug("This is a debug test", "test")
	logger.Errorf("This is a %s test", "test")
	logger.Warnf("This is a %s test", "test")
}

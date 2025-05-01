package tools

import (
	"github.com/yuanchangxing/study/xlog"
	"testing"
)

func TestNewRateLimit(t *testing.T) {
	var n = NewRateLimit(1, 2)
	for i := 0; i < 100; i++ {
		go func() {
			for i := 0; i < 10; i++ {
				n.Wait()
				xlog.Info("do something")
			}
		}()

	}
	<-chan bool(nil)
}

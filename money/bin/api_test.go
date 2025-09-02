package bin

import (
	"github.com/yuanchangxing/study/xlog"
	"testing"
)

func TestDoSom(t *testing.T) {
	Init(mockConf{})
	a, err := WalletBalance()
	xlog.PanicCheckErr(err)
	for _, v := range a {
		xlog.Info(v.Name)
	}
}

type mockConf struct{}

func (m mockConf) GetApiKey() string {
	return "jtEGSkmML0KSgzaPB51EbTWi7GFNtkSUwpQx8M8emu7R8pMPyrd8bTR52sN2Pc5n"
}

func (m mockConf) GetApiSecret() string {
	return "tN9oq3eMnlWx8tm3hW7EBep1KFM8COlDbHWWiueA5sG5BFlrpcKzWNOVDqoawMqf"
}

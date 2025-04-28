package algorithm

import (
	"fmt"
	"github.com/yuanchangxing/study/xlog"
	"testing"
)

func TestInsert(t *testing.T) {
	rank := NewRank(2)
	rank.Push(&xx{Data: 1})
	rank.Push(&xx{Data: 2})
	rank.Push(&xx{Data: 3})
	rank.Push(&xx{Data: 4})
	rank.Del(&xx{Data: 4})
	rank.Push(&xx{Data: 1})
	rank.Push(&xx{Data: 9})
	rank.Push(&xx{Data: 9})
	rank.Push(&xx{Data: 9})
	rank.Push(&xx{Data: 10})
	rank.Del(&xx{Data: 33})
	for data := rank.Pop(); data != nil; data = rank.Pop() {
		xlog.Logger.Info(data)
	}
}

func TestAl(t *testing.T) {
	var heap maxHeap
	heap.dMap = make(map[string]int)
	heap.Push(&xx{Data: 0})

	heap.Push(&xx{Data: 1})

	heap.Push(&xx{Data: 2})
	heap.Push(&xx{Data: 2})
	heap.Push(&xx{Data: 9})
	heap.Push(&xx{Data: 3})
	heap.Del("3")
	heap.Del("2")
	heap.Push(&xx{Data: -1})
	xlog.SwitchColor(true)
	for data := heap.Pop(); data != nil; data = heap.Pop() {
		xlog.Logger.Info(data)
	}
}

type xx struct {
	Data int
}

func (x *xx) Id() string {
	return fmt.Sprintf("%d", x.Data)
}

func (x *xx) Less(data maxHeapInterface) bool {
	return x.Data < data.(*xx).Data
}

func (x *xx) String() string {
	return fmt.Sprintf("{%v}", x.Data)
}

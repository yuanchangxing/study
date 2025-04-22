package algorithm

import (
	"fmt"
	"log"
	"testing"
)

func TestAl(t *testing.T) {
	var heap maxHeap
	heap.Push(&xx{Data: 0})

	heap.Push(&xx{Data: 1})

	heap.Push(&xx{Data: 2})
	heap.Push(&xx{Data: 2})
	heap.Modify(3, &xx{Data: 9})
	heap.Push(&xx{Data: 3})
	heap.Del(3)
	heap.Push(&xx{Data: -1})

	for data := heap.Pop(); data != nil; data = heap.Pop() {
		log.Println(data)
	}
}

type xx struct {
	Data int
}

func (x *xx) Less(data maxHeapInterface) bool {
	return x.Data < data.(*xx).Data
}

func (x *xx) String() string {
	return fmt.Sprintf("{%v}", x.Data)
}

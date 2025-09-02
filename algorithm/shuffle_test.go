package algorithm

import (
	"fmt"
	"testing"
	"time"
)

func TestShuffle(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	Init(time.Now().UnixMicro())
	Shuffle(a)
	fmt.Println(a)
}

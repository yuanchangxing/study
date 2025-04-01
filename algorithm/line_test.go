package algorithm

import (
	"log"
	"testing"
)

func TestLineRe(t *testing.T) {

	var n = NewLNode([]int{1, 2, 3, 4, 5, 6, 7, 8})
	log.Println(n)
	ReverseN(&n)
	log.Println(n)

	//x := n.Reverse2()
	//log.Println(x)
}

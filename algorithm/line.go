package algorithm

import (
	"fmt"
	"strconv"
	"strings"
)

// 单链表逆转
type LNode struct {
	next *LNode
	val  int
}

func NewLNode(val []int) *LNode {
	if len(val) == 0 {
		return nil
	}

	var node = &LNode{val: val[0]}
	var root = node
	for i := 1; i < len(val); i++ {
		newNode := &LNode{val: val[i], next: nil}
		node.next = newNode
		node = newNode
	}

	return root
}

func (l *LNode) String() string {
	var data = make([]string, 0, 16)
	var max = 2048
	for l != nil && max > 0 {
		data = append(data, strconv.Itoa(l.val))
		l = l.next
		max--
	}
	return fmt.Sprintf("(%s)", strings.Join(data, ","))
}

// 逆转
func (l *LNode) Reverse() *LNode {
	if l == nil {
		return nil
	}
	var current = l
	var next = current.next
	var before *LNode
	for current != nil {
		next = current.next
		current.next = before
		before = current
		current = next
	}

	/*
		 *l = *before
		 这里的修改是有问题的，这里这个节点的地址并没有发生变化，但是修改了它的值，
		而这个节点有使用在上面，因为地址没有变化，可能导致上面已经指向这个地址的指针的值发生变化。
	*/
	return before
}

// 不用return，我使用二级指针
func ReverseN(m **LNode) {
	if m == nil {
		return
	}
	var l = *m
	if l == nil {
		return
	}

	var current = l
	var next = current.next
	var before *LNode
	for current != nil {
		next = current.next
		current.next = before
		before = current
		current = next
	}

	// 我这里想把这个头节点赋值为这个新的头节点，我使用多级指针
	*m = before
}

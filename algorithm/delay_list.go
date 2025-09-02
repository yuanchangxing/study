package algorithm

import (
	"context"
	"sync/atomic"
	"time"
)

type delayList struct {
	list   *Heap // 单链表
	re     chan interface{}
	cancel context.CancelFunc
	fla    int32
}

type delayData struct {
	data interface{}
	t    int64
}

func (d *delayData) Less(h heapInterface) bool {
	return d.t > h.(*delayData).t
}

func (d *delayList) Push(v interface{}, timeout int64) {
	d.list.Push(&delayData{data: v, t: timeout})
}

func (d *delayList) loop() {
	if !atomic.CompareAndSwapInt32(&d.fla, 0, 1) {
		return
	}

	var ticker = time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	c, can := context.WithCancel(context.TODO())
	d.cancel = can
	for {
		select {
		case <-ticker.C:
		case <-c.Done():
			return
		}
		now := time.Now().Unix()
		res := d.list.RangAndDel(func(data heapInterface) bool {
			return data.(*delayData).t <= now
		})

		for _, v := range res {
			d.re <- v.(*delayData).data
		}

	}
}

func NewDelayList(reChanBufLen int) *delayList {
	var d = &delayList{}
	d.list = &Heap{}
	d.re = make(chan interface{}, reChanBufLen) //todo
	//d.check = make(chan struct{}, 1)
	go d.loop()
	return d
}

func (d *delayList) PopWait() interface{} {
	res := <-d.re
	return res
}

func (d *delayList) Stop() {
	if !atomic.CompareAndSwapInt32(&d.fla, 1, 2) {
		return
	}

	d.cancel()
}

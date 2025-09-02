package algorithm

import (
	"time"
)

type delayList struct {
	list *Heap // 单链表
	//timeout int64
	//mutex sync.Mutex
	re chan interface{}
	//done    chan struct{}
	//check chan struct{}
	id int64
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
	//select {
	//case d.check <- struct{}{}:
	//default:
	//}
}

func (d *delayList) loop() {
	var ticker = time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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

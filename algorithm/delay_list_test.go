package algorithm

import (
	"log"
	"testing"
	"time"
)

func TestNewDelayList(t *testing.T) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	d := NewDelayList(16)
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Duration(i) * time.Second)
			d.Push(i, time.Now().Unix()+5)
			log.Println("push ok", i)
		}(i)
	}
	d.Push(11, time.Now().Unix()+1)
	for {
		data := d.PopWait()

		log.Println("pop", data)
	}
}

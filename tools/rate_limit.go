package tools

import (
	"sync"
	"time"
)

type rateLimit struct {
	seconds    int
	maxNum     int
	recordTime int
	recordNum  int
	mutex      sync.Mutex
}

func (r *rateLimit) SetNum(num int) {
	if num < 0 {
		num = 0
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.recordNum = num
}

func NewRateLimit(seconds int, maxNum int) *rateLimit {
	if seconds < 0 {
		seconds = 0
	}
	if maxNum < 0 {
		maxNum = 0
	}
	return &rateLimit{seconds: seconds, maxNum: maxNum}
}

func (rl *rateLimit) Wait() {
	if rl.seconds <= 0 {
		return
	}

	for {
		var now = time.Now()
		rl.mutex.Lock()
		if now.Unix() == int64(rl.recordTime) {
			if rl.recordNum >= rl.maxNum {
				rl.mutex.Unlock()
				time.Sleep(time.Unix(int64(rl.recordTime+rl.seconds), 0).Sub(now))
				continue
			}

			rl.recordNum++
			rl.mutex.Unlock()
			break
		}

		rl.recordTime = int(now.Unix())
		rl.recordNum = 1
		rl.mutex.Unlock()
		break
	}
}

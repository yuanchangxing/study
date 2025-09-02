package algorithm

import (
	"sync"
)

type Heap struct {
	data  []heapInterface
	mutex sync.RWMutex
}

func NewHeap() *Heap {
	return &Heap{}
}

func (h *Heap) Top() heapInterface {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.data[0]
}

func (h *Heap) RangAndDel(f func(data heapInterface) bool) []heapInterface {
	var res []heapInterface
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for len(h.data) > 0 {
		re := h.topLocked()
		if f(re) {
			res = append(res, re)
			h.popLocked()
			continue
		}
		break
	}

	return res
}

func (h *Heap) topLocked() heapInterface {
	if len(h.data) == 0 {
		return nil
	}
	return h.data[0]
}

func (h *Heap) popLocked() heapInterface {
	if len(h.data) == 0 {
		return nil
	}

	res := h.data[0]
	h.data = h.data[1:]
	h.downLocked(0)

	return res
}

func (h *Heap) Pop() heapInterface {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.popLocked()
}

func (h *Heap) Push(data heapInterface) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.data = append(h.data, data)
	if len(h.data) == 1 {
		return
	}

	h.upLocked(len(h.data) - 1)

}

func (h *Heap) upLocked(ind int) {
	for ind > 0 {
		// 查找父节点
		father := (ind - 1) / 2
		if !h.compareAndSwapLocked(ind, father, true) {
			break
		}
		ind = father
	}
}

func (h *Heap) downLocked(ind int) {
	for ind < len(h.data)-1 {
		left := ind*2 + 1 // left
		if h.compareAndSwapLocked(ind, left, false) {
			ind = left
			continue
		}

		if h.compareAndSwapLocked(ind, left+1, false) {
			ind = left + 1
			continue
		}
		break
	}
}

func (h *Heap) compareAndSwapLocked(i, j int, isUp bool) bool {
	size := len(h.data) - 1
	if size < i || size < j {
		return false
	}

	if h.data[i] == nil || h.data[j] == nil {
		return false
	}

	if isUp {
		if !h.data[i].Less(h.data[j]) {
			h.data[i], h.data[j] = h.data[j], h.data[i]
			return true
		}
	} else {
		if h.data[i].Less(h.data[j]) {
			h.data[i], h.data[j] = h.data[j], h.data[i]
			return true
		}
	}

	return false
}

type rankHeapInterface interface {
	heapInterface
	Id() string
}

type heapInterface interface {
	Less(h heapInterface) bool
}

type maxHeap struct {
	data  []rankHeapInterface
	dMap  map[string]int
	less  func(a, b rankHeapInterface) bool
	mutex sync.RWMutex
}

func (m *maxHeap) GetSize() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.data)
}

func (m *maxHeap) Del(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	ind, ok := m.dMap[id]
	if !ok {
		return
	}

	ordData := m.data[ind]
	finalData := m.data[len(m.data)-1]
	delete(m.dMap, id)
	m.data = m.data[:len(m.data)-1]
	if len(m.data) == 0 || finalData == ordData {
		return
	}

	m.modifyLocked(ind, finalData, ordData)
}

func (m *maxHeap) Push(data rankHeapInterface) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.data) == 0 {
		m.data = append(m.data, data)
		m.dMap[data.Id()] = len(m.data) - 1
		return true
	}

	ordIndex, ok := m.dMap[data.Id()]
	if ok {
		m.modifyLocked(ordIndex, data, m.data[ordIndex])
		return false
	}

	m.data = append(m.data, data)
	m.dMap[data.Id()] = len(m.data) - 1
	ind := len(m.data) - 1
	m.upLocked(ind)
	return true
}

func (m *maxHeap) modifyLocked(ind int, data rankHeapInterface, ordData rankHeapInterface) {
	m.data[ind] = data
	m.dMap[data.Id()] = ind
	if m.less(ordData, data) {
		m.upLocked(ind)
	} else {
		m.downLocked(ind)
	}
}

func (m *maxHeap) Pop() rankHeapInterface {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.data) == 0 {
		return nil
	}

	res := m.data[0]
	m.data[0] = m.data[len(m.data)-1]
	m.data = m.data[:len(m.data)-1]
	delete(m.dMap, res.Id())
	if len(m.data) != 0 {
		m.downLocked(0)
	}
	return res

}

func (m *maxHeap) upLocked(ind int) {
	for ind > 0 {
		// 查找父节点
		father := (ind - 1) / 2
		if !m.compareAndSwapLocked(ind, father, true) {
			break
		}
		ind = father
	}
}

func (m *maxHeap) downLocked(ind int) {
	for ind < len(m.data)-1 {
		left := ind*2 + 1 // left
		if m.compareAndSwapLocked(ind, left, false) {
			ind = left
			continue
		}

		if m.compareAndSwapLocked(ind, left+1, false) {
			ind = left + 1
			continue
		}
		break

	}
}

func (m *maxHeap) compareAndSwapLocked(i, j int, isUp bool) bool {
	size := len(m.data) - 1
	if size < i || size < j {
		return false
	}

	if m.data[i] == nil || m.data[j] == nil {
		return false
	}

	if isUp {
		if !m.less(m.data[i], m.data[j]) {
			m.data[i], m.data[j] = m.data[j], m.data[i]
			m.dMap[m.data[i].Id()] = i
			m.dMap[m.data[j].Id()] = j
			return true
		}
	} else {
		if m.less(m.data[i], m.data[j]) {
			m.data[j], m.data[i] = m.data[i], m.data[j]
			m.dMap[m.data[i].Id()] = i
			m.dMap[m.data[j].Id()] = j
			return true
		}
	}

	return false
}

type Rank struct {
	maxHeap *maxHeap
	minHeap *maxHeap
	mu      sync.Mutex
	size    int
}

func NewMaxHeap(f func(a, b rankHeapInterface) bool) *maxHeap {
	return &maxHeap{dMap: make(map[string]int), less: f}
}

func NewRank(size int) *Rank {
	return &Rank{
		maxHeap: &maxHeap{dMap: make(map[string]int), less: func(a, b rankHeapInterface) bool {
			return a.Less(b)
		}},
		minHeap: &maxHeap{dMap: make(map[string]int), less: func(a, b rankHeapInterface) bool {
			return !a.Less(b)
		}},
		size: size,
	}
}

func (r *maxHeap) Top() rankHeapInterface {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.data) == 0 {
		return nil
	}
	return r.data[0]
}

func (r *Rank) Pop() rankHeapInterface {
	da := r.maxHeap.Pop()
	if da == nil {
		return nil
	}
	r.minHeap.Del(da.Id())
	return da
}

func (r *Rank) Del(data rankHeapInterface) {
	r.maxHeap.Del(data.Id())
	r.minHeap.Del(data.Id())
}

func (r *Rank) Push(data rankHeapInterface) {
	r.mu.Lock()
	defer r.mu.Unlock()

	isInsert := r.maxHeap.Push(data)
	r.minHeap.Push(data)
	if r.size == 0 || !isInsert {
		return
	}

	if r.size >= r.maxHeap.GetSize() {
		return
	}
	de := r.minHeap.Pop()
	r.maxHeap.Del(de.Id())
}

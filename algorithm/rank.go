package algorithm

import (
	"sync"
)

type maxHeapInterface interface {
	Less(h maxHeapInterface) bool
	Id() string
}
type maxHeap struct {
	data  []maxHeapInterface
	dMap  map[string]int
	less  func(a, b maxHeapInterface) bool
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

func (m *maxHeap) Push(data maxHeapInterface) bool {
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

func (m *maxHeap) modifyLocked(ind int, data maxHeapInterface, ordData maxHeapInterface) {
	m.data[ind] = data
	m.dMap[data.Id()] = ind
	if m.less(ordData, data) {
		m.upLocked(ind)
	} else {
		m.downLocked(ind)
	}
}

func (m *maxHeap) Pop() maxHeapInterface {
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

func NewRank(size int) *Rank {
	return &Rank{
		maxHeap: &maxHeap{dMap: make(map[string]int), less: func(a, b maxHeapInterface) bool {
			return a.Less(b)
		}},
		minHeap: &maxHeap{dMap: make(map[string]int), less: func(a, b maxHeapInterface) bool {
			return !a.Less(b)
		}},
		size: size,
	}
}

func (r *maxHeap) Top() maxHeapInterface {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.data) == 0 {
		return nil
	}
	return r.data[0]
}

func (r *Rank) Pop() maxHeapInterface {
	da := r.maxHeap.Pop()
	if da == nil {
		return nil
	}
	r.minHeap.Del(da.Id())
	return da
}

func (r *Rank) Del(data maxHeapInterface) {
	r.maxHeap.Del(data.Id())
	r.minHeap.Del(data.Id())
}

func (r *Rank) Push(data maxHeapInterface) {
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

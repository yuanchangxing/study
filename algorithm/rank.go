package algorithm

import (
	"sync"
)

type maxHeapInterface interface {
	Less(h maxHeapInterface) bool
}
type maxHeap struct {
	data  []maxHeapInterface
	size  int
	mutex sync.RWMutex
}

func (m *maxHeap) Del(ind int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.data) == 1 {
		m.data = m.data[:len(m.data)-1]
		return
	}

	if len(m.data) <= ind {
		return
	}

	m.data[ind] = m.data[len(m.data)-1]
	m.data = m.data[:len(m.data)-1]
	m.modifyLocked(ind, m.data[ind])
}

func (m *maxHeap) Push(data maxHeapInterface) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.data) == 0 {
		m.data = append(m.data, data)
		return
	}

	m.data = append(m.data, data)
	ind := len(m.data) - 1
	m.upLocked(ind)
}

func (m *maxHeap) modifyLocked(ind int, data maxHeapInterface) {
	ordData := m.data[ind]
	m.data[ind] = data
	if ordData.Less(data) {
		m.upLocked(ind)
	} else {
		m.downLocked(ind)
	}
}

func (m *maxHeap) Modify(ind int, data maxHeapInterface) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.modifyLocked(ind, data)
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
	m.downLocked(0)
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
		if !m.data[i].Less(m.data[j]) {
			m.data[i], m.data[j] = m.data[j], m.data[i]
			return true
		}
	} else {
		if m.data[i].Less(m.data[j]) {
			m.data[j], m.data[i] = m.data[i], m.data[j]
			return true
		}
	}

	return false
}

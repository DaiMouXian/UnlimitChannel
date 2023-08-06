package UnlimitChannel

import (
	"container/list"
	"sync"
)

type Queue struct {
	data  *list.List
	mutex sync.RWMutex
}

func NewQueue() *Queue {
	q := &Queue{data: list.New(), mutex: sync.RWMutex{}}
	return q
}

func (q *Queue) Push(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.data.PushBack(v)
}

func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.data.Len() > 0 {
		element := q.data.Front()
		v := element.Value
		q.data.Remove(element)
		return v
	} else {
		return nil
	}
}

func (q *Queue) Top() interface{} {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if q.data.Len() > 0 {
		element := q.data.Front()
		v := element.Value
		return v
	} else {
		return nil
	}
}

func (q *Queue) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return q.data.Len()
}

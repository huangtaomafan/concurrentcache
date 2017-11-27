package concurrentcache

import (
	"sync"
	"errors"
)

type Queue struct {
	head *QNode
	tail *QNode
	pool *sync.Pool
}

type QNode struct {
	prev *QNode
	next *QNode
	v interface{}
}

func NewQueue() *Queue {
	pool := &sync.Pool{
		New: func() interface{} {
			return &QNode{}
		},
	}
	return &Queue{pool:pool}
}

func (q *Queue) NewQNode(v interface{}) *QNode {
	n := q.pool.Get().(*QNode)
	n.reset()
	n.v = v
	return n
}

func (q *Queue) Recycle(n *QNode) {
	if n != nil {
		q.pool.Put(n)
	}
}

func (n *QNode) reset() {
	if n != nil {
		n.prev = nil
		n.next = nil
		n.v = nil
	}
}

func (q *Queue) LPush(n *QNode) {
	q.push(n, true)
}

func (q *Queue) RPush(n *QNode) {
	q.push(n, false)
}

func (q *Queue) push(n *QNode, left bool) {
	if q.Empty() {
		n.next, n.prev = nil, nil
		q.head, q.tail = n, n
		return
	} else {
		if left {
			n.next, n.prev = q.head, nil
			q.head.prev = n
			q.head = n
		} else {
			n.next, n.prev = nil, q.tail
			q.tail.next = n
			q.tail = n
		}
	}
}

func (q *Queue) Insert(prev, n *QNode) (bool, error) {
	return true,nil
}

func (q *Queue) LPop() *QNode {
	return q.pop(true)
}

func (q *Queue) RPop() *QNode {
	return q.pop(false)
}

func (q *Queue) pop(left bool) *QNode {
	if q.Empty() {
		return nil
	}
	if left {
		n := q.head
		if q.head == q.tail {
			q.head, q.tail = nil, nil
		} else {
			q.head = q.head.next
		}
		n.next, n.prev = nil, nil
		return n
	} else {
		n := q.tail
		if q.head == q.tail {
			q.head, q.tail = nil, nil
		} else {
			q.tail = q.tail.prev
		}
		n.next, n.prev = nil, nil
		return n
	}
}

func (q *Queue) Delete(n *QNode) (bool, error) {
	if q.Empty() || n == nil {
		return false, errors.New("invalid queue node")
	}
	if n.prev == nil && n != q.head {
		return false, errors.New("invalid queue node")
	} else if n.next == nil && n != q.tail {
		return false, errors.New("invalid queue node")
	}
	if n.prev != nil {
		n.prev.next = n.next
	} else if n == q.head && n != nil {
		q.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else if n == q.tail && n != nil {
		q.tail = n.prev
	}
	n.prev = nil
	n.next = nil
	return true, nil
}

func (q *Queue) Empty() bool {
	if q.head == q.tail && q.head == nil {
		return true
	}
	return false
}
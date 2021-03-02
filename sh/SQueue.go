package sh

import "sync"

type SQueue struct {
    items    []int
    capacity int
    tail     int
    head     int
    size     int
    access   sync.Mutex
    hasItems *sync.Cond
    hasRoom  *sync.Cond
}

func NewSQueue(capacity int) (sq *SQueue) {
    items := make([]int, capacity)
    sq = &SQueue{items:items, capacity:capacity}
    sq.hasRoom = sync.NewCond(&sq.access)
    sq.hasItems = sync.NewCond(&sq.access)
    return sq
}

func (sq *SQueue) Get(item *int) {
    sq.access.Lock()
    if sq.size == 0 {
        sq.hasItems.Wait()
    }
    (*item) = sq.items[sq.head]
    sq.head = (sq.head + 1) % sq.capacity
    sq.size--
    sq.hasRoom.Signal()
    sq.access.Unlock()
}

func (sq *SQueue) Add(item int) {
    sq.access.Lock()
    if sq.size == sq.capacity {
        sq.hasRoom.Wait()
    }
    sq.items[sq.tail] = item
    sq.tail = (sq.tail + 1) % sq.capacity
    sq.size++
    sq.hasItems.Signal()
    sq.access.Unlock()
}
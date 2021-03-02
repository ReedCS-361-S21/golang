package sh

import "sync"

type STally struct {
    sum     int
    done    int
    size    int
    lock    sync.Mutex
    allDone *sync.Cond
}

func NewSTally(groupSize int) (tally *STally) {
    tally = &STally{size:groupSize}
    tally.allDone = sync.NewCond(&tally.lock)
    return tally
}

func (tally *STally) Submit(value int) {
  tally.lock.Lock()
  tally.sum += value
  tally.done++
  if tally.done == tally.size { 
    tally.allDone.Broadcast()
  }
  tally.lock.Unlock()
}

func (tally *STally) Get(sum *int) {
  tally.lock.Lock()
  for tally.done < tally.size {
    tally.allDone.Wait()
  }
  (*sum) = tally.sum
  tally.lock.Unlock()
}

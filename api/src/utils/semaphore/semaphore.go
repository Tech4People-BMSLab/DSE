package semaphore

import (
	"sync"
)

// ------------------------------------------------------------
// : Semaphore
// ------------------------------------------------------------
type Semaphore struct {
    ch 		  chan struct{}
    mu        sync.Mutex
    count     int
    limit     int
}

// ------------------------------------------------------------
// : Constructor
// ------------------------------------------------------------
func New(limit int) *Semaphore {
    return &Semaphore{
        ch: make(chan struct{}, limit),
        limit:     limit,
    }
}

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func (s *Semaphore) Acquire() {
    s.ch <- struct{}{}
    s.mu.Lock()
    s.count++
    s.mu.Unlock()
}

func (s *Semaphore) Release() {
    <-s.ch
    s.mu.Lock()
    s.count--
    s.mu.Unlock()
}

func (s *Semaphore) GetCount() int {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.count
}

func (s *Semaphore) GetLimit() int {
    return s.limit
}

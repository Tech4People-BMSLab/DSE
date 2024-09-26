package gatekeeper

import (
	"sync"
)

// GateKeeper controls access to a resource or section of code among multiple goroutines.
type GateKeeper struct {
	mutex          sync.Mutex
	cond           *sync.Cond
	counter        int
	is_open        bool
}

// NewGateKeeper initializes a new GateKeeper. If `locked` is true, the gate starts in a locked state.
func NewGateKeeper(locked bool) *GateKeeper {
	gk := &GateKeeper{is_open: !locked}
	gk.cond = sync.NewCond(&gk.mutex)
	return gk
}

// Lock sets the gate to a locked state, preventing goroutines from passing until it is unlocked.
func (gk *GateKeeper) Lock() {
	gk.mutex.Lock()
	gk.is_open = false
	gk.mutex.Unlock()
}

// Unlock sets the gate to an open state, allowing all waiting goroutines to proceed.
func (gk *GateKeeper) Unlock() {
	gk.mutex.Lock()
	gk.is_open = true
	gk.cond.Broadcast()
	gk.mutex.Unlock()
}

// UnlockOne allows exactly one waiting goroutine to proceed, even if the gate is generally closed.
// It prioritizes one goroutine if multiple are waiting.
func (gk *GateKeeper) UnlockOne() {
	gk.mutex.Lock()
	defer gk.mutex.Unlock()
	if gk.counter > 0 {
		gk.counter--
		gk.cond.Signal()
	} else {
		gk.is_open = true
		gk.cond.Broadcast()
	}
}

// AllowIf lets a goroutine pass through the gate only if a specific condition is true.
// The condition is defined by the predicate function provided as an argument.
// If the gate is open, the predicate is ignored and the goroutine is allowed to proceed.
func (gk *GateKeeper) AllowIf(predicate func() bool) {
	gk.mutex.Lock()
	defer gk.mutex.Unlock()

	for !gk.is_open && !predicate() {
		gk.cond.Wait()
	}
}

// Wait blocks the calling goroutine until the gate is fully opened.
// It is useful when a goroutine needs to wait indefinitely until unrestricted access is allowed.
func (gk *GateKeeper) Wait() {
	gk.mutex.Lock()
	defer gk.mutex.Unlock()
	for !gk.is_open {
		gk.cond.Wait()
	}
}

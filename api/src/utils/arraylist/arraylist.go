package arraylist

import (
	"encoding/json"
	"sync"
)

// ArrayList is a generic list structure
type ArrayList[T comparable] struct {
    elements []T
    mu       sync.RWMutex
}

// NewArrayList creates a new instance of ArrayList
func NewArrayList[T comparable]() *ArrayList[T] {
    return &ArrayList[T]{elements: make([]T, 0)}
}

// Add appends an element to the list
func (a *ArrayList[T]) Add(element T) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.elements = append(a.elements, element)
}

// Get retrieves the element at the given index
func (a *ArrayList[T]) Get(index int) (T, bool) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if index >= 0 && index < len(a.elements) {
        return a.elements[index], true
    }
    var zeroValue T
    return zeroValue, false
}

// MustGet retrieves the element at the given index.
// It panics if the index is out of bounds.
func (a *ArrayList[T]) MustGet(index int) T {
    a.mu.RLock()
    defer a.mu.RUnlock()
    if index >= 0 && index < len(a.elements) {
        return a.elements[index]
    }
    panic("Index out of bounds in ArrayList")
}

// Remove removes the element at the given index
func (a *ArrayList[T]) Remove(index int) bool {
    a.mu.Lock()
    defer a.mu.Unlock()
    if index >= 0 && index < len(a.elements) {
        a.elements = append(a.elements[:index], a.elements[index+1:]...)
        return true
    }
    return false
}

// MustDelete removes the element at the given index.
// It panics if the index is out of bounds.
func (a *ArrayList[T]) MustDelete(index int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    if index >= 0 && index < len(a.elements) {
        a.elements = append(a.elements[:index], a.elements[index+1:]...)
    } else {
        panic("Index out of bounds in ArrayList")
    }
}

// Has checks if the element exists in the list
func (a *ArrayList[T]) Has(element T) bool {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for _, e := range a.elements {
        if e == element {
            return true
        }
    }
    return false
}

// Len returns the number of elements in the list
func (a *ArrayList[T]) Len() int {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return len(a.elements)
}

// Each applies a function to each element in the list.
// If the function returns false, the iteration stops.
func (a *ArrayList[T]) Each(f func(i int, v T) bool) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    for i, e := range a.elements {
        if !f(i, e) {
            break
        }
    }
}

// Filter returns a slice of elements that satisfy the given predicate, along with their indices.
func (a *ArrayList[T]) Filter(f func(element T, index int) bool) ([]T, []int) {
    a.mu.RLock()
    defer a.mu.RUnlock()

    var filtered []T
    var indices []int
    for i, e := range a.elements {
        if f(e, i) {
            filtered = append(filtered, e)
            indices = append(indices, i)
        }
    }
    return filtered, indices
}

// ArrayListIterator is an iterator for ArrayList
type ArrayListIterator[T comparable] struct {
    list  *ArrayList[T]
    index int
}

// Iterate returns a new iterator for the list
func (a *ArrayList[T]) Iterate() *ArrayListIterator[T] {
    a.mu.RLock()
    return &ArrayListIterator[T]{list: a, index: 0}
}

// Next moves the iterator to the next element and returns it
func (it *ArrayListIterator[T]) Next() (T, bool) {
    a := it.list
    if it.index < len(a.elements) {
        value := a.elements[it.index]
        it.index++
        if it.index == len(a.elements) {
            a.mu.RUnlock() // Release lock when iteration is done
        }
        return value, true
    }
    a.mu.RUnlock() // Release lock if iteration was incomplete
    var zeroValue T
    return zeroValue, false
}

// ToJSON serializes the ArrayList to JSON
func (a *ArrayList[T]) ToJSON() ([]byte, error) {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return json.Marshal(a.elements)
}

// FromJSON deserializes JSON data into the ArrayList
func (a *ArrayList[T]) FromJSON(data []byte) error {
    var elements []T
    if err := json.Unmarshal(data, &elements); err != nil {
        return err
    }
    a.mu.Lock()
    defer a.mu.Unlock()
    a.elements = elements
    return nil
}

// ToSlice returns a copy of the underlying slice.
func (a *ArrayList[T]) ToSlice() []T {
    a.mu.RLock()
    defer a.mu.RUnlock()

    // Create a copy of the slice to avoid exposing the original slice
    sliceCopy := make([]T, len(a.elements))
    copy(sliceCopy, a.elements)

    return sliceCopy
}

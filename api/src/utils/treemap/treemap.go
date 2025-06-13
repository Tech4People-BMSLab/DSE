package treemap

import (
	"encoding/json"
	"sync"
)

// Comparator is a function that compares two keys.
// It returns -1 if a < b, 0 if a == b, and 1 if a > b.
type Comparator[K any] func(a, b K) int

// TreeNode represents a node in the binary search tree.
type TreeNode[K any, V any] struct {
    Key   K             `json:"key"`
    Value V             `json:"value"`
    Left  *TreeNode[K, V] `json:"left,omitempty"`
    Right *TreeNode[K, V] `json:"right,omitempty"`
}

// Map is a generic ordered map implemented using a binary search tree.
type Map[K any, V any] struct {
    root       *TreeNode[K, V]
    comparator Comparator[K]
    mu         sync.RWMutex
}

// NewMap creates a new instance of Map with the given comparator.
func NewMap[K any, V any](comparator Comparator[K]) *Map[K, V] {
    return &Map[K, V]{comparator: comparator}
}

// Set adds a key-value pair to the map.
func (m *Map[K, V]) Set(key K, value V) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.root = m.set(m.root, key, value)
}

func (m *Map[K, V]) set(node *TreeNode[K, V], key K, value V) *TreeNode[K, V] {
    if node == nil {
        return &TreeNode[K, V]{Key: key, Value: value}
    }
    cmp := m.comparator(key, node.Key)
    if cmp < 0 {
        node.Left = m.set(node.Left, key, value)
    } else if cmp > 0 {
        node.Right = m.set(node.Right, key, value)
    } else {
        node.Value = value
    }
    return node
}

// Get retrieves the value associated with the key.
func (m *Map[K, V]) Get(key K) (V, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    node := m.get(m.root, key)
    if node != nil {
        return node.Value, true
    }
    var zeroValue V
    return zeroValue, false
}

func (m *Map[K, V]) get(node *TreeNode[K, V], key K) *TreeNode[K, V] {
    if node == nil {
        return nil
    }
    cmp := m.comparator(key, node.Key)
    if cmp < 0 {
        return m.get(node.Left, key)
    } else if cmp > 0 {
        return m.get(node.Right, key)
    } else {
        return node
    }
}

// Delete removes a key-value pair from the map.
func (m *Map[K, V]) Delete(key K) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.root = m.delete(m.root, key)
}

func (m *Map[K, V]) delete(node *TreeNode[K, V], key K) *TreeNode[K, V] {
    if node == nil {
        return nil
    }
    cmp := m.comparator(key, node.Key)
    if cmp < 0 {
        node.Left = m.delete(node.Left, key)
    } else if cmp > 0 {
        node.Right = m.delete(node.Right, key)
    } else {
        if node.Left == nil {
            return node.Right
        }
        if node.Right == nil {
            return node.Left
        }
        minNode := m.min(node.Right)
        node.Key = minNode.Key
        node.Value = minNode.Value
        node.Right = m.deleteMin(node.Right)
    }
    return node
}

func (m *Map[K, V]) min(node *TreeNode[K, V]) *TreeNode[K, V] {
    if node.Left == nil {
        return node
    }
    return m.min(node.Left)
}

func (m *Map[K, V]) deleteMin(node *TreeNode[K, V]) *TreeNode[K, V] {
    if node.Left == nil {
        return node.Right
    }
    node.Left = m.deleteMin(node.Left)
    return node
}

// Each applies a function to each key-value pair in the map in order.
func (m *Map[K, V]) Each(f func(key K, value V)) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    m.each(m.root, f)
}

func (m *Map[K, V]) each(node *TreeNode[K, V], f func(key K, value V)) {
    if node == nil {
        return
    }
    m.each(node.Left, f)
    f(node.Key, node.Value)
    m.each(node.Right, f)
}

// MapIterator is an iterator for Map.
type MapIterator[K any, V any] struct {
    stack []*TreeNode[K, V]
    mu    *sync.RWMutex
}

// Iterate returns a new iterator for the map.
func (m *Map[K, V]) Iterate() *MapIterator[K, V] {
    m.mu.RLock()
    iterator := &MapIterator[K, V]{mu: &m.mu}
    iterator.pushLeft(m.root)
    return iterator
}

func (it *MapIterator[K, V]) pushLeft(node *TreeNode[K, V]) {
    for node != nil {
        it.stack = append(it.stack, node)
        node = node.Left
    }
}

// Next moves the iterator to the next key-value pair and returns it.
func (it *MapIterator[K, V]) Next() (key K, value V, ok bool) {
    if len(it.stack) == 0 {
        it.mu.RUnlock() // Release lock when iteration is done
        return
    }
    node := it.stack[len(it.stack)-1]
    it.stack = it.stack[:len(it.stack)-1]
    it.pushLeft(node.Right)
    return node.Key, node.Value, true
}

// ToJSON serializes the Map to JSON
func (m *Map[K, V]) ToJSON() ([]byte, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return json.Marshal(m.root)
}

// FromJSON deserializes JSON data into the Map
func (m *Map[K, V]) FromJSON(data []byte) error {
    var root *TreeNode[K, V]
    if err := json.Unmarshal(data, &root); err != nil {
        return err
    }
    m.mu.Lock()
    defer m.mu.Unlock()
    m.root = root
    return nil
}

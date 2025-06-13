package hashmap

import (
	"encoding/json"
	"fmt"
	"sync"
)

// HashMap is a generic map structure
type HashMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

// NewHashMap creates a new instance of a generic hashmap
func NewHashMap[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{data: make(map[K]V)}
}

// Set adds a key-value pair to the map
func (h *HashMap[K, V]) Set(key K, value V) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data[key] = value
}

// Get retrieves the value associated with the key
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	value, ok := h.data[key]
	return value, ok
}

// MustGet retrieves the value associated with the key.
// It panics if the key does not exist.
func (h *HashMap[K, V]) MustGet(key K) V {
	h.mu.RLock()
	defer h.mu.RUnlock()
	value, ok := h.data[key]
	if !ok {
		panic("Key does not exist in HashMap")
	}
	return value
}

// Delete removes a key-value pair from the map
func (h *HashMap[K, V]) Delete(key K) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.data, key)
}

// MustDelete removes a key-value pair from the map.
// It panics if the key does not exist.
func (h *HashMap[K, V]) MustDelete(key K) {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.data[key]
	if !ok {
		panic("Key does not exist in HashMap")
	}
	delete(h.data, key)
}

// Has checks if the key exists in the map
func (h *HashMap[K, V]) Has(key K) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	_, ok := h.data[key]
	return ok
}

// Len returns the number of key-value pairs in the map
func (h *HashMap[K, V]) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.data)
}

// Each applies a function to each key-value pair in the map
func (h *HashMap[K, V]) Each(f func(key K, value V)) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for k, v := range h.data {
		f(k, v)
	}
}

// HashMapIterator is an iterator for HashMap
type HashMapIterator[K comparable, V any] struct {
	data []struct {
		key   K
		value V
	}
	index int
}

// Iterate returns a new iterator for the map
func (h *HashMap[K, V]) Iterate() *HashMapIterator[K, V] {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data := make([]struct {
		key   K
		value V
	}, 0, len(h.data))
	for k, v := range h.data {
		data = append(data, struct {
			key   K
			value V
		}{k, v})
	}
	return &HashMapIterator[K, V]{data: data, index: 0}
}

// Next moves the iterator to the next key-value pair and returns it
func (it *HashMapIterator[K, V]) Next() (key K, value V, ok bool) {
	if it.index < len(it.data) {
		kv := it.data[it.index]
		it.index++
		return kv.key, kv.value, true
	}
	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// ToValue returns a copy of the underlying map[K]V
func (h *HashMap[K, V]) ToMap() map[K]V {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a copy of the internal map to avoid exposing the original map directly
	copyMap := make(map[K]V, len(h.data))
	for key, value := range h.data {
		copyMap[key] = value
	}

	return copyMap
}

// ToAnyMap converts the HashMap[K, V] to a map[any]any
func (h *HashMap[K, V]) ToAnyMap() map[any]any {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a new map[any]any and copy all elements from the original map[K]V
	anyMap := make(map[any]any, len(h.data))
	for key, value := range h.data {
		anyMap[key] = value
	}

	return anyMap
}

// ToStringMap converts the HashMap[K, V] to a map[string]interface{}
func (h *HashMap[K, V]) ToStringMap() (map[string]interface{}, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a new map[string]interface{} and copy elements where key is string
	sm := make(map[string]interface{}, len(h.data))
	for key, value := range h.data {
		k, ok := any(key).(string) // Ensure the key is of type string
		if !ok {
			return nil, fmt.Errorf("key '%v' is not of type string", key)
		}
		sm[k] = value
	}

	return sm, nil
}


// ToJSON serializes the HashMap to JSON
func (h *HashMap[K, V]) ToJSON() ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return json.Marshal(h.data)
}

// FromJSON deserializes JSON data into the HashMap
func (h *HashMap[K, V]) FromJSON(data []byte) error {
	var tempMap map[K]V
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.data = tempMap
	return nil
}

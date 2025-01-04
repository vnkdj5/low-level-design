package utils

import "sync"

// Set is a thread-safe implementation of a generic Set.
type Set[T comparable] struct {
	mu    sync.RWMutex
	items map[T]struct{}
}

// NewSet creates and returns a new Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Add adds an item to the Set.
func (s *Set[T]) Add(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
}

// Remove removes an item from the Set.
func (s *Set[T]) Remove(item T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, item)
}

// Contains checks if an item exists in the Set.
func (s *Set[T]) Contains(item T) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.items[item]
	return exists
}

// Size returns the number of elements in the Set.
func (s *Set[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}

// Clear removes all items from the Set.
func (s *Set[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = make(map[T]struct{})
}

// Items returns all items in the Set as a slice.
func (s *Set[T]) Items() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]T, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}
	return keys
}

// Union returns a new Set that is the union of two Sets.
func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	for item := range s.items {
		result.Add(item)
	}
	s.mu.RUnlock()

	other.mu.RLock()
	for item := range other.items {
		result.Add(item)
	}
	other.mu.RUnlock()
	return result
}

// Intersection returns a new Set that is the intersection of two Sets.
func (s *Set[T]) Intersection(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	other.mu.RLock()
	for item := range s.items {
		if _, exists := other.items[item]; exists {
			result.Add(item)
		}
	}
	other.mu.RUnlock()
	s.mu.RUnlock()
	return result
}

// Difference returns a new Set that is the difference of two Sets (s - other).
func (s *Set[T]) Difference(other *Set[T]) *Set[T] {
	result := NewSet[T]()
	s.mu.RLock()
	other.mu.RLock()
	for item := range s.items {
		if _, exists := other.items[item]; !exists {
			result.Add(item)
		}
	}
	other.mu.RUnlock()
	s.mu.RUnlock()
	return result
}

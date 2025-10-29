// Package set provides a generic set implementation using Go's generics.
//
// Set is a memory-efficient set data structure that offers O(1) average time
// complexity for add, remove, and lookup operations. It uses Go's built-in map with
// empty structs as values to minimize memory usage while maintaining high performance.
//
// Thread Safety:
//
// Set is NOT thread-safe. Concurrent access must be synchronized externally
// using sync.Mutex or similar synchronization mechanisms. For example:
//
//	type ThreadSafeSet[T comparable] struct {
//	    mu sync.RWMutex
//	    table *set.Set[T]
//	}
//
// Memory Usage:
//
// The Set implementation uses empty structs (struct{}) as map values, which
// consume 0 bytes, making it memory-efficient for storing sets of data.
//
// Basic Example:
//
//	set := set.New[string]()
//	set.Add("apple")
//	set.Add("banana", "cherry")
//
//	if set.Contains("apple") {
//	    fmt.Println("Found apple!")
//	}
//
//	fruits := set.Items()
//
// Set Operations Example:
//
//	// Create two sets
//	fruits := set.New[string]()
//	fruits.Add("apple", "banana", "cherry")
//
//	tropical := set.New[string]()
//	tropical.Add("banana", "mango", "pineapple")
//
//	// Find intersection (common elements)
//	intersection := set.New[string]()
//	for _, item := range tropical.Items() {
//	    if fruits.Contains(item) {
//	        intersection.Add(item)
//	    }
//	}
//
//	// Find difference (in fruits but not in tropical)
//	difference := set.New[string]()
//	for _, item := range fruits.Items() {
//	    if !tropical.Contains(item) {
//	        difference.Add(item)
//	    }
//	}
package set

// Set is a generic set implementation that stores unique elements of type T.
// It uses a map with empty struct values to provide efficient memory usage.
// The type parameter T must satisfy the comparable constraint.
//
// Note: This implementation is not thread-safe. Concurrent access must be
// synchronized externally.
type Set[T comparable] struct {
	data map[T]struct{}
}

// New creates and returns a new empty Set for type T.
// The returned Set is ready to use with all methods.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

// Add inserts one or more items into the set.
// If any item already exists, the set remains unchanged for that item.
// This operation has O(n) time complexity where n is the number of items.
func (h *Set[T]) Add(items ...T) {
    for _, item := range items {
        h.data[item] = struct{}{}
    }
}

// Remove deletes one or more items from the set if they exist.
// If an item doesn't exist, this operation has no effect for that item.
// This operation has O(n) time complexity where n is the number of items.
func (h *Set[T]) Remove(items ...T) {
    for _, item := range items {
        delete(h.data, item)
    }
}

// Contains checks if an item exists in the set.
// Returns true if the item is present, false otherwise.
// This operation has O(1) average time complexity.
func (h *Set[T]) Contains(item T) bool {
    _, ok := h.data[item]
    return ok
}


// ContainsAny checks if any of the provided items exist in the set.
// Returns true if at least one item is found, false otherwise.
// This operation has O(n) time complexity where n is the number of items checked.
func (h *Set[T]) ContainsAny(items ...T) bool {
    for _, item := range items {
        if h.Contains(item) {
            return true
        }
    }
    return false
}

// ContainsAll checks if all of the provided items exist in the set.
// Returns true only if every item is found, false otherwise.
// This operation has O(n) time complexity where n is the number of items checked.
func (h *Set[T]) ContainsAll(items ...T) bool {
    for _, item := range items {
        if !h.Contains(item) {
            return false
        }
    }
    return true
}

// Size returns the number of items in the set.
// This operation has O(1) time complexity.
func (h *Set[T]) Size() int {
    return len(h.data)
}

// Clear removes all items from the set, resetting it to empty.
// This operation has O(1) time complexity but the garbage collector
// will need to free the memory of the previous map.
func (h *Set[T]) Clear() {
    h.data = make(map[T]struct{})
}

// Items returns all items in the set as a slice.
// The order of items in the returned slice is not guaranteed.
// This operation has O(n) time complexity where n is the number of items.
func (h *Set[T]) Items() []T {
    items := make([]T, 0, len(h.data))
    for item := range h.data {
        items = append(items, item)
    }
    return items
}

// Union returns a new Set containing all items from both this Set and another.
// This operation has O(n+m) time complexity where n and m are the sizes of the two sets.
func (h *Set[T]) Union(other *Set[T]) *Set[T] {
    result := New[T]()

    // Add all items from this set
    for item := range h.data {
        result.Add(item)
    }

    // Add all items from the other set
    for item := range other.data {
        result.Add(item)
    }

    return result
}

// Intersection returns a new Set containing only items that exist in both
// this Set and another.
// This operation has O(min(n,m)) time complexity where n and m are the sizes of the two sets.
func (h *Set[T]) Intersection(other *Set[T]) *Set[T] {
    result := New[T]()

    // Iterate over the smaller set for efficiency
    if h.Size() <= other.Size() {
        for item := range h.data {
            if other.Contains(item) {
                result.Add(item)
            }
        }
    } else {
        for item := range other.data {
            if h.Contains(item) {
                result.Add(item)
            }
        }
    }

    return result
}

// Difference returns a new Set containing items that exist in this Set
// but not in the other Set.
// This operation has O(n) time complexity where n is the size of this set.
func (h *Set[T]) Difference(other *Set[T]) *Set[T] {
    result := New[T]()

    for item := range h.data {
        if !other.Contains(item) {
            result.Add(item)
        }
    }

    return result
}

// IsSubset returns true if this Set is a subset of the other Set
// (i.e., all elements in this set are also in the other set).
// This operation has O(n) time complexity where n is the size of this set.
func (h *Set[T]) IsSubset(other *Set[T]) bool {
    // A set cannot be a subset of a smaller set
    if h.Size() > other.Size() {
        return false
    }

    for item := range h.data {
        if !other.Contains(item) {
            return false
        }
    }

    return true
}

// Equals returns true if this Set contains exactly the same elements as the other Set.
// This operation has O(n) time complexity where n is the size of this set.
func (h *Set[T]) Equals(other *Set[T]) bool {
    // If sizes differ, sets cannot be equal
    if h.Size() != other.Size() {
        return false
    }

    // Check if every element in this set is in the other set
    for item := range h.data {
        if !other.Contains(item) {
            return false
        }
    }

    return true
}

package collection

import (
	"sort"

	"golang.org/x/exp/constraints"
)

type Sortable[T interface{}] struct {
	items  []T
	sortFn func(a, b T) bool
}

// Len is the number of elements in the collection.
func (s *Sortable[T]) Len() int {
	return len(s.items)
}

// Less reports whether the element with index i
// must sort before the element with index j.
//
// If both Less(i, j) and Less(j, i) are false,
// then the elements at index i and j are considered equal.
// Sort may place equal elements in any order in the final result,
// while Stable preserves the original input order of equal elements.
//
// Less must describe a transitive ordering:
//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
//
// Note that floating-point comparison (the < operator on float32 or float64 values)
// is not a transitive ordering when not-a-number (NaN) values are involved.
// See Float64Slice.Less for a correct implementation for floating-point values.
func (s *Sortable[T]) Less(i int, j int) bool {
	return s.sortFn(s.items[i], s.items[j])
}

// Swap swaps the elements with indexes i and j.
func (s *Sortable[T]) Swap(i int, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// SortBy sorts the collection in-place using the given sort function.
func (s *Sortable[T]) SortBy(sortFn func(a, b T) bool) []T {
	s.sortFn = sortFn
	sort.Sort(s)
	return s.items
}

func Map[T interface{}, R interface{}](items []T, fn func(T) R) []R {
	result := make([]R, len(items))
	for _, item := range items {
		result = append(result, fn(item))
	}

	return result
}

func Filter[T interface{}](items []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}

	return result
}

func Find[T interface{}](items []T, fn func(T) bool) (T, bool) {
	var result T
	for _, item := range items {
		if fn(item) {
			return item, true
		}
	}

	return result, false
}

func Contains[T comparable](items []T, item T) bool {
	return IndexOf(items, item) != -1
}

func Reduce[T interface{}, R constraints.Ordered](items []T, fn func(R, T) R, first R) R {
	result := fn(first, items[0])
	for _, item := range items[1:] {
		result = fn(result, item)
	}

	return result
}

func IndexOf[T comparable](items []T, search T) int {
	for i, item := range items {
		if search == item {
			return i
		}
	}

	return -1
}

func LastIndexOf[T comparable](items []T, search T) int {
	for i := len(items) - 1; i >= 0; i-- {
		if search == items[i] {
			return i
		}
	}

	return -1
}

func Every[T interface{}](items []T, fn func(T) bool) bool {
	for _, item := range items {
		if !fn(item) {
			return false
		}
	}

	return true
}

func SortBy[T interface{}](items []T, fn func(a, b T) bool) []T {
	return NewSortable(items).SortBy(fn)
}

func GroupBy[T interface{}, K comparable](items []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range items {
		key := fn(item)
		result[key] = append(result[key], item)
	}

	return result
}

func NewSortable[T interface{}](items []T) *Sortable[T] {
	data := make([]T, len(items))
	copy(data, items)
	return &Sortable[T]{items: data}
}

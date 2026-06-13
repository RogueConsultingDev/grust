package it

import (
	"fmt"
	"reflect"
	"slices"
)

// Equal is an interface for types that can be compared for equality.
type Equal[T any] interface {
	Equal(other T) bool
}

// Unique filters out elements that have already been produced once during the iteration.
// For deduplication, the type T must be comparable or implement the Equal[T] interface.
// If the type T is comparable, a map is used to keep track of seen items.
// If the type T is not comparable, but implements the Equal[T] interface, the Equal method is used to compare against
// all previously seen items.
// Otherwise, an error is returned.
func (i *Iterator[T]) Unique() *Iterator[T] {
	var t T
	v := reflect.ValueOf(t)

	if v.Comparable() {
		return i.uniqueCmp()
	}

	_, ok := any(&t).(Equal[T])
	if ok {
		return i.uniqueEq()
	}

	it := func(yield func(T, error) bool) {
		yield(t, fmt.Errorf("can't use unique on non-comparable, non-Equal type: %T", t))
	}

	return &Iterator[T]{it}
}

func (i *Iterator[T]) uniqueCmp() *Iterator[T] {
	it := func(yield func(T, error) bool) {
		seen := make(map[any]struct{})

		for v, err := range i.it {
			if err != nil {
				var zero T
				yield(zero, err)

				return
			}

			_, ok := seen[v]
			if ok {
				continue
			}

			seen[v] = struct{}{}

			if !yield(v, err) {
				return
			}
		}
	}

	return &Iterator[T]{it}
}

func (i *Iterator[T]) uniqueEq() *Iterator[T] {
	it := func(yield func(T, error) bool) {
		var seen []T

		for v, err := range i.it {
			if err != nil {
				var zero T
				yield(zero, err)

				return
			}

			eq, _ := any(&v).(Equal[T])

			ok := slices.ContainsFunc(seen, eq.Equal)

			if ok {
				continue
			}

			seen = append(seen, v)

			if !yield(v, err) {
				return
			}
		}
	}

	return &Iterator[T]{it}
}

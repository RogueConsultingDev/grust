// Package it implements Rust-style Iterators
package it

import (
	"errors"
	"iter"
)

type Iterator[T any, U any] struct {
	it iter.Seq2[T, error]
}

type Tuple[T any, U any] struct {
	A T
	B U
}

// New creates an iterator from the given slice.
func New[T any](values []T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for i := range values {
				if !yield(values[i], nil) {
					return
				}
			}
		},
	}
}

// New2 creates an iterator from the given slice, with a different secondary type.
// This is useful for creating a mapping iterator.

// Deprecated: Un-necessary since generic methods.
func New2[T any, U any](values []T) *Iterator[T, U] {
	return &Iterator[T, U]{
		it: func(yield func(T, error) bool) {
			for i := range values {
				if !yield(values[i], nil) {
					return
				}
			}
		},
	}
}

// Reversed creates an iterator from the given slice in reverse order.
func Reversed[T any](values []T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			n := len(values) - 1
			for i := range values {
				if !yield(values[n-i], nil) {
					return
				}
			}
		},
	}
}

// From creates an iterator from the given iterator, with a different secondary type.
// This is useful for creating a mapping iterator.

// Deprecated: Un-necessary since generic methods.
func From[T any, U any](source *Iterator[T, any]) *Iterator[T, U] {
	return &Iterator[T, U]{
		it: func(yield func(T, error) bool) {
			for i, err := range source.it {
				if !yield(i, err) {
					return
				}
			}
		},
	}
}

// Repeat creates an iterator than endlessly repeats a unique element.
func Repeat[T any](v T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

// RepeatN creates an iterator than repeats a unique element a given number of times.
func RepeatN[T any](v T, n int) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for range n {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

// Incr creates an iterator that returns numbers, starting at 0 and endlessly incrementing by 1.
func Incr() *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

// IncrN creates an iterator that returns numbers, starting at 0 and endlessly incrementing by `n`.
func IncrN(n int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

// IncrFrom creates an iterator that returns numbers, starting at `start` and endlessly incrementing by 1.
func IncrFrom(start int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

// IncrNFrom creates an iterator that returns numbers, starting at `start` and endlessly incrementing by `n`.
func IncrNFrom(start int, n int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

// Range creates an iterator that returns numbers, starting at `start` (inclusive) and incrementing by 1 until `end`
// (exclusive).
func Range(start int, end int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; i < end; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

// Cycle creates an iterator that endlessly cycles through the given values.
func Cycle[T any](values []T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for {
				for i := range values {
					if !yield(values[i], nil) {
						return
					}
				}
			}
		},
	}
}

// Chain creates an iterator that yields all elements from the given slices.
func Chain[T any](slices ...[]T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for _, s := range slices {
				for i := range s {
					if !yield(s[i], nil) {
						return
					}
				}
			}
		},
	}
}

// Product creates an iterator that yields all possible pairs of elements from the given slices.
func Product[T any, U any](p []T, q []U) *Iterator[Tuple[T, U], any] {
	return &Iterator[Tuple[T, U], any]{
		it: func(yield func(Tuple[T, U], error) bool) {
			for i := range p {
				for j := range q {
					t := Tuple[T, U]{p[i], q[j]}
					if !yield(t, nil) {
						return
					}
				}
			}
		},
	}
}

// Zip creates an iterator that yields elements of both slices, one by one, until either slice is exhausted.
func Zip[T any, U any](a []T, b []U) *Iterator[Tuple[T, U], any] {
	return &Iterator[Tuple[T, U], any]{
		it: func(yield func(Tuple[T, U], error) bool) {
			for idx := range a {
				if idx >= len(b) {
					return
				}

				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
	}
}

// ZipEq creates an iterator that yields elements of both slices, one by one, as long as both slices are the same
// length.
func ZipEq[T any, U any](a []T, b []U) *Iterator[Tuple[T, U], any] {
	return &Iterator[Tuple[T, U], any]{
		it: func(yield func(Tuple[T, U], error) bool) {
			if len(a) != len(b) {
				var t Tuple[T, U]
				yield(t, errors.New("slices are not the same length"))

				return
			}

			for idx := range a {
				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
	}
}

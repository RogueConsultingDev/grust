package it

// Enumerator is a wrapper around the iterator's value and its index.
type Enumerator[T any] struct {
	Idx   int
	Value T
}

// NewEnumeratorFrom creates an iterator that wraps the given iterator's values with their index.
func NewEnumeratorFrom[T any](iter *Iterator[T]) *Iterator[Enumerator[T]] {
	it := func(yield func(Enumerator[T], error) bool) {
		idx := 0
		for v, err := range iter.it {
			enum := Enumerator[T]{idx, v}

			if !yield(enum, err) || err != nil {
				return
			}

			idx++
		}
	}

	return &Iterator[Enumerator[T]]{it}
}

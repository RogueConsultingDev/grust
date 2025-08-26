package it

func Filter[T any](iterator Iterator[T], f func(T) bool) Iterator[T] {
	inner := func(yield func(T, error) bool) {
		for v, err := range iterator.it {
			if f(v) && !yield(v, err) {
				return
			}
		}
	}

	return Iterator[T]{
		inner,
		nil,
	}
}

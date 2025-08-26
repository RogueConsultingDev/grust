package it

func Map[T any, U any](iterator Iterator[T], f func(T) (U, error)) Iterator[U] {
	inner := func(yield func(U, error) bool) {
		for v := range iterator.it {
			if !yield(f(v)) {
				return
			}
		}
	}

	return Iterator[U]{
		inner,
		nil,
	}
}

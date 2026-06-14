//go:build !go1.27

package it

// FilterMap applies a filtering and mapping function to each element and yields only the elements satisfy the filter.
// Unlike its Go 1.27 counterpart, the return type of the mapping function must be the same as the one of the iterator.
func (i *Iterator[T]) FilterMap(f func(T) (T, bool, error)) *Iterator[T] {
	inner := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if err != nil {
				var zero T
				yield(zero, err)

				return
			}

			u, ok, err := f(v)
			if ok || err != nil {
				if !yield(u, err) {
					return
				}
			}
		}
	}

	return &Iterator[T]{inner}
}

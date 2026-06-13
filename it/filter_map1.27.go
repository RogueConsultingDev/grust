//go:build go1.27

package it

// FilterMap applies a filtering and mapping function to each element and yields only the elements satisfy the filter.
// Unlike its Go <=1.26 counterpart, the return type of the mapping function can be different from the one of the
// iterator.
func (i *Iterator[T]) FilterMap[U any](f func(T) (U, bool, error)) *Iterator[U] {
	inner := func(yield func(U, error) bool) {
		for v, err := range i.it {
			if err != nil {
				var zero U
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

	return &Iterator[U]{inner}
}

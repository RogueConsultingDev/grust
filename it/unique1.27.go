//go:build go1.27

package it

// UniqueBy filters out elements that have already been produced once during the iteration. The deduplication is done
// by executing the provided function on all items and by using the result as a differentiator.
func (i *Iterator[T]) UniqueBy[U comparable](f func(T) U) *Iterator[T] {
	it := func(yield func(T, error) bool) {
		seen := make(map[U]struct{})

		for v, err := range i.it {
			if err != nil {
				var zero T
				yield(zero, err)

				return
			}

			k := f(v)
			_, ok := seen[k]
			if ok {
				continue
			}

			seen[k] = struct{}{}

			if !yield(v, err) {
				return
			}
		}
	}

	return &Iterator[T]{it}
}

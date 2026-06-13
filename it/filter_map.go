//go:build !go1.27

package it

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

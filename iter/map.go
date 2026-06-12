//go:build !go1.27

package it

func (i *Iterator[T]) Map(f func(T) (T, error)) *Iterator[T] {
	it := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if err != nil {
				var zero T
				yield(zero, err)

				return
			}

			u, err := f(v)
			if !yield(u, err) {
				return
			}
		}
	}

	return &Iterator[T]{it}
}

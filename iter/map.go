//go:build go1.27
package it

func (i *Iterator[T]) Map[U any](f func(T) (U, error)) *Iterator[U] {
	it := func(yield func(U, error) bool) {
		for v, err := range i.it {
			if err != nil {
				var zero U
				yield(zero, err)

				return
			}

			u, err := f(v)
			if !yield(u, err) {
				return
			}
		}
	}

	return &Iterator[U]{it}
}

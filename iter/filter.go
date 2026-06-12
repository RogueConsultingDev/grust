package it

func (i *Iterator[T]) Filter(predicate func(T) bool) *Iterator[T] {
	inner := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if err != nil {
				yield(v, err)

				return
			}

			if predicate(v) && !yield(v, err) {
				return
			}
		}
	}

	return &Iterator[T]{inner}
}

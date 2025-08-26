package it

func (i *Iterator[T, U]) Filter(predicate func(T) bool) *Iterator[T, U] {
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

	return &Iterator[T, U]{inner}
}

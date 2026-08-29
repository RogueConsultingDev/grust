package it

// Take creates an iterator that yields at most N elements.
func (i *Iterator[T]) Take(n int) *Iterator[T] {
	inner := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if err != nil {
				yield(v, err)

				return
			}

			if !yield(v, err) {
				return
			}

			n--

			if n == 0 {
				return
			}
		}
	}

	return &Iterator[T]{inner}
}

// Skip creates an iterator that skips up to N elements.
func (i *Iterator[T]) Skip(n int) *Iterator[T] {
	inner := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if n > 0 {
				n--

				continue
			}

			if err != nil {
				yield(v, err)

				return
			}

			if !yield(v, err) {
				return
			}
		}
	}

	return &Iterator[T]{inner}
}

package betteriter

// Take creates an iterator that yields at most N elements.
func (i *Iterator[T, U]) Take(n int) *Iterator[T, U] {
	inner := func(yield func(T, error) bool) {
		for v, err := range i.it {
			if err != nil {
				yield(v, err)
				return
			}

			if !yield(v, err) {
				return
			}

			n -= 1

			if n == 0 {
				return
			}
		}
	}

	return &Iterator[T, U]{inner}
}

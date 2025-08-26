package betteriter

func (i *Iterator[T, U]) FilterMap(f func(T) (U, bool, error)) *Iterator[*U, U] {
	inner := func(yield func(*U, error) bool) {
		for v, err := range i.it {
			if err != nil {
				yield(nil, err)
				return
			}

			u, ok, err := f(v)
			if ok || err != nil {
				if !yield(&u, err) {
					return
				}
			}
		}
	}

	return &Iterator[*U, U]{inner}
}

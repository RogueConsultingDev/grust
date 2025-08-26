package betteriter

type Enumerator[T any] struct {
	Idx   int
	Value T
}

func NewEnumeratorFrom[T any, U any](iter *Iterator[T, U]) *Iterator[Enumerator[T], U] {
	it := func(yield func(Enumerator[T], error) bool) {
		idx := 0
		for v, err := range iter.it {
			enum := Enumerator[T]{idx, v}

			if !yield(enum, err) || err != nil {
				return
			}

			idx += 1
		}
	}

	return &Iterator[Enumerator[T], U]{it}
}

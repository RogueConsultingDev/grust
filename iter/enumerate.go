package it

type Enumerator[T any] struct {
	Idx   int
	Value T
}

func NewEnumeratorFrom[T any](iter *Iterator[T]) *Iterator[Enumerator[T]] {
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

	return &Iterator[Enumerator[T]]{it}
}

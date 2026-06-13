//go:build go1.27

package it

// Fold applies a function against an accumulator and each element in the iterator, from left to right, to reduce it to
// a single value.
func (i *Iterator[T]) Fold[U any](init U, adder func(cur U, item T) U) (U, error) {
	current := init

	for v, err := range i.it {
		if err != nil {
			return current, err
		}

		current = adder(current, v)
	}

	return current, nil
}

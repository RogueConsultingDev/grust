package betteriter

import "errors"

func (i *Iterator[T, U]) Collect() ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}
		output = append(output, v)
	}

	return output, nil
}

func (i *Iterator[T, U]) Reversed() ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}
		output = append([]T{v}, output...)
	}

	return output, nil
}

func (i *Iterator[T, U]) Any(predicate func(T) bool) (bool, error) {
	for v, err := range i.it {
		if err != nil {
			return false, err
		}

		if predicate(v) {
			return true, nil
		}
	}

	return false, nil
}

func (i *Iterator[T, U]) All(predicate func(T) bool) (bool, error) {
	for v, err := range i.it {
		if err != nil {
			return false, err
		}

		if !predicate(v) {
			return false, nil
		}
	}

	return true, nil
}

func (i *Iterator[T, U]) First() (T, error) {
	for v, err := range i.it {
		if err != nil {
			var t T
			return t, err
		}

		return v, nil
	}

	var t T
	return t, errors.New("empty iterator")
}

func (i *Iterator[T, U]) Last() (T, error) {
	var t T
	found := false

	for v, err := range i.it {
		found = true
		if err != nil {
			return t, err
		}

		t = v

	}

	if !found {
		return t, errors.New("empty iterator")
	}

	return t, nil
}

func (i *Iterator[T, U]) Find(predicate func(T) bool) (T, bool, error) {
	for v, err := range i.it {
		if err != nil {
			var t T
			return t, false, err
		}

		if predicate(v) {
			return v, true, nil
		}

	}

	var t T

	return t, false, nil
}

func (i *Iterator[T, U]) Position(predicate func(T) bool) (int, bool, error) {
	idx := 0
	for v, err := range i.it {
		if err != nil {
			return 0, false, err
		}

		if predicate(v) {
			return idx, true, nil
		}

		idx += 1
	}

	return 0, false, nil
}

func (i *Iterator[T, U]) Fold(init U, adder func(cur U, item T) U) (U, error) {
	current := init

	for v, err := range i.it {
		if err != nil {
			return current, err
		}

		current = adder(current, v)
	}

	return current, nil
}

func Copied[T any, U any](i *Iterator[*T, U]) ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}
		output = append(output, *v)
	}

	return output, nil
}

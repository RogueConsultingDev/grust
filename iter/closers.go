package it

import (
	"errors"
	"iter"
)

// Iter returns the raw iterator
func (i *Iterator[T, U]) Iter() iter.Seq2[T, error] {
	return i.it
}

// Collect collects all elements from the iterator into a slice.
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

// Reversed collects all elements from the iterator into a slice in reverse order.
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

// Any checks if any element in the iterator satisfies the given predicate.
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

// All checks if all elements in the iterator satisfy the given predicate.
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

// First returns the first element of the iterator, or an error if the iterator is empty.
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

// Last returns the last element of the iterator, or an error if the iterator is empty.
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

// Find returns the first element of the iterator that satisfies the given predicate, if any.
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

// Position returns the index of the element of the iterator that satisfies the given predicate, if any.
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

// ForEach calls the given function for each element in the iterator.
func (i *Iterator[T, U]) ForEach(f func(T)) error {
	for v, err := range i.it {
		if err != nil {
			return err
		}

		f(v)
	}

	return nil
}

// Fold applies a function against an accumulator and each element in the iterator, from left to right, to reduce it to a single value.
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

// Copied dereferences all elements from the iterator into a slice.
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

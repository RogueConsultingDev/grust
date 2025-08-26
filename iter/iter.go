package it

import (
	"errors"
	"iter"
)

type Iterator[T any] struct {
	it  iter.Seq2[T, error]
	err error
}

type Tuple[T any, U any] struct {
	A T
	B U
}

// TODO: make this a functions?
func (i Iterator[T]) Collect() ([]T, error) {
	output := make([]T, 0)

	for v, err := range i.it {
		if err != nil {
			return nil, err
		}

		output = append(output, v)
	}

	return output, nil
}

func New[T any](values []T) Iterator[T] {
	return Iterator[T]{
		it: func(yield func(T, error) bool) {
			for _, v := range values {
				if !yield(v, nil) {
					return
				}
			}
		},
		err: nil,
	}
}

func NewRepeat[T any](v T) Iterator[T] {
	return Iterator[T]{
		it: func(yield func(T, error) bool) {
			for {
				if !yield(v, nil) {
					return
				}
			}
		},
		err: nil,
	}
}

func NewRepeatN[T any](v T, n int) Iterator[T] {
	return Iterator[T]{
		it: func(yield func(T, error) bool) {
			for range n {
				if !yield(v, nil) {
					return
				}
			}
		},
		err: nil,
	}
}

func Zip[T any, U any](a []T, b []U) Iterator[Tuple[T, U]] {
	return Iterator[Tuple[T, U]]{
		it: func(yield func(Tuple[T, U], error) bool) {
			for idx := range a {
				if idx >= len(b) {
					return
				}

				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
		err: nil,
	}
}

func ZipEq[T any, U any](a []T, b []U) Iterator[Tuple[T, U]] {
	return Iterator[Tuple[T, U]]{
		it: func(yield func(Tuple[T, U], error) bool) {
			if len(a) != len(b) {
				var t Tuple[T, U]
				yield(t, errors.New("slices are not the same length"))

				return
			}

			for idx := range a {
				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
		err: nil,
	}
}

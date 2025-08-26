package betteriter

import (
	"errors"
	"iter"
)

type Iterator[T any, U any] struct {
	it iter.Seq2[T, error]
}

type Tuple[T any, U any] struct {
	A T
	B U
}

func New[T any](values []T) *Iterator[*T, any] {
	return &Iterator[*T, any]{
		it: func(yield func(*T, error) bool) {
			for i := range values {
				if !yield(&values[i], nil) {
					return
				}
			}
		},
	}
}

func New2[T any, U any](values []T) *Iterator[*T, U] {
	return &Iterator[*T, U]{
		it: func(yield func(*T, error) bool) {
			for i := range values {
				if !yield(&values[i], nil) {
					return
				}
			}
		},
	}
}

func Reversed[T any](values []T) *Iterator[*T, any] {
	return &Iterator[*T, any]{
		it: func(yield func(*T, error) bool) {
			n := len(values) - 1
			for i := range values {
				if !yield(&values[n-i], nil) {
					return
				}
			}
		},
	}
}

func From[T any, U any](source *Iterator[T, any]) *Iterator[T, U] {
	return &Iterator[T, U]{
		it: func(yield func(T, error) bool) {
			for i, err := range source.it {
				if !yield(i, err) {
					return
				}
			}
		},
	}
}

func Repeat[T any](v T) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

func RepeatN[T any](v T, n int) *Iterator[T, any] {
	return &Iterator[T, any]{
		it: func(yield func(T, error) bool) {
			for i := 0; i < n; i++ {
				if !yield(v, nil) {
					return
				}
			}
		},
	}
}

func Incr() *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrN(n int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := 0; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrFrom(start int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func IncrNFrom(start int, n int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; ; i += n {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func Range(start int, end int) *Iterator[int, any] {
	return &Iterator[int, any]{
		it: func(yield func(int, error) bool) {
			for i := start; i < end; i++ {
				if !yield(i, nil) {
					return
				}
			}
		},
	}
}

func Cycle[T any](values []T) *Iterator[*T, any] {
	return &Iterator[*T, any]{
		it: func(yield func(*T, error) bool) {
			for {
				for i := range values {
					if !yield(&values[i], nil) {
						return
					}
				}
			}
		},
	}
}

func Chain[T any](slices ...[]T) *Iterator[*T, any] {
	return &Iterator[*T, any]{
		it: func(yield func(*T, error) bool) {
			for _, s := range slices {
				for i := range s {
					if !yield(&s[i], nil) {
						return
					}
				}
			}
		},
	}
}

func Product[T any, U any](p []T, q []U) *Iterator[*Tuple[*T, *U], any] {
	return &Iterator[*Tuple[*T, *U], any]{
		it: func(yield func(*Tuple[*T, *U], error) bool) {
			for i := range p {
				for j := range q {
					t := Tuple[*T, *U]{&p[i], &q[j]}
					if !yield(&t, nil) {
						return
					}
				}
			}
		},
	}
}

func Zip[T any, U any](a []T, b []U) *Iterator[Tuple[T, U], any] {
	return &Iterator[Tuple[T, U], any]{
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
	}
}

func ZipEq[T any, U any](a []T, b []U) *Iterator[Tuple[T, U], any] {
	return &Iterator[Tuple[T, U], any]{
		it: func(yield func(Tuple[T, U], error) bool) {
			if len(a) != len(b) {
				yield(Tuple[T, U]{}, errors.New("slices are not the same length"))
				return
			}

			for idx := range a {
				t := Tuple[T, U]{a[idx], b[idx]}

				if !yield(t, nil) {
					return
				}
			}
		},
	}
}

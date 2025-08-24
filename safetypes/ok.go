package st

import "fmt"

// Ok creates an Ok variant of Result from the value.
func Ok[T any](val T) Result[T] {
	return &ok[T]{
		val: val,
	}
}

type ok[T any] struct {
	val T
}

func (o *ok[T]) IsOk() bool {
	return true
}

func (o *ok[T]) IsOkAnd(f func(T) bool) bool {
	return f(o.val)
}

func (o *ok[T]) IsErr() bool {
	return false
}

func (o *ok[T]) IsErrAnd(_ func(error) bool) bool {
	return false
}

func (o *ok[T]) Expect(_ string) T {
	return o.val
}

func (o *ok[T]) ExpectErr(msg string) error {
	panic(fmt.Errorf("%s: %v", msg, o.val))
}

func (o *ok[T]) Inspect(f func(*T)) Result[T] {
	f(&o.val)

	return o
}

func (o *ok[T]) InspectErr(_ func(error)) Result[T] {
	return o
}

func (o *ok[T]) Unwrap() T {
	return o.val
}

func (o *ok[T]) UnwrapOr(_ T) T {
	return o.val
}

func (o *ok[T]) UnwrapOrElse(_ func() T) T {
	return o.val
}

func (o *ok[T]) UnwrapOrDefault() T {
	return o.val
}

func (o *ok[T]) UnwrapErr() error {
	panic(fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", o.val))
}

func (o *ok[T]) WrapErr(_ string) Result[T] {
	return o
}

func (o *ok[T]) Expand() (T, error) {
	return o.val, nil
}

func (o *ok[T]) String() string {
	return fmt.Sprintf("Ok(%v)", o.val)
}

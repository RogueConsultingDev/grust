package st

import "fmt"

// Err creates an Err variant of Result from the error.
func Err[T any](err error) Result[T] {
	return &errT[T]{
		err: err,
	}
}

type errT[T any] struct {
	err error
}

func (e *errT[T]) IsOk() bool {
	return false
}

func (e *errT[T]) IsOkAnd(_ func(T) bool) bool {
	return false
}

func (e *errT[T]) IsErr() bool {
	return true
}

func (e *errT[T]) IsErrAnd(f func(error) bool) bool {
	return f(e.err)
}

func (e *errT[T]) Expect(msg string) T {
	panic(fmt.Errorf("%s: %w", msg, e.err))
}

func (e *errT[T]) ExpectErr(_ string) error {
	return e.err
}

func (e *errT[T]) Inspect(_ func(*T)) Result[T] {
	return e
}

func (e *errT[T]) InspectErr(f func(error)) Result[T] {
	f(e.err)

	return e
}

func (e *errT[T]) Unwrap() T {
	panic(fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", e.err))
}

func (e *errT[T]) UnwrapOr(def T) T {
	return def
}

func (e *errT[T]) UnwrapOrElse(f func() T) T {
	return f()
}

func (e *errT[T]) UnwrapOrDefault() T {
	var v T

	return v
}

func (e *errT[T]) UnwrapErr() error {
	return e.err
}

// AsOptionValue converts a Result to a Some when res is result.Ok or None when res is result.Err.
func (e *errT[T]) AsOptionValue() Option[T] {
	return None[T]()
}

// AsOptionErr converts a Result to a Some when res is result.Err or None when res is result.Ok.
func (e *errT[T]) AsOptionErr() Option[error] {
	return Some[error](e.err)
}

func (e *errT[T]) WrapErr(msg string) Result[T] {
	return Err[T](fmt.Errorf("%s: %w", msg, e.err))
}

func (e *errT[T]) Expand() (T, error) {
	var zero T

	return zero, e.err
}

func (e *errT[T]) String() string {
	return fmt.Sprintf("Err(%v)", e.err)
}

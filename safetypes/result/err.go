// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

import "fmt"

// Err creates an Err variant of Result from the error.
func Err[T any, E error](err E) Result[T, E] {
	return errT[T, E]{
		err: err,
	}
}

type errT[T any, E error] struct {
	err E
}

func (e errT[T, E]) IsOk() bool {
	return false
}

func (e errT[T, E]) IsOkAnd(_ func(T) bool) bool {
	return false
}

func (e errT[T, E]) IsErr() bool {
	return true
}

func (e errT[T, E]) IsErrAnd(f func(error) bool) bool {
	return f(e.err)
}

func (e errT[T, E]) Expect(msg string) T {
	panic(fmt.Errorf("%s: %w", msg, e.err))
}

func (e errT[T, E]) ExpectErr(_ string) E {
	return e.err
}

func (e errT[T, E]) Inspect(_ func(*T)) Result[T, E] {
	return e
}

func (e errT[T, E]) InspectErr(f func(*E)) Result[T, E] {
	f(&e.err)

	return e
}

func (e errT[T, E]) Unwrap() T {
	panic(fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", e.err))
}

func (e errT[T, E]) UnwrapOr(def T) T {
	return def
}

func (e errT[T, E]) UnwrapOrElse(f func() T) T {
	return f()
}

func (e errT[T, E]) UnwrapOrDefault() T {
	var v T

	return v
}

func (e errT[T, E]) UnwrapErr() E {
	return e.err
}

func (e errT[T, E]) String() string {
	return fmt.Sprintf("Err(%v)", e.err)
}

package st

import (
	"fmt"
	"reflect"
)

// Ok creates an Ok variant of Result from the value.
func Ok[T any](val T) *Result[T] {
	return &Result[T]{
		ok:  true,
		val: val,
		err: nil,
	}
}

// Err creates an Err variant of Result from the error.
func Err[T any](err error) *Result[T] {
	var v T

	return &Result[T]{
		ok:  false,
		val: v,
		err: err,
	}
}

// ResultOf creates a Result from the given value and error.
func ResultOf[T any](val T, err error) *Result[T] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return Err[T](err)
	}

	return Ok(val)
}

// MapResult maps a Result<T> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func MapResult[T any, U any](res *Result[T], f func(T) U) *Result[U] {
	if res.IsErr() {
		return Err[U](res.UnwrapErr())
	}

	return Ok(f(res.Unwrap()))
}

// MapResultOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapResultOr[T any, U any](res *Result[T], def U, f func(T) U) U {
	if res.IsErr() {
		return def
	}

	return f(res.Unwrap())
}

// MapResultOrElse maps a Result<T> to U by applying fallback function default to a contained Err value, or function
// f to a contained Ok value.
func MapResultOrElse[T any, U any](
	res *Result[T],
	factory func() U,
	mapper func(T) U,
) U {
	if res.IsErr() {
		return factory()
	}

	return mapper(res.Unwrap())
}

// MapResultErr maps a Result<T> to Result<T, F> by applying a function to a contained Err value, leaving an Ok value
// untouched.
func MapResultErr[T any](res *Result[T], f func(error) error) *Result[T] {
	if res.IsOk() {
		return res
	}

	return Err[T](f(res.UnwrapErr()))
}

// Result is a type that represents either success (Ok) or failure (Err).
type Result[T any] struct {
	ok  bool
	val T
	err error
}

// IsOk returns `true` if the result is Ok.
func (r *Result[T]) IsOk() bool {
	return r.ok
}

// IsOkAnd returns `true` if the result is Ok and the value inside of it matches a predicate.
func (r *Result[T]) IsOkAnd(f func(T) bool) bool {
	return r.ok && f(r.val)
}

// IsErr returns `true` if the result is Err.
func (r *Result[T]) IsErr() bool {
	return !r.ok
}

// IsErrAnd returns `true` if the result is Err and the value inside of it matches a predicate.
func (r *Result[T]) IsErrAnd(f func(error) bool) bool {
	return !r.ok && f(r.err)
}

// Expect returns the contained Ok value, consuming the self value. Panics if the value is an Err, with a panic
// message including the passed message, and the content of the Err.
func (r *Result[T]) Expect(msg string) T {
	if r.ok {
		return r.val
	}

	panic(fmt.Errorf("%s: %w", msg, r.err))
}

// ExpectErr returns the contained Err value, consuming the self value. Panics if the value is an Ok, with a panic
// message including the passed message, and the content of the Ok.
func (r *Result[T]) ExpectErr(msg string) error {
	if !r.ok {
		return r.err
	}

	panic(fmt.Errorf("%s: %v", msg, r.val))
}

// Unwrap returns the contained Ok value, consuming the self value. Panics if the value is an Err, with a panic
// message provided by the Err's value.
func (r *Result[T]) Unwrap() T {
	if r.ok {
		return r.val
	}

	panic(fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", r.err))
}

// UnwrapOr returns the contained Ok value or a provided default.
func (r *Result[T]) UnwrapOr(def T) T {
	if r.ok {
		return r.val
	}

	return def
}

// UnwrapOrElse returns the contained Ok value or computes it from a closure.
func (r *Result[T]) UnwrapOrElse(f func() T) T {
	if r.ok {
		return r.val
	}

	return f()
}

// UnwrapOrDefault returns the contained Ok value or a default.
func (r *Result[T]) UnwrapOrDefault() T {
	if r.ok {
		return r.val
	}

	var def T

	return def
}

// UnwrapErr returns the contained Err value, consuming the self value. Panics if the value is an Ok, with a custom
// panic message provided by the Ok's value.
func (r *Result[T]) UnwrapErr() error {
	if r.ok {
		panic(fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", r.val))
	}

	return r.err
}

// Inspect calls a function with a reference to the contained value if Ok. Returns the original result.
func (r *Result[T]) Inspect(f func(*T)) *Result[T] {
	if r.ok {
		f(&r.val)
	}

	return r
}

// InspectErr calls a function with a reference to the contained value if Err. Returns the original result.
func (r *Result[T]) InspectErr(f func(error)) *Result[T] {
	if !r.ok {
		f(r.err)
	}

	return r
}

// AsOptionValue converts a Result to a Some when res is result.Ok or None when res is result.Err.
func (r *Result[T]) AsOptionValue() *Option[T] {
	if r.ok {
		return Some(r.val)
	}

	return None[T]()
}

// AsOptionErr converts a Result to a Some when res is result.Err or None when res is result.Ok.
func (r *Result[T]) AsOptionErr() *Option[error] {
	if !r.ok {
		return Some(r.err)
	}

	return None[error]()
}

// WrapErr wraps the error of an Err, leaving Ok untouched.
func (r *Result[T]) WrapErr(msg string) *Result[T] {
	if !r.ok {
		return Err[T](fmt.Errorf("%s: %w", msg, r.err))
	}

	return r
}

// Expand returns the Result as a standard Go (T, error).
func (r *Result[T]) Expand() (T, error) {
	return r.val, r.err
}

func (r *Result[T]) String() string {
	if r.ok {
		return fmt.Sprintf("Ok(%v)", r.val)
	}

	return fmt.Sprintf("Err(%v)", r.err)
}

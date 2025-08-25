package st

import (
	"fmt"
	"reflect"
)

// Result is a type that represents either success (Ok) or failure (Err).
type Result[T any] interface {
	// IsOk returns `true` if the result is Ok.
	IsOk() bool
	// IsOkAnd returns `true` if the result is Ok and the value inside of it matches a predicate.
	IsOkAnd(f func(T) bool) bool
	// IsErr returns `true` if the result is Err.
	IsErr() bool
	// IsErrAnd returns `true` if the result is Err and the value inside of it matches a predicate.
	IsErrAnd(f func(error) bool) bool
	// Expect returns the contained Ok value, consuming the self value. Panics if the value is an Err, with a panic
	// message including the passed message, and the content of the Err.
	Expect(msg string) T
	// ExpectErr returns the contained Err value, consuming the self value. Panics if the value is an Ok, with a panic
	// message including the passed message, and the content of the Ok.
	ExpectErr(msg string) error
	// Unwrap returns the contained Ok value, consuming the self value. Panics if the value is an Err, with a panic
	// message provided by the Err's value.
	Unwrap() T
	// UnwrapOr returns the contained Ok value or a provided default.
	UnwrapOr(def T) T
	// UnwrapOrElse returns the contained Ok value or computes it from a closure.
	UnwrapOrElse(f func() T) T
	// UnwrapOrDefault returns the contained Ok value or a default.
	UnwrapOrDefault() T
	// UnwrapErr returns the contained Err value, consuming the self value. Panics if the value is an Ok, with a custom
	// panic message provided by the Ok's value.
	UnwrapErr() error
	// Inspect calls a function with a reference to the contained value if Ok. Returns the original result.
	Inspect(f func(*T)) Result[T]
	// InspectErr calls a function with a reference to the contained value if Err. Returns the original result.
	InspectErr(f func(error)) Result[T]
	// AsOptionValue converts a Result to a Some when res is result.Ok or None when res is result.Err.
	AsOptionValue() Option[T]
	// AsOptionErr converts a Result to a Some when res is result.Err or None when res is result.Ok.
	AsOptionErr() Option[error]

	// WrapErr wraps the error of an Err, leaving Ok untouched
	WrapErr(msg string) Result[T]

	// Expand returns the Result as a standard Go (T, error)
	Expand() (T, error)

	fmt.Stringer
}

// Ok creates an Ok variant of Result from the value.
func Ok[T any](val T) Result[T] {
	return &result[T]{
		ok:  true,
		val: val,
		err: nil,
	}
}

// Err creates an Err variant of Result from the error.
func Err[T any](err error) Result[T] {
	var v T

	return &result[T]{
		ok:  false,
		val: v,
		err: err,
	}
}

// ResultOf creates a Result from the given value and error.
func ResultOf[T any](val T, err error) Result[T] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return Err[T](err)
	}

	return Ok(val)
}

// MapResult maps a Result<T> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func MapResult[T any, U any](res Result[T], f func(T) U) Result[U] {
	if res.IsErr() {
		return Err[U](res.UnwrapErr())
	}

	return Ok(f(res.Unwrap()))
}

// MapResultOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapResultOr[T any, U any](res Result[T], def U, f func(T) U) U {
	if res.IsErr() {
		return def
	}

	return f(res.Unwrap())
}

// MapResultOrElse maps a Result<T> to U by applying fallback function default to a contained Err value, or function
// f to a contained Ok value.
func MapResultOrElse[T any, U any](
	res Result[T],
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
func MapResultErr[T any](res Result[T], f func(error) error) Result[T] {
	if res.IsOk() {
		return res
	}

	return Err[T](f(res.UnwrapErr()))
}

type result[T any] struct {
	ok  bool
	val T
	err error
}

func (r *result[T]) IsOk() bool {
	return r.ok
}

func (r *result[T]) IsOkAnd(f func(T) bool) bool {
	return r.ok && f(r.val)
}

func (r *result[T]) IsErr() bool {
	return !r.ok
}

func (r *result[T]) IsErrAnd(f func(error) bool) bool {
	return !r.ok && f(r.err)
}

func (r *result[T]) Inspect(f func(*T)) Result[T] {
	if r.ok {
		f(&r.val)
	}

	return r
}

func (r *result[T]) InspectErr(f func(error)) Result[T] {
	if !r.ok {
		f(r.err)
	}

	return r
}

func (r *result[T]) Expect(msg string) T {
	if r.ok {
		return r.val
	}

	panic(fmt.Errorf("%s: %w", msg, r.err))
}

func (r *result[T]) ExpectErr(msg string) error {
	if !r.ok {
		return r.err
	}

	panic(fmt.Errorf("%s: %v", msg, r.val))
}

func (r *result[T]) Unwrap() T {
	if r.ok {
		return r.val
	}

	panic(fmt.Errorf("called `Result.Unwrap()` on an `Err` value: %w", r.err))
}

func (r *result[T]) UnwrapOr(def T) T {
	if r.ok {
		return r.val
	}

	return def
}

func (r *result[T]) UnwrapOrElse(f func() T) T {
	if r.ok {
		return r.val
	}

	return f()
}

func (r *result[T]) UnwrapOrDefault() T {
	if r.ok {
		return r.val
	}

	var def T

	return def
}

func (r *result[T]) UnwrapErr() error {
	if r.ok {
		panic(fmt.Errorf("called `Result.UnwrapErr()` on an `Ok` value: %v", r.val))
	}

	return r.err
}

func (r *result[T]) AsOptionValue() Option[T] {
	if r.ok {
		return Some(r.val)
	}

	return None[T]()
}

func (r *result[T]) AsOptionErr() Option[error] {
	if !r.ok {
		return Some(r.err)
	}

	return None[error]()
}

func (r *result[T]) WrapErr(msg string) Result[T] {
	if !r.ok {
		return Err[T](fmt.Errorf("%s: %w", msg, r.err))
	}

	return r
}

func (r *result[T]) Expand() (T, error) {
	return r.val, r.err
}

func (r *result[T]) String() string {
	if r.ok {
		return fmt.Sprintf("Ok(%v)", r.val)
	}

	return fmt.Sprintf("Err(%v)", r.err)
}

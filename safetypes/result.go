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
	// Inspect calls a function with a reference to the contained value if Ok. Returns the original result.
	Inspect(f func(*T)) Result[T]
	// InspectErr calls a function with a reference to the contained value if Err. Returns the original result.
	InspectErr(f func(error)) Result[T]
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

	// WrapErr wraps the error of an Err, leaving Ok untouched
	WrapErr(msg string) Result[T]

	// Expand returns the Result as a standard Go (T, error)
	Expand() (T, error)

	fmt.Stringer
}

// ResultOf creates a Result from the given value and error.
func ResultOf[T any](val T, err error) Result[T] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return &errT[T]{
			err: err,
		}
	}

	return &ok[T]{
		val: val,
	}
}

// MapResult maps a Result<T> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func MapResult[T any, U any](res Result[T], f func(T) U) Result[U] {
	s, isOk := res.(*ok[T])
	if !isOk {
		return &errT[U]{res.UnwrapErr()}
	}

	val := f(s.val)

	return &ok[U]{val}
}

// MapResultOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapResultOr[T any, U any](res Result[T], def U, f func(T) U) U {
	s, isOk := res.(*ok[T])
	if !isOk {
		return def
	}

	return f(s.val)
}

// MapResultOrElse maps a Result<T> to U by applying fallback function default to a contained Err value, or function
// f to a contained Ok value.
func MapResultOrElse[T any, U any](
	res Result[T],
	factory func() U,
	mapper func(T) U,
) U {
	s, isOk := res.(*ok[T])
	if !isOk {
		return factory()
	}

	return mapper(s.val)
}

// MapResultErr maps a Result<T> to Result<T, F> by applying a function to a contained Err value, leaving an Ok value
// untouched.
func MapResultErr[T any](res Result[T], f func(error) error) Result[T] {
	s, isOk := res.(*ok[T])
	if isOk {
		return s
	}

	return &errT[T]{f(res.UnwrapErr())}
}

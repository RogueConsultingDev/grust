package st

import (
	"fmt"
	"reflect"
)

// Result is a type that represents either success (Ok) or failure (Err).
type Result[T any, E error] interface {
	// IsOk returns `true` if the result is Ok.
	IsOk() bool
	// IsOkAnd returns `true` if the result is Ok and the value inside of it matches a predicate.
	IsOkAnd(f func(T) bool) bool
	// IsErr returns `true` if the result is Err.
	IsErr() bool
	// IsErrAnd returns `true` if the result is Err and the value inside of it matches a predicate.
	IsErrAnd(f func(error) bool) bool
	// Inspect calls a function with a reference to the contained value if Ok. Returns the original result.
	Inspect(f func(*T)) Result[T, E]
	// InspectErr calls a function with a reference to the contained value if Err. Returns the original result.
	InspectErr(f func(*E)) Result[T, E]
	// Expect returns the contained Ok value, consuming the self value. Panics if the value is an Err, with a panic
	// message including the passed message, and the content of the Err.
	Expect(msg string) T
	// ExpectErr returns the contained Err value, consuming the self value. Panics if the value is an Ok, with a panic
	// message including the passed message, and the content of the Ok.
	ExpectErr(msg string) E
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
	UnwrapErr() E

	fmt.Stringer
}

// ResultOf creates a Result from the given value and error.
func ResultOf[T any, E error](val T, err E) Result[T, E] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return errT[T, E]{
			err: err,
		}
	}

	return ok[T, E]{
		val: val,
	}
}

// MapResult maps a Result<T, E> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func MapResult[T any, U any, E error](res Result[T, E], f func(T) U) Result[U, E] {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return errT[U, E]{res.UnwrapErr()}
	}

	val := f(s.val)

	return ok[U, E]{val}
}

// MapResultOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapResultOr[T any, U any, E error](res Result[T, E], def U, f func(T) U) U {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return def
	}

	return f(s.val)
}

// MapResultOrElse maps a Result<T, E> to U by applying fallback function default to a contained Err value, or function
// f to a contained Ok value.
func MapResultOrElse[T any, U any, E error](
	res Result[T, E],
	factory func() U,
	mapper func(T) U,
) U {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return factory()
	}

	return mapper(s.val)
}

// MapResultErr maps a Result<T, E> to Result<T, F> by applying a function to a contained Err value, leaving an Ok value
// untouched.
func MapResultErr[T any, E error, F error](res Result[T, E], f func(E) F) Result[T, F] {
	s, isOk := res.(ok[T, E])
	if isOk {
		return (ok[T, F])(s)
	}

	return errT[T, F]{f(res.UnwrapErr())}
}

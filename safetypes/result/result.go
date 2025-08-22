// Package result implements https://doc.rust-lang.org/std/result/enum.Result.html
package result

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
	Inspect(f func(*T)) Result[T, E]
	InspectErr(f func(*E)) Result[T, E]
	Expect(msg string) T
	ExpectErr(msg string) E
	Unwrap() T
	UnwrapOr(def T) T
	UnwrapOrElse(f func() T) T
	UnwrapOrDefault() T
	UnwrapErr() E

	fmt.Stringer
}

// Of creates a Result from the given value and error.
func Of[T any, E error](val T, err E) Result[T, E] {
	if !reflect.ValueOf(&err).Elem().IsNil() {
		return errT[T, E]{
			err: err,
		}
	}

	return ok[T, E]{
		val: val,
	}
}

// Map maps a Result<T, E> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func Map[T any, U any, E error](res Result[T, E], f func(T) U) Result[U, E] {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return errT[U, E]{res.UnwrapErr()}
	}

	val := f(s.val)

	return ok[U, E]{val}
}

// MapOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func MapOr[T any, U any, E error](res Result[T, E], def U, f func(T) U) U {
	s, isOk := res.(ok[T, E])
	if !isOk {
		return def
	}

	return f(s.val)
}

// MapOrElse maps a Result<T, E> to U by applying fallback function default to a contained Err value, or function f to a
// contained Ok value.
func MapOrElse[T any, U any, E error](
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

// MapErr maps a Result<T, E> to Result<T, F> by applying a function to a contained Err value, leaving an Ok value
// untouched.
func MapErr[T any, E error, F error](res Result[T, E], f func(E) F) Result[T, F] {
	s, isOk := res.(ok[T, E])
	if isOk {
		return (ok[T, F])(s)
	}

	return errT[T, F]{f(res.UnwrapErr())}
}

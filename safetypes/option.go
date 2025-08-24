package st

import (
	"fmt"
	"reflect"
)

// Option is a type that represents either a value (Some) or not (None).
type Option[T any] interface {
	IsNone() bool
	IsNoneOr(f func(T) bool) bool
	IsSome() bool
	IsSomeAnd(f func(T) bool) bool
	Expect(msg string) T
	Unwrap() T
	UnwrapOr(def T) T
	UnwrapOrElse(f func() T) T
	UnwrapOrDefault() T
	Inspect(f func(T)) Option[T]
	Filter(f func(T) bool) Option[T]
	Or(other Option[T]) Option[T]
	OrElse(f func() Option[T]) Option[T]
	Xor(other Option[T]) Option[T]

	fmt.Stringer
}

// OptionOf creates an Option from the given value.
func OptionOf[T any](val T) Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return none[T]{}
	}

	return Some(val)
}

// MapOption maps an Option<T> to Option<U> by applying a function to a contained value (if Some) or returns None
// (if None).
func MapOption[T any, U any](opt Option[T], f func(T) U) Option[U] {
	s, ok := opt.(some[T])
	if !ok {
		return none[U]{}
	}

	return Some(f(s.val))
}

// MapOptionOr returns the provided default result (if None), or applies a function to the contained value (if Some).
func MapOptionOr[T any, U any](opt Option[T], def U, f func(T) U) U {
	s, ok := opt.(some[T])
	if !ok {
		return def
	}

	return f(s.val)
}

// MapOptionOrElse computes a default function result (if None), or applies a different function to the contained value
// (if Some).
func MapOptionOrElse[T any, U any](opt Option[T], factory func() U, f func(T) U) U {
	s, ok := opt.(some[T])
	if !ok {
		return factory()
	}

	return f(s.val)
}

// And returns None if the option is None, otherwise returns `optb`.
func And[T any, U any](opt Option[T], other Option[U]) Option[U] {
	_, ok := opt.(some[T])
	if !ok {
		return none[U]{}
	}

	return other
}

// AndThen returns None if the option is None, otherwise calls `f` with the wrapped value and returns the result.
func AndThen[T any, U any](opt Option[T], f func(T) Option[U]) Option[U] {
	s, ok := opt.(some[T])
	if !ok {
		return none[U]{}
	}

	return f(s.val)
}

package st

import (
	"fmt"
	"reflect"
)

// Option is a type that represents either a value (Some) or not (None).
type Option[T any] interface {
	// IsNone returns true if the option is a None value.
	IsNone() bool
	// IsNoneOr returns true if the option is a None or the value inside of it matches a predicate.
	IsNoneOr(f func(T) bool) bool
	// IsSome returns true if the option is a Some value.
	IsSome() bool
	// IsSomeAnd returns true if the option is a Some and the value inside of it matches a predicate.
	IsSomeAnd(f func(T) bool) bool
	// Expect returns the contained Some value, consuming the self value. Panics if the value is a None with a custom
	// panic message provided by msg.
	Expect(msg string) T
	// Unwrap returns the contained Some value, consuming the self value. Panics if the self value equals None.
	Unwrap() T
	// UnwrapOr returns the contained Some value or a provided default.
	UnwrapOr(def T) T
	// UnwrapOrElse returns the contained Some value or computes it from a closure.
	UnwrapOrElse(f func() T) T
	// UnwrapOrDefault returns the contained Some value or a default.
	UnwrapOrDefault() T
	// Inspect calls a function with a reference to the contained value if Some. Returns the original option.
	Inspect(f func(T)) Option[T]
	// Filter returns None if the option is None, otherwise calls predicate with the wrapped value and returns:
	//  * Some(t) if predicate returns true (where t is the wrapped value), and
	//  * None if predicate returns false.
	Filter(f func(T) bool) Option[T]
	// Or returns the option if it contains a value, otherwise returns optb.
	Or(other Option[T]) Option[T]
	// OrElse returns the option if it contains a value, otherwise calls f and returns the result.
	OrElse(f func() Option[T]) Option[T]
	// Xor returns Some if exactly one of self, optb is Some, otherwise returns None.
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

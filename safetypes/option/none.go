package option

import (
	"errors"
)

type none[T any] struct{}

// None creates a None variant of Option.
func None() Option[any] {
	return none[any]{}
}

func (n none[T]) IsNone() bool {
	return true
}

func (n none[T]) IsNoneOr(func(T) bool) bool {
	return true
}

func (n none[T]) IsSome() bool {
	return false
}

func (n none[T]) IsSomeAnd(func(T) bool) bool {
	return false
}

func (n none[T]) Expect(msg string) T {
	panic(errors.New(msg))
}

func (n none[T]) Unwrap() T {
	panic(errors.New("called `Option.Unwrap()` on a `None` value"))
}

func (n none[T]) UnwrapOr(def T) T {
	return def
}

func (n none[T]) UnwrapOrElse(f func() T) T {
	return f()
}

func (n none[T]) UnwrapOrDefault() T {
	var v T

	return v
}

func (n none[T]) Inspect(_ func(T)) Option[T] {
	return n
}

func (n none[T]) Filter(_ func(T) bool) Option[T] {
	return n
}

func (n none[T]) Or(other Option[T]) Option[T] {
	return other
}

func (n none[T]) OrElse(f func() Option[T]) Option[T] {
	return f()
}

func (n none[T]) Xor(other Option[T]) Option[T] {
	return other
}

func (n none[T]) String() string {
	return "None"
}

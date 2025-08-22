package option

import "fmt"

type some[T any] struct {
	val T
}

// Some creates a Some variant of Option from the value.
func Some[T any](val T) Option[T] {
	return some[T]{
		val: val,
	}
}

func (s some[T]) IsNone() bool {
	return false
}

func (s some[T]) IsNoneOr(f func(T) bool) bool {
	return f(s.val)
}

func (s some[T]) IsSome() bool {
	return true
}

func (s some[T]) IsSomeAnd(f func(T) bool) bool {
	return f(s.val)
}

func (s some[T]) Expect(_ string) T {
	return s.val
}

func (s some[T]) Unwrap() T {
	return s.val
}

func (s some[T]) UnwrapOr(_ T) T {
	return s.val
}

func (s some[T]) UnwrapOrElse(_ func() T) T {
	return s.val
}

func (s some[T]) UnwrapOrDefault() T {
	return s.val
}

func (s some[T]) Inspect(f func(T)) Option[T] {
	f(s.val)

	return s
}

func (s some[T]) Filter(f func(T) bool) Option[T] {
	if f(s.val) {
		return s
	}

	return none[T]{}
}

func (s some[T]) Or(_ Option[T]) Option[T] {
	return s
}

func (s some[T]) OrElse(_ func() Option[T]) Option[T] {
	return s
}

func (s some[T]) Xor(other Option[T]) Option[T] {
	if other.IsNone() {
		return s
	}

	return none[T]{}
}

func (s some[T]) String() string {
	return fmt.Sprintf("Some(%v)", s.val)
}

package st

import (
	"errors"
	"fmt"
	"reflect"
)

// None creates a None variant of Option.
func None[T any]() *Option[T] {
	var v T

	return &Option[T]{
		ok:  false,
		val: v,
	}
}

// Some creates a Some variant of Option from the value.
func Some[T any](val T) *Option[T] {
	return &Option[T]{
		ok:  true,
		val: val,
	}
}

// OptionOf creates an Option from the given value.
func OptionOf[T any](val T) *Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return None[T]()
	}

	return Some(val)
}

// MapOption maps an Option<T> to Option<U> by applying a function to a contained value (if Some) or returns None
// (if None).
//
// Deprecated: Use Option[T].Map()
func MapOption[T any, U any](opt *Option[T], f func(T) U) *Option[U] {
	if opt.IsNone() {
		return None[U]()
	}

	return Some(f(opt.Unwrap()))
}

// MapOptionOr returns the provided default result (if None), or applies a function to the contained value (if Some).
//
// Deprecated: Use Option[T].MapOr()
func MapOptionOr[T any, U any](opt *Option[T], def U, f func(T) U) U {
	if opt.IsNone() {
		return def
	}

	return f(opt.Unwrap())
}

// MapOptionOrElse computes a default function result (if None), or applies a different function to the contained value
// (if Some).
//
// Deprecated: Use Option[T].MapOrElse()
func MapOptionOrElse[T any, U any](opt *Option[T], factory func() U, f func(T) U) U {
	if opt.IsNone() {
		return factory()
	}

	return f(opt.Unwrap())
}

// And returns None if the option is None, otherwise returns `optb`.
//
// Deprecated: Use Option[T].And()
func And[T any, U any](opt *Option[T], other *Option[U]) *Option[U] {
	if opt.IsNone() {
		return None[U]()
	}

	return other
}

// AndThen returns None if the option is None, otherwise calls `f` with the wrapped value and returns the result.
//
// Deprecated: Use Option[T].AndThen()
func AndThen[T any, U any](opt *Option[T], f func(T) *Option[U]) *Option[U] {
	if opt.IsNone() {
		return None[U]()
	}

	return f(opt.Unwrap())
}

// Option is a type that represents either a value (Some) or not (None).
type Option[T any] struct {
	ok  bool
	val T
}

// IsNone returns true if the option is a None value.
func (o *Option[T]) IsNone() bool {
	return !o.ok
}

// IsNoneOr returns true if the option is a None or the value inside of it matches a predicate.
func (o *Option[T]) IsNoneOr(f func(T) bool) bool {
	if !o.ok {
		return true
	}

	return f(o.val)
}

// IsSome returns true if the option is a Some value.
func (o *Option[T]) IsSome() bool {
	return o.ok
}

// IsSomeAnd returns true if the option is a Some and the value inside of it matches a predicate.
func (o *Option[T]) IsSomeAnd(f func(T) bool) bool {
	if !o.ok {
		return false
	}

	return f(o.val)
}

// Expect returns the contained Some value, consuming the self value. Panics if the value is a None with a custom
// panic message provided by msg.
func (o *Option[T]) Expect(msg string) T {
	if o.ok {
		return o.val
	}

	panic(errors.New(msg))
}

// Unwrap returns the contained Some value, consuming the self value. Panics if the self value equals None.
func (o *Option[T]) Unwrap() T {
	if o.ok {
		return o.val
	}

	panic(errors.New("called `Option.Unwrap()` on a `None` value"))
}

// UnwrapOr returns the contained Some value or a provided default.
func (o *Option[T]) UnwrapOr(def T) T {
	if o.ok {
		return o.val
	}

	return def
}

// UnwrapOrElse returns the contained Some value or computes it from a closure.
func (o *Option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.val
	}

	return f()
}

// UnwrapOrDefault returns the contained Some value or a default.
func (o *Option[T]) UnwrapOrDefault() T {
	if o.ok {
		return o.val
	}

	var def T

	return def
}

// Map maps an Option<T> to Option<U> by applying a function to a contained value (if Some) or returns None
// (if None).
func (o *Option[T]) Map[U any](f func(T) U) *Option[U] {
	if !o.ok {
		return None[U]()
	}

	return Some(f(o.Unwrap()))
}

// MapOr returns the provided default result (if None), or applies a function to the contained value (if Some).
func (o *Option[T]) MapOr[U any](def U, f func(T) U) U {
	if !o.ok {
		return def
	}

	return f(o.Unwrap())
}

// MapOrElse computes a default function result (if None), or applies a different function to the contained value
// (if Some).
func (o *Option[T]) MapOrElse[U any](factory func() U, f func(T) U) U {
	if !o.ok {
		return factory()
	}

	return f(o.Unwrap())
}

// AsOkOr converts an option to a Ok when opt is Some or result.Err when opt is None.
func (o *Option[T]) AsOkOr(err error) *Result[T] {
	if o.ok {
		return Ok[T](o.val)
	}

	return Err[T](err)
}

// AsOkOrElse converts an option to a Ok when opt is Some or result.Err when opt is None.
func (o *Option[T]) AsOkOrElse(f func() error) *Result[T] {
	if o.ok {
		return Ok[T](o.val)
	}

	return Err[T](f())
}

// Inspect calls a function with a reference to the contained value if Some. Returns the original option.
func (o *Option[T]) Inspect(f func(T)) *Option[T] {
	if o.ok {
		f(o.val)
	}

	return o
}

// And returns None if the option is None, otherwise returns `optb`.
func (o *Option[T]) And[U any](other *Option[U]) *Option[U] {
	if !o.ok {
		return None[U]()
	}

	return other
}

// AndThen returns None if the option is None, otherwise calls `f` with the wrapped value and returns the result.
func (o *Option[T]) AndThen[U any](f func(T) *Option[U]) *Option[U] {
	if !o.ok {
		return None[U]()
	}

	return f(o.Unwrap())
}

// Filter returns None if the option is None, otherwise calls predicate with the wrapped value and returns:
//   - Some(t) if predicate returns true (where t is the wrapped value), and
//   - None if predicate returns false.
func (o *Option[T]) Filter(f func(T) bool) *Option[T] {
	if o.ok && f(o.val) {
		return o
	}

	return None[T]()
}

// Or returns the option if it contains a value, otherwise returns optb.
func (o *Option[T]) Or(other *Option[T]) *Option[T] {
	if o.ok {
		return o
	}

	return other
}

// OrElse returns the option if it contains a value, otherwise calls f and returns the result.
func (o *Option[T]) OrElse(f func() *Option[T]) *Option[T] {
	if o.ok {
		return o
	}

	return f()
}

// Xor returns Some if exactly one of self, optb is Some, otherwise returns None.
func (o *Option[T]) Xor(other *Option[T]) *Option[T] {
	if o.ok && other.IsNone() {
		return o
	}

	if !o.ok && other.IsSome() {
		return other
	}

	return None[T]()
}

// Insert inserts value into the option, then returns a mutable reference to it.
//
// If the option already contains a value, the old value is dropped.
//
// See also GetOrInsert, which doesn’t update the value if the option already contains Some.
func (o *Option[T]) Insert(val T) *T {
	o.ok = true
	o.val = val

	return &o.val
}

// GetOrInsert inserts value into the option if it is None, then returns a pointer to the contained value.
//
// See also Insert, which updates the value even if the option already contains Some.
func (o *Option[T]) GetOrInsert(val T) *T {
	if o.ok {
		return &o.val
	}

	o.ok = true
	o.val = val

	return &o.val
}

// GetOrInsertDefault inserts the default value into the option if it is None, then returns a pointer to the
// contained value.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.ok {
		return &o.val
	}

	var def T

	o.ok = true
	o.val = def

	return &o.val
}

// GetOrInsertWith inserts a value computed from f into the option if it is None, then returns a pointer to the
// contained value.
func (o *Option[T]) GetOrInsertWith(f func() T) *T {
	if o.ok {
		return &o.val
	}

	o.ok = true
	o.val = f()

	return &o.val
}

// Take takes the value out of the option, leaving a None in its place.
func (o *Option[T]) Take() *Option[T] {
	if o.ok {
		res := *o

		o.ok = false

		return &res
	}

	return o
}

// TakeIf takes the value out of the option, but only if the predicate evaluates to true on a mutable reference to
// the value.
//
// In other words, replaces self with None if the predicate returns true. This method operates similar to take but
// conditional.
func (o *Option[T]) TakeIf(f func(T) bool) *Option[T] {
	if o.ok && f(o.val) {
		res := *o

		o.ok = false

		return &res
	}

	return None[T]()
}

func (o *Option[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.val)
	}

	return "None"
}

package st

import (
	"errors"
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
	// AsOkOr converts an option to a Ok when opt is Some or result.Err when opt is None.
	AsOkOr(err error) Result[T]
	// AsOkOrElse converts an option to a Ok when opt is Some or result.Err when opt is None.
	AsOkOrElse(f func() error) Result[T]
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
	// Insert inserts value into the option, then returns a mutable reference to it.
	//
	// If the option already contains a value, the old value is dropped.
	//
	// See also GetOrInsert, which doesn’t update the value if the option already contains Some.
	Insert(value T) *T
	// GetOrInsert inserts value into the option if it is None, then returns a pointer to the contained value.
	//
	// See also Insert, which updates the value even if the option already contains Some.
	GetOrInsert(value T) *T
	// GetOrInsertDefault inserts the default value into the option if it is None, then returns a pointer to the
	// contained value.
	GetOrInsertDefault() *T
	// GetOrInsertWith inserts a value computed from f into the option if it is None, then returns a pointer to the
	// contained value.
	GetOrInsertWith(f func() T) *T
	// Take takes the value out of the option, leaving a None in its place.
	Take() Option[T]
	// TakeIf takes the value out of the option, but only if the predicate evaluates to true on a mutable reference to
	// the value.
	//
	// In other words, replaces self with None if the predicate returns true. This method operates similar to take but
	// conditional.
	TakeIf(f func(T) bool) Option[T]

	fmt.Stringer
}

// None creates a None variant of Option.
func None[T any]() Option[T] {
	return newNone[T]()
}

// Some creates a Some variant of Option from the value.
func Some[T any](val T) Option[T] {
	return newSome[T](val)
}

// OptionOf creates an Option from the given value.
func OptionOf[T any](val T) Option[T] {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return newNone[T]()
	}

	return Some(val)
}

// MapOption maps an Option<T> to Option<U> by applying a function to a contained value (if Some) or returns None
// (if None).
func MapOption[T any, U any](opt Option[T], f func(T) U) Option[U] {
	if opt.IsNone() {
		return newNone[U]()
	}

	return Some(f(opt.Unwrap()))
}

// MapOptionOr returns the provided default result (if None), or applies a function to the contained value (if Some).
func MapOptionOr[T any, U any](opt Option[T], def U, f func(T) U) U {
	if opt.IsNone() {
		return def
	}

	return f(opt.Unwrap())
}

// MapOptionOrElse computes a default function result (if None), or applies a different function to the contained value
// (if Some).
func MapOptionOrElse[T any, U any](opt Option[T], factory func() U, f func(T) U) U {
	if opt.IsNone() {
		return factory()
	}

	return f(opt.Unwrap())
}

// And returns None if the option is None, otherwise returns `optb`.
func And[T any, U any](opt Option[T], other Option[U]) Option[U] {
	if opt.IsNone() {
		return newNone[U]()
	}

	return other
}

// AndThen returns None if the option is None, otherwise calls `f` with the wrapped value and returns the result.
func AndThen[T any, U any](opt Option[T], f func(T) Option[U]) Option[U] {
	if opt.IsNone() {
		return newNone[U]()
	}

	return f(opt.Unwrap())
}

func newSome[T any](val T) *option[T] {
	return &option[T]{
		ok:  true,
		val: val,
	}
}

func newNone[T any]() *option[T] {
	var v T

	return &option[T]{
		ok:  false,
		val: v,
	}
}

type option[T any] struct {
	ok  bool
	val T
}

func (o *option[T]) IsNone() bool {
	return !o.ok
}

func (o *option[T]) IsNoneOr(f func(T) bool) bool {
	if !o.ok {
		return true
	}

	return f(o.val)
}

func (o *option[T]) IsSome() bool {
	return o.ok
}

func (o *option[T]) IsSomeAnd(f func(T) bool) bool {
	if !o.ok {
		return false
	}

	return f(o.val)
}

func (o *option[T]) Expect(msg string) T {
	if o.ok {
		return o.val
	}

	panic(errors.New(msg))
}

func (o *option[T]) Unwrap() T {
	if o.ok {
		return o.val
	}

	panic(errors.New("called `Option.Unwrap()` on a `None` value"))
}

func (o *option[T]) UnwrapOr(def T) T {
	if o.ok {
		return o.val
	}

	return def
}

func (o *option[T]) UnwrapOrElse(f func() T) T {
	if o.ok {
		return o.val
	}

	return f()
}

func (o *option[T]) UnwrapOrDefault() T {
	if o.ok {
		return o.val
	}

	var def T

	return def
}

func (o *option[T]) AsOkOr(err error) Result[T] {
	if o.ok {
		return Ok[T](o.val)
	}

	return Err[T](err)
}

func (o *option[T]) AsOkOrElse(f func() error) Result[T] {
	if o.ok {
		return Ok[T](o.val)
	}

	return Err[T](f())
}

func (o *option[T]) Inspect(f func(T)) Option[T] {
	if o.ok {
		f(o.val)
	}

	return o
}

func (o *option[T]) Filter(f func(T) bool) Option[T] {
	if o.ok && f(o.val) {
		return o
	}

	return newNone[T]()
}

func (o *option[T]) Or(other Option[T]) Option[T] {
	if o.ok {
		return o
	}

	return other
}

func (o *option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.ok {
		return o
	}

	return f()
}

func (o *option[T]) Xor(other Option[T]) Option[T] {
	if o.ok && other.IsNone() {
		return o
	}

	if !o.ok && other.IsSome() {
		return other
	}

	return newNone[T]()
}

func (o *option[T]) Insert(val T) *T {
	o.ok = true
	o.val = val

	return &o.val
}

func (o *option[T]) GetOrInsert(val T) *T {
	if o.ok {
		return &o.val
	}

	o.ok = true
	o.val = val

	return &o.val
}

func (o *option[T]) GetOrInsertDefault() *T {
	if o.ok {
		return &o.val
	}

	var def T

	o.ok = true
	o.val = def

	return &o.val
}

func (o *option[T]) GetOrInsertWith(f func() T) *T {
	if o.ok {
		return &o.val
	}

	o.ok = true
	o.val = f()

	return &o.val
}

func (o *option[T]) Take() Option[T] {
	if o.ok {
		res := *o

		o.ok = false

		return &res
	}

	return o
}

func (o *option[T]) TakeIf(f func(T) bool) Option[T] {
	if o.ok && f(o.val) {
		res := *o

		o.ok = false

		return &res
	}

	return newNone[T]()
}

func (o *option[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.val)
	}

	return "None"
}

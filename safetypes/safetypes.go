// Package st implements Rust-style Option[T] and Result[T]
package st

// AsOkOr converts an option to a Ok when opt is Some or result.Err when opt is None.
func AsOkOr[T any](opt Option[T], err error) Result[T] {
	if opt.IsSome() {
		return Ok[T](opt.Unwrap())
	}

	return Err[T](err)
}

// AsOkOrElse converts an option to a Ok when opt is Some or result.Err when opt is None.
func AsOkOrElse[T any](opt Option[T], f func() error) Result[T] {
	if opt.IsSome() {
		return Ok[T](opt.Unwrap())
	}

	return Err[T](f())
}

// AsOptionValue converts a Result to a Some when res is result.Ok or None when res is result.Err.
func AsOptionValue[T any](res Result[T]) Option[T] {
	if res.IsOk() {
		return Some(res.Unwrap())
	}

	var v T

	return OptionOf(v)
}

// AsOptionErr converts a Result to a Some when res is result.Err or None when res is result.Ok.
func AsOptionErr[T any](res Result[T]) Option[error] {
	if res.IsErr() {
		return OptionOf(res.UnwrapErr())
	}

	var v error

	return OptionOf(v)
}

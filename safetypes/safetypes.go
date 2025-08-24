// Package st implements Rust-style Option[T] and Result[T, E]
package st

// AsOkOr converts an option to a Ok when opt is Some or result.Err when opt is None.
func AsOkOr[T any, E error](opt Option[T], err E) Result[T, E] {
	if opt.IsSome() {
		return Ok[T, E](opt.Unwrap())
	}

	return Err[T, E](err)
}

// AsOkOrElse converts an option to a Ok when opt is Some or result.Err when opt is None.
func AsOkOrElse[T any, E error](opt Option[T], f func() E) Result[T, E] {
	if opt.IsSome() {
		return Ok[T, E](opt.Unwrap())
	}

	return Err[T, E](f())
}

// AsOptionValue converts a Result to a Some when res is result.Ok or None when res is result.Err.
func AsOptionValue[T any, E error](res Result[T, E]) Option[T] {
	if res.IsOk() {
		return Some(res.Unwrap())
	}

	var v T

	return OptionOf(v)
}

// AsOptionErr converts a Result to a Some when res is result.Err or None when res is result.Ok.
func AsOptionErr[T any, E error](res Result[T, E]) Option[E] {
	if res.IsErr() {
		return OptionOf(res.UnwrapErr())
	}

	var v E

	return OptionOf(v)
}

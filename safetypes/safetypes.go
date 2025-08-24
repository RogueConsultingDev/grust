// Package st implements Rust-style Option[T] and Result[T]
package st

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

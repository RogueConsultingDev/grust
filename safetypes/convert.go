// Package safetypes implements Rust-style Option[T] and Result[T, E]
package safetypes

import (
	"github.com/mathieu-lemay/grust/safetypes/option"
	"github.com/mathieu-lemay/grust/safetypes/result"
)

// AsOkOr converts an option to a result.Ok when opt is option.Some or result.Err when opt is option.None.
func AsOkOr[T any, E error](opt option.Option[T], err E) result.Result[T, E] {
	if opt.IsSome() {
		return result.Ok[T, E](opt.Unwrap())
	}

	return result.Err[T, E](err)
}

// AsOkOrElse converts an option to a result.Ok when opt is option.Some or result.Err when opt is option.None.
func AsOkOrElse[T any, E error](opt option.Option[T], f func() E) result.Result[T, E] {
	if opt.IsSome() {
		return result.Ok[T, E](opt.Unwrap())
	}

	return result.Err[T, E](f())
}

// AsOptionValue converts a result.Result to an option.Some when res is result.Ok or option.None when res is result.Err.
func AsOptionValue[T any, E error](res result.Result[T, E]) option.Option[T] {
	if res.IsOk() {
		return option.Some(res.Unwrap())
	}

	var v T

	return option.Of(v)
}

// AsOptionErr converts a result.Result to an option.Some when res is result.Err or option.None when res is result.Ok.
func AsOptionErr[T any, E error](res result.Result[T, E]) option.Option[E] {
	if res.IsErr() {
		return option.Of(res.UnwrapErr())
	}

	var v E

	return option.Of(v)
}

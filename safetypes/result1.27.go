//go:build go1.27

package st

// Map maps a Result<T> to Result<U, E> by applying a function to a contained Ok value, leaving an Err value
// untouched.
func (r *Result[T]) Map[U any](f func(T) U) *Result[U] {
	if !r.ok {
		return Err[U](r.UnwrapErr())
	}

	return Ok(f(r.Unwrap()))
}

// MapOr returns the provided default (if Err), or applies a function to the contained value (if Ok).
func (r *Result[T]) MapOr[U any](def U, f func(T) U) U {
	if !r.ok {
		return def
	}

	return f(r.Unwrap())
}

// MapOrElse maps a Result<T> to U by applying fallback function default to a contained Err value, or function
// f to a contained Ok value.
func (r *Result[T]) MapOrElse[U any](factory func() U, mapper func(T) U) U {
	if !r.ok {
		return factory()
	}

	return mapper(r.Unwrap())
}

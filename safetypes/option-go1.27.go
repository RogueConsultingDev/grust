//go:build go1.27

package st

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

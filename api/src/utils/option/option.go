package option

// Option represents an optional value: every Option is either Some and contains a value, or None, and does not.
type Option[T any] struct {
	value  T
	isSome bool
}

// Some creates an Option containing a value.
func Some[T any](value T) Option[T] {
	return Option[T]{value: value, isSome: true}
}

// None creates an Option with no value.
func None[T any]() Option[T] {
	var zero T // Zero value of T
	return Option[T]{value: zero, isSome: false}
}

// IsSome returns true if the Option is a Some value.
func (o Option[T]) IsSome() bool {
	return o.isSome
}

// IsNone returns true if the Option is a None value.
func (o Option[T]) IsNone() bool {
	return !o.isSome
}

// Unwrap returns the contained value.
// It panics if the Option is None.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("called `Option.Unwrap()` on a `None` value")
	}
	return o.value
}

// UnwrapOr returns the contained value or a default if None.
func (o Option[T]) UnwrapOr(defaultValue T) T {
	if o.IsSome() {
		return o.value
	}
	return defaultValue
}

// UnwrapOrElse returns the contained value or computes it from a closure if None.
func (o Option[T]) UnwrapOrElse(f func() T) T {
	if o.IsSome() {
		return o.value
	}
	return f()
}

// Map applies a function to the contained value (if any), and returns a new Option containing the result.
// If the Option is None, it returns None.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsSome() {
		return Some(f(o.value))
	}
	return None[U]()
}

// AndThen (flatMap) applies a function that returns an Option, and returns the result.
// If the Option is None, it returns None.
func AndThen[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if o.IsSome() {
		return f(o.value)
	}
	return None[U]()
}

// Or returns the Option if it contains a value, otherwise returns the alternative Option provided.
func (o Option[T]) Or(opt Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return opt
}

// OrElse returns the Option if it contains a value, otherwise computes it from a closure.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	}
	return f()
}

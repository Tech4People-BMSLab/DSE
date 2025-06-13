package result

// Result represents the result of an operation, containing either a value or an error.
type Result[T any, E error] struct {
    value T
    err   E
    isErr bool
}

// Ok creates a Result containing a value.
func Ok[T any, E error](value T) Result[T, E] {
    var zeroE E // Declare zero value of E
    return Result[T, E]{value: value, err: zeroE, isErr: false}
}

// Err creates a Result containing an error.
func Err[T any, E error](err E) Result[T, E] {
    var zeroT T // Declare zero value of T
    return Result[T, E]{value: zeroT, err: err, isErr: true}
}

// IsOk returns true if the Result is Ok.
func (r Result[T, E]) IsOk() bool {
    return !r.isErr
}

// IsErr returns true if the Result is Err.
func (r Result[T, E]) IsErr() bool {
    return r.isErr
}

// Unwrap returns the contained value.
// It panics if the Result is Err.
func (r Result[T, E]) Unwrap() T {
    if r.IsErr() {
        panic("called `Result.Unwrap()` on an `Err` value")
    }
    return r.value
}

// UnwrapErr returns the contained error.
// It panics if the Result is Ok.
func (r Result[T, E]) UnwrapErr() E {
    if r.IsOk() {
        panic("called `Result.UnwrapErr()` on an `Ok` value")
    }
    return r.err
}

// UnwrapOr returns the contained value or a default if Err.
func (r Result[T, E]) UnwrapOr(value T) T {
    if r.IsOk() {
        return r.value
    }
    return value
}

// Map applies a function to the contained value and returns a new Result.
// If the original Result is Err, it propagates the error.
func Map[T, U any, E error](r Result[T, E], f func(T) U) Result[U, E] {
    if r.IsOk() {
        return Ok[U, E](f(r.value))
    }
    return Err[U, E](r.err)
}

// MapErr applies a function to the contained error and returns a new Result.
// If the original Result is Ok, it propagates the value.
func MapErr[T any, E error, F error](r Result[T, E], f func(E) F) Result[T, F] {
    if r.IsErr() {
        return Err[T, F](f(r.err))
    }
    return Ok[T, F](r.value)
}

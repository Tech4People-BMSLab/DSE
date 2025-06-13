package arrays

func First[T any](s []T) *T {
	if len(s) == 0 {
		return nil
	}

	return &s[0]
}

func Last[T any](s []T) *T {
	if len(s) == 0 {
		return nil
	}

	return &s[len(s)-1]
}

func Find[T any](s []T, predicate func(T) bool) *T {
	for _, v := range s {
		if predicate(v) {
			return &v
		}
	}
	return nil
}

func Contains[T any](s []T, predicate func(T) bool) bool {
	for _, v := range s {
		if predicate(v) {
			return true
		}
	}
	return false
}

func Filter[T any](s []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range s {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

func IndexOf[T any](s []T, predicate func(T) bool) int {
	for i, v := range s {
		if predicate(v) {
			return i
		}
	}
	return -1
}

func Map[T any, U any](s []T, mapper func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = mapper(v)
	}
	return result
}

func Reduce[T any, U any](s []T, reducer func(U, T) U, initial U) U {
	result := initial
	for _, v := range s {
		result = reducer(result, v)
	}
	return result
}

func Reverse[T any](s []T) []T {
	length := len(s)
	result := make([]T, length)
	for i, v := range s {
		result[length-i-1] = v
	}
	return result
}

func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for size < len(s) {
		s, chunks = s[size:], append(chunks, s[0:size:size])
	}
	chunks = append(chunks, s)
	return chunks
}

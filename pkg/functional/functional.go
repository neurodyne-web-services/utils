package functional

// Identity is used to match types for other generic funcs.
func Identity[T any](v T) T {
	return v
}

// Contains - checks whether an object exists in a slice.
func Contains[E comparable](s []E, v E) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}

// Map - generic map function for slices.
func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

// Filter - generic filter function.
func Filter[E any](s []E, f func(E) bool) []E {
	result := []E{}
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Flatten - flattens 2D slice into 1D.
func Flatten[T any](lists [][]T) []T {
	var res []T
	for i := 0; i < len(lists); i++ {
		res = append(res, lists[i]...)
	}
	return res
}

// Fold - generic reduce function.
func Fold[E any](s []E, init E, f func(curr, next E) E) E {
	cur := init
	for _, v := range s {
		cur = f(cur, v)
	}
	return cur
}

// FoldMap - maps each element in array and folds result into a single value.
func FoldMap[T, U any](ts []T, init T, mapFunc func(T) U, reduceFunc func(curr, next U) U) U {
	cur := mapFunc(init)
	for _, v := range ts {
		cur = reduceFunc(cur, mapFunc(v))
	}
	return cur
}

// MapMapValues - maps over generic map values.
func MapMapValues[K comparable, V, U any](ts map[K]V, f func(V) U) []U {
	us := make([]U, len(ts))
	for _, v := range ts {
		us = append(us, f(v))
	}
	return us
}

// MapMapKeys - maps over generic map keys.
func MapMapKeys[K comparable, V, U any](ts map[K]V, f func(K) U) []U {
	us := make([]U, len(ts))
	for k := range ts {
		us = append(us, f(k))
	}
	return us
}

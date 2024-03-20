package slice

import (
	"reflect"
)

const (
	defaultSliceLen = 8
)

//-----------------------------------------------------------------------------------------------------------
// Slice utils
//-----------------------------------------------------------------------------------------------------------

type GenericSliceFunc func(a, b any) []any

// Taken from https://github.com/juliangruber/go-intersect
// Intersect - finds an intersect of two slices of any type.
func Intersect(a, b any) []any {
	set := make([]any, 0, defaultSliceLen)
	hash := make(map[any]bool, defaultSliceLen)

	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		if _, found := hash[el]; found {
			set = append(set, el)
		}
	}

	return set
}

// Diff - finds an intersect of two slices of any type
// Diff with hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2).
func Diff(a, b any) []any {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	set := make([]any, 0, defaultSliceLen)
	hash := make(map[any]bool, defaultSliceLen)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		hash[el] = true
	}

	for i := 0; i < bv.Len(); i++ {
		el := bv.Index(i).Interface()
		if _, found := hash[el]; !found {
			set = append(set, el)
		}
	}

	return set
}

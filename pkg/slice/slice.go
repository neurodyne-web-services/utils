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

type GenericSliceFunc func(a, b interface{}) []interface{}

// DiffSlices - compares two string slices. False if equal
func EqualSlices(a, b interface{}) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	if av.Len() != bv.Len() {
		return false
	}

	for i := 0; i < av.Len(); i++ {
		at := av.Index(i).Interface()
		bt := bv.Index(i).Interface()
		if at != bt {
			return false
		}
	}

	return true
}

// Taken from https://github.com/juliangruber/go-intersect
// SliceIntersect - finds an intersect of two slices of any type
func SliceIntersect(a, b interface{}) []interface{} {
	set := make([]interface{}, 0, defaultSliceLen)
	hash := make(map[interface{}]bool, defaultSliceLen)

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

// SliceDiff - finds an intersect of two slices of any type
// SliceDiff with hash has complexity: O(n * x) where x is a factor of hash function efficiency (between 1 and 2)
func SliceDiff(a, b interface{}) []interface{} {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	set := make([]interface{}, 0, defaultSliceLen)
	hash := make(map[interface{}]bool, defaultSliceLen)

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

func Contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}

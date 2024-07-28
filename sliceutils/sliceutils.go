package sliceutils

import (
	"math/rand"
	"reflect"
)



func Shuffle(slice interface{}) {
	switch v := slice.(type) {
	case []map[string]interface{}:
		rand.Shuffle(len(v), func(i, j int) { v[i], v[j] = v[j], v[i] })
	case []string:
		rand.Shuffle(len(v), func(i, j int) { v[i], v[j] = v[j], v[i] })
	default:
		panic("unsupported type")
	}
}

// RandomSubset returns a random subset of the given size from the slice.
func RandomSubset(slice interface{}, subset int) interface{} {
	switch v := slice.(type) {
	case []map[string]interface{}:
		if subset > len(v) {
			return v
		}
		Shuffle(v)
		return v[:subset]
	case []string:
		if subset > len(v) {
			return v
		}
		Shuffle(v)
		return v[:subset]
	default:
		panic("unsupported type")
	}
}


func IsSlice(variable interface{}) bool {
	return reflect.TypeOf(variable).Kind() == reflect.Slice
}
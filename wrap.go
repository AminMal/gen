package gen

import (
	"reflect"
)

type WrappedGen struct {
	tpe reflect.Type
	vg  Gen[reflect.Value]
}

type valueGen[T any] struct {
	underlying Gen[T]
}

func (t *valueGen[T]) Generate() reflect.Value {
	return reflect.ValueOf(t.underlying.Generate())
}

func (t *valueGen[T]) GenerateN(n uint) []reflect.Value {
	result := make([]reflect.Value, n, n)
	for i := uint(0); i < n; i++ {
		result[i] = t.Generate()
	}
	return result
}

func Wrap[T any](g Gen[T]) *WrappedGen {
	return &WrappedGen{reflect.TypeOf(*new(T)), &valueGen[T]{g}}
}

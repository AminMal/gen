package gen

import (
	"reflect"
)

// WrappedGen basically wraps `Gen`s to provide a generator that works with `reflect.Value`
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

// Wrap wraps around a `Gen` and returns a *WrappedGen.
func Wrap[T any](g Gen[T]) *WrappedGen {
	return &WrappedGen{reflect.TypeOf(*new(T)), &valueGen[T]{g}}
}

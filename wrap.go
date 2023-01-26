package gen

import (
	"fmt"
	"reflect"
	"strings"
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

func getFunctionSignature(ft reflect.Type) string {
	ins := []string{}
	outs := []string{}
	for i := 0; i < ft.NumIn(); i++ {
		in := ft.In(i).Name()
		if in == "" {
			in = "any"
		}
		ins = append(ins, in)
	}
	for i := 0; i < ft.NumOut(); i++ {
		out := ft.Out(i).Name()
		if out == "" {
			out = "any"
		}
		outs = append(outs, out)
	}

	return fmt.Sprintf("(%s) => (%s)", strings.Join(ins, ", "), strings.Join(outs, ", "))
}

func notInferrable(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Func:
		return fmt.Errorf("cannot infer functions yet: %s", getFunctionSignature(t))
	default:
		return fmt.Errorf("cannot infer `%s` yet", t)
	}
}

func Infer[T any](valueGenerators ...*WrappedGen) (Gen[T], error) {
	tpe := reflect.TypeOf(*new(T))

	valueGeneratorsByType := make(map[reflect.Type]*WrappedGen)

	for _, vg := range valueGenerators {
		valueGeneratorsByType[vg.tpe] = vg
	}

	for i := 0; i < tpe.NumField(); i++ {
		if tpe.Field(i).Type.Kind() == reflect.Func {
			return nil, notInferrable(tpe.Field(i).Type)
		}
	}

	if tpe.Kind() == reflect.Func {
		return nil, notInferrable(tpe)
	}
	return &adhocGen[T]{valueGeneratorsByType}, nil
}

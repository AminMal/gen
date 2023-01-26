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

type adhocGen[T any] struct {
	generatorsByType map[reflect.Type]*WrappedGen
}

func (g *adhocGen[T]) Generate() T {
	actual := new(T)
	tpe := reflect.TypeOf(actual)
	v := reflect.ValueOf(actual).Elem()

	concrete := tpe.Elem()

	for i := 0; i < concrete.NumField(); i++ {
		fieldType := concrete.Field(i)
		if g, found := g.generatorsByType[fieldType.Type]; found {
			if v.Field(i).CanSet() {
				v.Field(i).Set(g.vg.Generate())
				continue
			}
		} else {
			if fieldValue, ok := sizedValue(fieldType.Type, complexSize); ok {
				if v.Field(i).CanSet() {
					v.Field(i).Set(fieldValue)
				}
			} else {
				panic(notImplemented(fieldType.Type))
			}
		}
	}
	return *actual
}

func (g *adhocGen[T]) GenerateN(n uint) []T {
	result := make([]T, n, n)
	for i := uint(0); i < n; i++ {
		result[i] = g.Generate()
	}
	return result
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

func notImplemented(t reflect.Type) error {
	switch t.Kind() {
	case reflect.Func:
		return fmt.Errorf("cannot construct functions yet: %s", getFunctionSignature(t))
	default:
		return fmt.Errorf("cannot construct this type yet (it's probably a function or anonymous): `%s`", t)
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
			return nil, notImplemented(tpe.Field(i).Type)
		}
	}

	if tpe.Kind() == reflect.Func {
		return nil, notImplemented(tpe)
	}
	return &adhocGen[T]{valueGeneratorsByType}, nil
}

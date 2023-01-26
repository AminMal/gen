package gen

import (
	"math/rand"
	"reflect"
	"fmt"
	"strings"
)

var complexSize = 50

func sizedValue(t reflect.Type, size int) (value reflect.Value, ok bool) {
	v := reflect.New(t).Elem()
	switch concrete := t; concrete.Kind() {
	case reflect.Bool:
		v.SetBool(ArbitraryInt.Generate()&1 == 0)
	case reflect.Float32:
		v.SetFloat(float64(ArbitraryFloat32.Generate()))
	case reflect.Float64:
		v.SetFloat(ArbitraryFloat64.Generate())
	case reflect.Complex64:
		v.SetComplex(complex(ArbitraryFloat64.Generate(), ArbitraryFloat64.Generate()))
	case reflect.Complex128:
		v.SetComplex(complex(ArbitraryFloat64.Generate(), ArbitraryFloat64.Generate()))
	case reflect.Int16:
		v.SetInt(ArbitraryInt64.Generate())
	case reflect.Int32:
		v.SetInt(ArbitraryInt64.Generate())
	case reflect.Int64:
		v.SetInt(ArbitraryInt64.Generate())
	case reflect.Int8:
		v.SetInt(ArbitraryInt64.Generate())
	case reflect.Int:
		v.SetInt(ArbitraryInt64.Generate())
	case reflect.Uint16:
		v.SetUint(uint64(ArbitraryUint16.Generate()))
	case reflect.Uint32:
		v.SetUint(uint64(ArbitraryUint32.Generate()))
	case reflect.Uint64:
		v.SetUint(ArbitraryUint64.Generate())
	case reflect.Uint8:
		v.SetUint(uint64(ArbitraryUint8.Generate()))
	case reflect.Uint:
		v.SetUint(uint64(ArbitraryUint.Generate()))
	case reflect.Uintptr:
		v.SetUint(ArbitraryUint64.Generate())
	case reflect.Map:
		numElems := Between(0, size).Generate()
		v.Set(reflect.MakeMap(concrete))
		for i := 0; i < numElems; i++ {
			key, ok1 := sizedValue(concrete.Key(), size)
			value, ok2 := sizedValue(concrete.Elem(), size)
			if !ok1 || !ok2 {
				return reflect.Value{}, false
			}
			v.SetMapIndex(key, value)
		}
	case reflect.Pointer:
		if Between(0, size).Generate() == 0 {
			v.Set(reflect.Zero(concrete)) // Generate nil pointer.
		} else {
			elem, ok := sizedValue(concrete.Elem(), size)
			if !ok {
				return reflect.Value{}, false
			}
			v.Set(reflect.New(concrete.Elem()))
			v.Elem().Set(elem)
		}
	case reflect.Slice:
		numElems := Between(0, size).Generate()
		sizeLeft := size - numElems
		v.Set(reflect.MakeSlice(concrete, numElems, numElems))
		for i := 0; i < numElems; i++ {
			elem, ok := sizedValue(concrete.Elem(), sizeLeft)
			if !ok {
				return reflect.Value{}, false
			}
			v.Index(i).Set(elem)
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			elem, ok := sizedValue(concrete.Elem(), size)
			if !ok {
				return reflect.Value{}, false
			}
			v.Index(i).Set(elem)
		}
	case reflect.String:
		numChars := rand.Intn(complexSize)
		codePoints := make([]rune, numChars)
		for i := 0; i < numChars; i++ {
			codePoints[i] = rune(Between(0, 0x10ffff).Generate())
		}
		v.SetString(string(codePoints))
	case reflect.Struct:
		n := v.NumField()
		// Divide sizeLeft evenly among the struct fields.
		sizeLeft := size
		if n > sizeLeft {
			sizeLeft = 1
		} else if n > 0 {
			sizeLeft /= n
		}
		for i := 0; i < n; i++ {
			elem, ok := sizedValue(concrete.Field(i).Type, sizeLeft)
			if !ok {
				return reflect.Value{}, false
			}
			v.Field(i).Set(elem)
		}
	default:
		return reflect.Value{}, false
	}

	return v, true
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
				panic(notInferrable(fieldType.Type))
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


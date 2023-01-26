package gen

import "fmt"

type Gen[T any] interface {
	Generate() T
	GenerateN(n uint) []T
}

type only[T any] struct {
	value T
}

func (o *only[T]) Generate() T { return o.value }

func (o *only[T]) GenerateN(n uint) []T {
	values := make([]T, n, n)
	for i := uint(0); i < n; i++ {
		values[i] = o.value
	}
	return values
}

func Only[T any](value T) Gen[T] {
	return &only[T]{value}
}

type oneOf[T any] struct {
	choices    []T
	numChoices int
}

func (o *oneOf[T]) Generate() T {
	return o.choices[random.Intn(o.numChoices)]
}

func (o *oneOf[T]) GenerateN(n uint) []T {
	values := make([]T, n, n)
	for i := uint(0); i < n; i++ {
		values[i] = o.Generate()
	}
	return values
}

func OneOf[T any](values ...T) Gen[T] {
	return &oneOf[T]{values, len(values)}
}

type between[T Numeric] struct {
	min, max T
}

func (r *between[T]) Generate() T {
	// todo, the below line causes subtraction overflow, fix it
	switch diff := any(r.max - r.min).(type) {
	case uint8:
		return any(randUint8(diff)).(T) + r.min
	case uint16:
		return any(randUint16(diff)).(T) + r.min
	case uint32:
		return any(randUint32(diff)).(T) + r.min
	case uint64:
		return any(randUint64(diff)).(T) + r.min
	case uint:
		return any(randUint(diff)).(T) + r.min
	case int8:
		return any(randInt8(diff)).(T) + r.min
	case int16:
		return any(randInt16(diff)).(T) + r.min
	case int32:
		return any(randInt32(diff)).(T) + r.min
	case int64:
		return any(randInt64(diff)).(T) + r.min
	case int:
		return any(randInt(diff)).(T) + r.min
	case float32:
		return any(randFloat32(diff)).(T) + r.min
	case float64:
		return any(randFloat64(diff)).(T) + r.min
	default:
		panic(fmt.Errorf("match error: unrecognized Numeric type %t", diff))
	}
}

func (r *between[T]) GenerateN(n uint) []T {
	res := make([]T, n)
	for i := uint(0); i < n; i++ {
		res[i] = r.Generate()
	}
	return res
}

func Between[T Numeric](min, max T) Gen[T] {
	actualMin := numericMin(min, max)
	actualMax := numericMax(min, max)
	return &between[T]{actualMin, actualMax}
}

package gen

import "fmt"

// Gen describes how to generate a value of a specific type `T`.
// The behavior of the Gen only depends on the structs implementing it.
type Gen[T any] interface {
	// Generate generates a single value of type `T`
	Generate() T
	// GenerateN generates `n` values of type `T`
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

// Only can generate only the value it's given.
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

// OneOf picks out a value among those values that it's given.
// If the values contain only one element, it returns an Only generator.
func OneOf[T any](values ...T) Gen[T] {
	if len(values) == 1 {
		return Only(values[0])
	}
	return &oneOf[T]{values, len(values)}
}

type between[T Numeric] struct {
	min, max T
}

func (r *between[T]) Generate() T {
	// todo, the below line causes subtraction overflow, fix it
	switch diff := any(r.max - r.min).(type) {
	case uint8:
		return T(randUint8(diff) + uint8(r.min))
	case uint16:
		return T(randUint16(diff) + uint16(r.min))
	case uint32:
		return T(randUint32(diff) + uint32(r.min))
	case uint64:
		return T(randUint64(diff) + uint64(r.min))
	case uint:
		return T(randUint(diff) + uint(r.min))
	case int8:
		return T(randInt8(diff) + int8(r.min))
	case int16:
		return T(randInt16(diff) + int16(r.min))
	case int32:
		return T(randInt32(diff) + int32(r.min))
	case int64:
		return T(randInt64(diff) + int64(r.min))
	case int:
		return T(randInt(diff) + int(r.min))
	case float32:
		return T(randFloat32(diff) + float32(r.min))
	case float64:
		return T(randFloat64(diff) + float64(r.min))
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

// Between generates values within the given range.
// The order of the parameters doesn't actually matter, but it's more convenient to pass them properly.
// If max equals min, it returns an Only generator
func Between[T Numeric](min, max T) Gen[T] {
	if min == max {
		return Only(min)
	}
	actualMin := numericMin(min, max)
	actualMax := numericMax(min, max)
	return &between[T]{actualMin, actualMax}
}

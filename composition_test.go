package gen

import (
	"fmt"
	"testing"
)

type genWithPanicSideEffect[T any] struct{}

func (genWithPanicSideEffect[T]) Generate() T {
	panic("I'm supposed to panic")
}

func (genWithPanicSideEffect[T]) GenerateN(n uint) []T {
	panic(fmt.Sprintf("I'm supposed to panic %d times", n))
}

type Person struct {
	Name string
	Age  int
}

func TestLazyGenBeingAnEffect(t *testing.T) {

	var nameGen Gen[string] = genWithPanicSideEffect[string]{}
	var ageGen Gen[int] = genWithPanicSideEffect[int]{}

	_ = UsingGen(nameGen, func(name string) Gen[Person] {
		return Using(ageGen, func(age int) Person {
			return Person{name, age}
		})
	})

	// if we reach this point and the program has not pannicked, it means that `genWithPanicSideEffect`'s Generate function
	// has not been called even once. Otherwise, test is failed with panic
}

func TestComposedOnlyGenBeingOnly(t *testing.T) {
	nameGen := Only("John")
	ageGen := Only(42)

	personGen := UsingGen(nameGen, func(name string) Gen[Person] {
		return Using(ageGen, func(age int) Person {
			return Person{name, age}
		})
	})

	expected := Person{"John", 42}

	for _, person := range personGen.GenerateN(100) {
		if person != expected {
			t.Errorf("lazyGen with `Only` as it's base generators did not produce the expected result. expected: %v, got: %v", expected, person)
		}
	}
}

package gen

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func TestOnly(t *testing.T) {
	expected := 2
	g := Only(expected)

	for _, value := range GenerateN(g, 100) {
		if value != expected {
			t.Fatalf("only returned unexpected value. expected: %d, actual: %d", expected, value)
		}
	}
}

func TestOneOf(t *testing.T) {
	possibilities := []string{"Margot", "John", "Jack", "Peter", "Bob", "Anne", "Marry"}
	g := OneOf(possibilities...)

	possibilitiesLookup := make(map[string]struct{})

	for _, p := range possibilities {
		possibilitiesLookup[p] = struct{}{}
	}

	for _, value := range GenerateN(g, 100) {
		if _, exists := possibilitiesLookup[value]; !exists {
			t.Fatalf("OneOf returned %s which was not in the possibilities", value)
		}
	}
}

func isBetween(num, actualMin, actualMax int) bool {
	return num >= actualMin && num <= actualMax
}

func TestBetween(t *testing.T) {
	min := 100000  // this is intentional
	max := -790832 // this is intentional
	g := Between(min, max)

	actualMin := max
	actualMax := min

	for _, value := range GenerateN(g, 100) {
		if !isBetween(value, actualMin, actualMax) {
			t.Fatalf("%d is not actually between %d and %d", value, actualMin, actualMax)
		}
	}
}

func TestBetweenWithOnePossibility(t *testing.T) {
	amount := 10
	g := Between(amount, amount)

	for _, value := range GenerateN(g, 100) {
		if value != amount {
			t.Fatalf("Between with the same arguments did not act as Only")
		}
	}
}

func TestOneOfWithOnePossibility(t *testing.T) {
	type Human struct{ Name string }
	onlyPossibility := Human{"John"}

	choices := OneOf(onlyPossibility)
	for _, value := range GenerateN(choices, 100) {
		if value != onlyPossibility {
			t.Fatal("OneOf did not act as Only when only one possibility exists")
		}
	}
}

type quickOnly[T any] struct {
	value T
}

func (qo *quickOnly[T]) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(qo.value)
}

func BenchmarkOnly(b *testing.B) {
	value := 42
	o := Only(value)
	var qo quick.Generator = &quickOnly[int]{value}

	b.Run("gen-only", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			o.Generate()
		}
	})

	b.Run("quick-only", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qo.Generate(random, 50).Int()
		}
	})
}

type quickBetween[T Numeric] struct {
	start, end int
}

// This is a much-simpler implementation of Between, so it makes sense if it can compete with gen's Between
func (qb *quickBetween[T]) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(
		rand.Intn(int(qb.end - qb.start)) + qb.start,
	)
}

func BenchmarkBetween(b *testing.B) {

	start := 1
	end := 100000000
	
	g := Between(start, end)
	var qb quick.Generator = &quickBetween[int]{start, end}

	b.Run("gen-between", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Generate()
		}
	})

	b.Run("quick-between", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qb.Generate(random, 50).Int()
		}
	})
}

type quickOneOf[T any] struct {
	slice []T
	length int
}

func (qo *quickOneOf[T]) Generate(rand *rand.Rand, size int) reflect.Value {
	index := rand.Intn(qo.length)
	return reflect.ValueOf(qo.slice[index])
}

func BenchmarkOneOf(b *testing.B) {
	slice := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		slice[i] = i
	}

	g := OneOf(slice)
	var qo quick.Generator = &quickOneOf[int]{slice, len(slice)}

	b.Run("gen-one-of", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Generate()
		}
	})

	b.Run("quick-one-of", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			qo.Generate(random, 50).Int()
		}
	})
}

type TestPerson struct {
	Name 		string
	Surname 	string
	Age 		int
}

func BenchmarkComposition(b *testing.B) {

	nameChoices := []string{"John", "Bob", "Alice", "Anne", "Sherlock", "Jack", "Ross", "Brian"}
	surnameChoices := []string{"Cole", "Dylan", "Chains", "Marry", "Holmes", "Black", "Geller", "UNKNOWN"}
	age := 20

	g := UsingGen(OneOf(nameChoices...), func(name string) Gen[TestPerson] {
		return UsingGen(OneOf(surnameChoices...), func (surname string) Gen[TestPerson] {
			return Using(Only(age), func(age int) TestPerson {
				return TestPerson{name, surname, age}
			})
		})
	})

	inferedG, _ := Infer[TestPerson](Wrap(OneOf(nameChoices)), Wrap(Only(age)))
	
	b.Run("gen-composition-functional", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Generate()
		}
	})

	b.Run("gen-infered-composition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			inferedG.Generate()
		}
	})
}

type personGen struct {
	nameGen, surnameGen Gen[string]
	ageGen 	Gen[int]
}

func (pg *personGen) Generate() TestPerson {
	return TestPerson{pg.nameGen.Generate(), pg.surnameGen.Generate(), pg.ageGen.Generate()}
}

type quickPersonGen struct {
	nameGen, surnameGen, ageGen quick.Generator
}

func (qpg quickPersonGen) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(
		TestPerson{qpg.nameGen.Generate(random, 50).String(), qpg.surnameGen.Generate(random, 50).String(), int(qpg.ageGen.Generate(random, 50).Int())},
	)
}

func BenchmarkComposition2(b *testing.B) {
	names := []string{"John", "Jack", "Beth", "Anne", "Freddie", "Ross", "Rachel", "Tom"}
	surnames := []string{"Smith", "Hart", "Lauren", "Marry", "Simpson", "Geller", "Hanks", "Mercury"}
	age := 42

	g := &personGen{OneOf(names...), OneOf(surnames...), Only(age)}
	qg := &quickPersonGen{&quickOneOf[string]{names, len(names)}, &quickOneOf[string]{surnames, len(surnames)}, &quickOnly[int]{age}}

	b.Run("gen-composition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			g.Generate()
		}
	})

	b.Run("quick-composition", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = qg.Generate(random, 50).Interface().(TestPerson)
		}
	})
}

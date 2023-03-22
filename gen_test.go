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
		rand.Intn(int(qb.end-qb.start)) + qb.start,
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
	slice  []T
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
	Name    string
	Surname string
	Age     int
}

func BenchmarkComposition(b *testing.B) {

	nameChoices := []string{"John", "Bob", "Alice", "Anne", "Sherlock", "Jack", "Ross", "Brian"}
	surnameChoices := []string{"Cole", "Dylan", "Chains", "Marry", "Holmes", "Black", "Geller", "UNKNOWN"}
	age := 20

	g := FlatMap(OneOf(nameChoices...), func(name string) Gen[TestPerson] {
		return FlatMap(OneOf(surnameChoices...), func(surname string) Gen[TestPerson] {
			return Map(Only(age), func(age int) TestPerson {
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
	ageGen              Gen[int]
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

type Programmer struct {
	Name       string
	Surname    string
	GithubUrl  string
	FavLang    string
	Origin     string
	Age        int
	Experiance int
}

type programmerGen struct {
	nameGen, surnameGen, gitGen, langGen, originGen Gen[string]
	ageGen, experienceGen                           Gen[int]
}

func (pg *programmerGen) Generate() Programmer {
	return Programmer{
		pg.nameGen.Generate(),
		pg.surnameGen.Generate(),
		pg.gitGen.Generate(),
		pg.langGen.Generate(),
		pg.originGen.Generate(),
		pg.ageGen.Generate(),
		pg.experienceGen.Generate(),
	}
}

type quickProgrammerGen struct {
	nameGen, surnameGen, gitGen, langGen, originGen, ageGen, experienceGen quick.Generator
}

func (qpg quickProgrammerGen) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(
		Programmer{
			qpg.nameGen.Generate(random, 50).String(),
			qpg.surnameGen.Generate(random, 50).String(),
			qpg.gitGen.Generate(random, 50).String(),
			qpg.langGen.Generate(random, 50).String(),
			qpg.originGen.Generate(random, 50).String(),
			int(qpg.ageGen.Generate(random, 50).Int()),
			int(qpg.experienceGen.Generate(random, 50).Int()),
		},
	)
}

func BenchmarkCompositionMultipleFields(b *testing.B) {
	nameGen := OneOf("John", "Jack", "Beth", "Anne", "Freddie", "Ross", "Rachel", "Tom")
	surnameGen := OneOf("Smith", "Hart", "Lauren", "Marry", "Simpson", "Geller", "Hanks", "Mercury")
	gitGen := Only("https://github.com/AminMal")
	langGen := OneOf("Scala", "Rust", "Go", "Java", "Python")
	originGen := OneOf("USA", "Iran", "Spain", "France")
	ageGen := Between(16, 79)
	experienceGen := Between(1, 8)

	var directGen Gen[Programmer] = &programmerGen{nameGen, surnameGen, gitGen, langGen, originGen, ageGen, experienceGen}

	qg := quickProgrammerGen{
		nameGen:       &quickOneOf[string]{[]string{"John", "Jack", "Beth", "Anne", "Freddie", "Ross", "Rachel", "Tom"}, 8},
		surnameGen:    &quickOneOf[string]{[]string{"Smith", "Hart", "Lauren", "Marry", "Simpson", "Geller", "Hanks", "Mercury"}, 8},
		gitGen:        &quickOnly[string]{"https://github.com/AminMal"},
		langGen:       &quickOneOf[string]{[]string{"Scala", "Rust", "Go", "Java", "Python"}, 5},
		originGen:     &quickOneOf[string]{[]string{"USA", "Iran", "Spain", "France"}, 4},
		ageGen:        &quickBetween[int]{16, 79},
		experienceGen: &quickBetween[int]{1, 8},
	}

	var functional1 Gen[Programmer] = FlatMap(nameGen, func(name string) Gen[Programmer] {
		return FlatMap(surnameGen, func(surname string) Gen[Programmer] {
			return FlatMap(gitGen, func(git string) Gen[Programmer] {
				return FlatMap(langGen, func(lang string) Gen[Programmer] {
					return FlatMap(originGen, func(origin string) Gen[Programmer] {
						return FlatMap(ageGen, func(age int) Gen[Programmer] {
							return Map(experienceGen, func(experience int) Programmer {
								return Programmer{
									name, surname, git, lang, origin, age, experience,
								}
							})
						})
					})
				})
			})
		})
	})

	var functional2 Gen[Programmer] = Map7(
		nameGen, surnameGen, gitGen, langGen, originGen, ageGen, experienceGen,
		func(name, surname, git, lang, origin string, age, experience int) Programmer {
			return Programmer{name, surname, git, lang, origin, age, experience}
		},
	)

	b.Run("gen-composition-7-fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			directGen.Generate()
		}
	})

	b.Run("gen-composition-map/flatmap-7-fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			functional1.Generate()
		}
	})

	b.Run("gen-composition-mapN-7-fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			functional2.Generate()
		}
	})

	b.Run("quick-composition-7-fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = qg.Generate(random, 50).Interface().(Programmer)
		}
	})

}

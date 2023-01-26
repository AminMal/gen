package gen

import (
	"testing"
)

func TestOnly(t *testing.T) {
	expected := 2
	g := Only(expected)

	for _, value := range g.GenerateN(100) {
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

	for _, value := range g.GenerateN(100) {
		if _, exists := possibilitiesLookup[value]; !exists {
			t.Fatalf("OneOf returned %s which was not in the possibilities", value)
		}
	}
}

func isBetween(num, actualMin, actualMax int) bool {
	return num >= actualMin && num <= actualMax
}

func TestBetween(t *testing.T) {
	min := 100000 // this is intentional
	max := -790832 // this is intentional
	g := Between(min, max)

	actualMin := max
	actualMax := min

	for _, value := range g.GenerateN(100) {
		if !isBetween(value, actualMin, actualMax) {
			t.Fatalf("%d is not actually between %d and %d", value, actualMin, actualMax)
		}
	}
}

func TestBetweenWithOnePossibility(t *testing.T) {
	amount := 10
	g := Between(amount, amount)

	for _, value := range g.GenerateN(100) {
		if value != amount {
			t.Fatalf("Between with the same arguments did not act as Only")
		}
	}
}

func TestOneOfWithOnePossibility(t *testing.T) {
	type Human struct { Name string }
	onlyPossibility := Human{"John"}

	choices := OneOf(onlyPossibility)
	for _, value := range choices.GenerateN(100) {
		if value != onlyPossibility {
			t.Fatal("OneOf did not act as Only when only one possibility exists")
		}
	}
}

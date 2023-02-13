package gen

import (
	"reflect"
	"testing"
)

func TestSeqState(t *testing.T) {

	from := 1
	to := 100
	step := 1

	s := Sequential(from, to, step)

	take := 59
	_ = GenerateN(s, uint(take))

	expected := take + 1
	actual := s.Generate()
	if actual != expected {
		t.Errorf("sequential generator failed to hold the state, expected: `%d`, got `%d`", expected, actual)
	}
}

func TestSeqForZeroStep(t *testing.T) {
	from := 1
	to := 100
	step := 0

	s := Sequential(from, to, step)
	o := Only(from)
	take := 50

	if !reflect.DeepEqual(GenerateN(s, uint(take)), GenerateN(o, uint(take))) {
		t.Errorf("sequential generator with 0 step did not generate only one value")
	}
}

func TestSeqForStepNotInTheSameDirectionAsStartToEndVector(t *testing.T) {
	from := 1
	to := 100
	step := -2

	s := Sequential(from, to, step)
	o := Only(from)

	take := 50

	if !reflect.DeepEqual(GenerateN(s, uint(take)), GenerateN(o, uint(take))) {
		t.Errorf("sequential generator with step not in the same direction as start to end vector did not generate only one value")
	}
}

// todo, add propert-based testings for seq state management

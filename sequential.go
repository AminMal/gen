package gen

import "time"

type seq[T Numeric] struct {
	from, to, step, current T
}

func (s *seq[T]) Generate() T {
	current := s.current
	if current+s.step > s.to {
		s.current = s.from + (s.step - (s.to - current) - 1)
		return current
	}
	s.current += s.step
	return current
}

func (s *seq[T]) GenerateN(n uint) []T {
	res := make([]T, n, n)
	for i := uint(0); i < n; i++ {
		// todo, inline Generate here?
		res[i] = s.Generate()
	}
	return res
}

// Sequential is a sequential generator that holds the current state of the generator.
// It will generate numerics, between `from` and `to` (inclusive), with the given `step` size.
func Sequential[T Numeric](from, to, step T) Gen[T] {
	if (from > to && step < 0) || (from < to && step > 0) {
		return &seq[T]{from, to, step, from}
	}

	// `from` equals `to` or `step` is zero
	// Or step is not as the same direction as `from` to `to`, so we can only return `from`
	return Only(from)
}

type timeSeq struct {
	from, to, current time.Time
	step              time.Duration
}

func (ts *timeSeq) Generate() time.Time {
	current := ts.current
	if current.Add(ts.step).After(ts.to) {
		ts.current = ts.from.Add(ts.step - (ts.to.Sub(current)))
		return current
	}
	ts.current = ts.current.Add(ts.step)
	return current
}

func (ts *timeSeq) GenerateN(n uint) []time.Time {
	res := make([]time.Time, n, n)
	for i := uint(0); i < n; i++ {
		res[i] = ts.Generate()
	}
	return res
}

// TimeSeq is a sequential time generator, it generates `time.Time`s within the given range and step.
func TimeSeq(from, to time.Time, step time.Duration) Gen[time.Time] {
	if (from.After(to) && step < 0) || (from.Before(to) && step > 0) {
		return &timeSeq{from, to, from, step}
	}

	// `from` equals `to` or `step` is zero
	// Or step is not as the same direction as `from` to `to`, so we can only return `from`
	return Only(from)
}

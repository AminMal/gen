package gen

import "time"

type timeBetween struct {
	start       time.Time
	durationGen Gen[int64]
}

func (t timeBetween) Generate() time.Time {
	newDuration := t.durationGen.Generate()
	return t.start.Add(time.Duration(newDuration))
}

func (t timeBetween) GenerateN(n uint) []time.Time {
	res := make([]time.Time, n)
	for i := uint(0); i < n; i++ {
		res[i] = t.Generate()
	}
	return res
}

func TimeBetween(start time.Time, end time.Time) Gen[time.Time] {
	if start.Equal(end) {
		return Only(start)
	}
	actualStart := start
	if start.After(end) {
		actualStart = end
	}
	actualEnd := end
	if end.Before(start) {
		actualEnd = start
	}

	dur := actualEnd.Sub(actualStart)
	return timeBetween{actualStart, Between(int64(0), int64(dur))}
}

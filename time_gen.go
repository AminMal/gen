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

// TimeBetween is a generator for `time.Time` that will generate random `time.Time`s between the given start and end.
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

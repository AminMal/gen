package gen

func GenerateN[T any](g Gen[T], n uint) []T {
	if g == nil { return nil }
	res := make([]T, n, n)
	for i := uint(0); i < n; i++ {
		res[i] = g.Generate()
	}
	return res
}

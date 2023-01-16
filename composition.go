package gen

type lazyGen[K any] struct {
	genOneFunc func() K
}

func (lg lazyGen[K]) Generate() K { return lg.genOneFunc() }

func (lg lazyGen[K]) GenerateN(n uint) []K {
	res := make([]K, n)
	for i := uint(0); i < n; i++ {
		res[i] = lg.genOneFunc()
	}
	return res
}

func Using[T any, K any](gen Gen[T], compositionAction func(T) K) Gen[K] {
	return lazyGen[K]{genOneFunc: func() K { return compositionAction(gen.Generate()) }}
}

type flattenedLazyGen[K any, T any] struct {
	tGen Gen[T]
	gen  func(T) Gen[K]
}

func (f flattenedLazyGen[K, T]) Generate() K {
	tInstance := f.tGen.Generate()
	return f.gen(tInstance).Generate()
}

func (f flattenedLazyGen[K, T]) GenerateN(n uint) []K {
	res := make([]K, n)
	ts := f.tGen.GenerateN(n)
	for i := uint(0); i < n; i++ {
		res[i] = f.gen(ts[i]).Generate()
	}
	return res
}

func UsingGen[T any, K any](gen Gen[T], flatMapFunc func(T) Gen[K]) Gen[K] {
	return flattenedLazyGen[K, T]{gen, flatMapFunc}
}

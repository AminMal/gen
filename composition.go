package gen

type lazyGen[K any] struct {
	genOneFunc func() K
}

func (lg lazyGen[K]) Generate() K { return lg.genOneFunc() }

// Map creates a lazy generator, which when it's Generate method is invoked, it does the composition action on generated value by gen.
func Map[T any, K any](gen Gen[T], compositionAction func(T) K) Gen[K] {
	return lazyGen[K]{genOneFunc: func() K { return compositionAction(gen.Generate())} }
}

// Using is just the same as `Map`, it's for keeping the APIs consistant and a better naming rather than FP's Functor.
func Using[T any, K any](gen Gen[T], compositionAction func(T) K) Gen[K] {
	return Map(gen, compositionAction)
}

type flattenedLazyGen[K any, T any] struct {
	tGen Gen[T]
	gen  func(T) Gen[K]
}

func (f flattenedLazyGen[K, T]) Generate() K {
	tInstance := f.tGen.Generate()
	return f.gen(tInstance).Generate()
}

// FlatMap creates a flattened lazy generator given the base generator as `gen`, and a bind function.
func FlatMap[T any, K any](gen Gen[T], flatMapFunc func(T) Gen[K]) Gen[K] {
	return flattenedLazyGen[K, T]{gen, flatMapFunc}
}

// UsingGen is the same as FlatMap, it's here for consistant APIs and a better naming rather than FP's Monad.
func UsingGen[T any, K any](gen Gen[T], flatMapFunc func(T) Gen[K]) Gen[K] {
	return FlatMap(gen, flatMapFunc)
}

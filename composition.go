package gen

// Map creates a lazy generator, which when it's Generate method is invoked, it does the composition action on generated value by gen.
func Map[T any, K any](gen Gen[T], compositionAction func(T) K) Gen[K] {
    return Pure(func() K { return compositionAction(gen.Generate()) })
}

// Using is just the same as `Map`, it's for keeping the APIs consistant and a better naming rather than FP's Functor.
func Using[T any, K any](gen Gen[T], compositionAction func(T) K) Gen[K] {
	return Map(gen, compositionAction)
}

// FlatMap creates a flattened lazy generator given the base generator as `gen`, and a bind function.
func FlatMap[T any, K any](gen Gen[T], flatMapFunc func(T) Gen[K]) Gen[K] {
    return Pure(func() K {
        return flatMapFunc(gen.Generate()).Generate()
    })
}

// UsingGen is the same as FlatMap, it's here for consistant APIs and a better naming rather than FP's Monad.
func UsingGen[T any, K any](gen Gen[T], flatMapFunc func(T) Gen[K]) Gen[K] {
	return FlatMap(gen, flatMapFunc)
}

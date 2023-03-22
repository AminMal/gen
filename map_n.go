package gen

// Map2 takes 2 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map2[T1 any, T2 any, K any](g1 Gen[T1], g2 Gen[T2], compose func(T1, T2) K) Gen[K] {
	return Pure(func() K { return compose(g1.Generate(), g2.Generate()) })
}

// Map3 takes 3 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map3[T1 any, T2 any, T3 any, K any](g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], compose func(T1, T2, T3) K) Gen[K] {
	return Pure(func() K { return compose(g1.Generate(), g2.Generate(), g3.Generate()) })
}

// Map4 takes 4 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map4[T1, T2, T3, T4, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], compose func(T1, T2, T3, T4) K,
) Gen[K] {
	return Pure(func() K {
		return compose(g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate())
	})
}

// Map5 takes 5 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map5[T1, T2, T3, T4, T5, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], compose func(T1, T2, T3, T4, T5) K,
) Gen[K] {
	return Pure(func() K {
		return compose(g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate())
	})
}

// Map6 takes 6 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map6[T1, T2, T3, T4, T5, T6, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], compose func(T1, T2, T3, T4, T5, T6) K,
) Gen[K] {
	return Pure(func() K {
		return compose(g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate())
	})
}

// Map7 takes 7 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map7[T1, T2, T3, T4, T5, T6, T7, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7], compose func(T1, T2, T3, T4, T5, T6, T7) K,
) Gen[K] {
	return Pure(func() K {
		return compose(g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate())
	})
}

// Map8 takes 8 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map8[T1, T2, T3, T4, T5, T6, T7, T8, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6],
	g7 Gen[T7], g8 Gen[T8], compose func(T1, T2, T3, T4, T5, T6, T7, T8) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(),
			g6.Generate(), g7.Generate(), g8.Generate(),
		)
	})
}

// Map9 takes 9 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map9[T1, T2, T3, T4, T5, T6, T7, T8, T9, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(),
			g7.Generate(), g8.Generate(), g9.Generate(),
		)
	})
}

// Map10 takes 10 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(),
		)
	})
}

// Map11 takes 11 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map11[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], g11 Gen[T11], compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(), g11.Generate(),
		)
	})
}

// Map12 takes 12 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map12[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], g11 Gen[T11], g12 Gen[T12],
	compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(), g11.Generate(), g12.Generate(),
		)
	})
}

// Map13 takes 13 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map13[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], g11 Gen[T11], g12 Gen[T12], g13 Gen[T13],
	compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(), g11.Generate(), g12.Generate(), g13.Generate(),
		)
	})
}

// Map14 takes 14 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map14[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], g11 Gen[T11], g12 Gen[T12], g13 Gen[T13], g14 Gen[T14],
	compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(), g11.Generate(), g12.Generate(), g13.Generate(), g14.Generate(),
		)
	})
}

// Map15 takes 15 generators, and a composition action, and returns a generator which when invoked,
// will use the composition action and the given generators to generate new values
func Map15[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15, K any](
	g1 Gen[T1], g2 Gen[T2], g3 Gen[T3], g4 Gen[T4], g5 Gen[T5], g6 Gen[T6], g7 Gen[T7],
	g8 Gen[T8], g9 Gen[T9], g10 Gen[T10], g11 Gen[T11], g12 Gen[T12], g13 Gen[T13], g14 Gen[T14], g15 Gen[T15],
	compose func(T1, T2, T3, T4, T5, T6, T7, T8, T9, T10, T11, T12, T13, T14, T15) K,
) Gen[K] {
	return Pure(func() K {
		return compose(
			g1.Generate(), g2.Generate(), g3.Generate(), g4.Generate(), g5.Generate(), g6.Generate(), g7.Generate(),
			g8.Generate(), g9.Generate(), g10.Generate(), g11.Generate(), g12.Generate(), g13.Generate(), g14.Generate(), g15.Generate(),
		)
	})
}

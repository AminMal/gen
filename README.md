# Gen

Gen is a type-safe random value generator library for Go, built on generics (Go 1.22+). It is concurrency-safe, has zero external dependencies, and is significantly faster than `testing/quick.Generator` for most operations.

## Installation

```bash
go get github.com/AminMal/gen
```

## Quick Start

```go
nameGen := gen.OneOf("Alice", "Bob", "Charlie")
ageGen := gen.Between(18, 65)

personGen := gen.Map2(nameGen, ageGen, func(name string, age int) Person {
    return Person{Name: name, Age: age}
})

person := personGen.Generate()
people := gen.GenerateN(personGen, 100)
```

## Core Interface

Gen has a single generic interface:
```go
type Gen[T any] interface {
    Generate() T
}
```

All generators implement this interface. You can create custom generators by implementing it yourself, or compose existing ones using the functional combinators below.

## Built-in Generators

### Pure

The basic lazy constructor. Wraps any generation function into a `Gen[T]`:
```go
intGen := gen.Pure(func() int { return rand.IntN(100) })
```

### Only

Always generates the same value:
```go
onlyTwo := gen.Only(2)
values := gen.GenerateN(onlyTwo, 10000) // all 2s
```

### OneOf

Randomly selects from a set of values:
```go
nameGen := gen.OneOf("Bob", "Alice", "Peter", "John")
```

### Between (Numeric)

Generates values within a range (inclusive on both ends). Works with all numeric types — `int`, `uint`, `float32`, `float64`, and their sized variants. Argument order doesn't matter:
```go
ageGen := gen.Between(1, 100)
tempGen := gen.Between(float64(-40.0), float64(50.0))
byteGen := gen.Between(uint8(0), uint8(255)) // full range works correctly
```

The `Numeric` type constraint:
```go
type Numeric interface {
    uint8 | uint16 | uint32 | uint64 | uint | int8 | int16 | int32 | int64 | int | float32 | float64
}
```

### TimeBetween

Random `time.Time` within a range:
```go
now := time.Now()
tenDaysAgo := now.Add(-10 * 24 * time.Hour)
timeGen := gen.TimeBetween(tenDaysAgo, now)
```

### Sequential

Stateful generator producing sequential numeric values. Wraps around when the range is exceeded:
```go
seq := gen.Sequential(1, 100, 2)
a := seq.Generate() // 1
b := seq.Generate() // 3
c := seq.Generate() // 5
```

### TimeSeq

Sequential time generator:
```go
start := time.Now()
end := start.Add(30 * 24 * time.Hour)
step := 24 * time.Hour

g := gen.TimeSeq(start, end, step)
// generates consecutive days within the range
```

## Composition

### Map and FlatMap

Functional composition for building complex generators from simple ones:
```go
nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)

personGen := gen.FlatMap(nameGen, func(name string) gen.Gen[Person] {
    return gen.Map(ageGen, func(age int) Person {
        return Person{Name: name, Age: age}
    })
})
```

### MapN (Recommended)

`Map2` through `Map15` provide a cleaner and significantly faster alternative to nested FlatMap chains. Each takes N generators plus a composition function:

```go
type Programmer struct {
    Name, Surname, GithubUrl, FavLang, Origin string
    Age, Experience                            int
}

programmerGen := gen.Map7(
    nameGen, surnameGen, gitGen, langGen, originGen, ageGen, experienceGen,
    func(name, surname, git, lang, origin string, age, experience int) Programmer {
        return Programmer{name, surname, git, lang, origin, age, experience}
    },
)
```

MapN is the recommended approach for composing multiple generators — it's nearly as fast as a hand-written generator struct and avoids all the allocation overhead of nested FlatMap.

### Infer (Reflection-based)

For quick prototyping, `Infer` uses reflection to automatically generate structs. Provide wrapped generators for field types:
```go
nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)
personGen, err := gen.Infer[Person](
    gen.Wrap(nameGen), gen.Wrap(ageGen),
)
```

**Limitations of Infer:**
- Matches generators to fields by type, not by name — if a struct has two `string` fields, the same generator is used for both.
- Cannot generate function types.
- Significantly slower than MapN due to reflection (see benchmarks below).

For production use, prefer MapN over Infer.

## Arbitrary Values

Pre-built generators for common types:

```go
gen.ArbitraryInt      // int within half-range
gen.ArbitraryInt32    // int32 within half-range
gen.ArbitraryInt64    // int64 within half-range
gen.ArbitraryUint     // uint full range
gen.ArbitraryUint8    // uint8 full range [0, 255]
gen.ArbitraryUint16   // uint16 full range
gen.ArbitraryUint32   // uint32 full range
gen.ArbitraryUint64   // uint64 full range
gen.ArbitraryFloat32  // float32
gen.ArbitraryFloat64  // float64
gen.ArbitraryRune     // rune (int32)
```

### StringGen

Generate random strings with a custom alphabet and length constraints:
```go
hexGen := gen.StringGen("0123456789abcdef", 8, 16)
alphaGen := gen.StringGen("abcdefghijklmnopqrstuvwxyz", 3, 10)
```

## Generating Multiple Values

```go
persons := gen.GenerateN(personGen, 100) // []Person of length 100
```

## Concurrency Safety

All generators are safe to use concurrently from multiple goroutines. The library uses `math/rand/v2` top-level functions which are concurrency-safe by design.

## Benchmarks

Comparison with `testing/quick.Generator` (Apple M5 Pro, Go 1.26):

| Benchmark | gen | quick | Speedup | gen allocs | quick allocs |
|-----------|-----|-------|---------|------------|--------------|
| Only (constant value) | 0.65 ns/op | 0.67 ns/op | ~1x | 0 | 0 |
| Between (random int) | 9.8 ns/op | 8.0 ns/op | ~1x | 1 | 0 |
| OneOf (pick from 1000) | 0.63 ns/op | 10.0 ns/op | **16x** | 0 | 0 |
| Struct (3 fields) | 10.7 ns/op | 40.9 ns/op | **3.8x** | 0 | 3 |
| Struct (7 fields) | 32.5 ns/op | 87.0 ns/op | **2.7x** | 0 | 6 |

Internal comparison of composition strategies (7-field struct):

| Strategy | ns/op | allocs/op |
|----------|-------|-----------|
| Direct (custom Gen impl) | 32.5 | 0 |
| MapN (Map7) | 43.9 | 0 |
| Nested FlatMap/Map | 310 | 12 |
| Infer (reflection) | 1,251 | 104 |

Run benchmarks yourself:
```bash
go test -bench=. -benchmem
```

## Custom Generator Example

Implement the `Gen[T]` interface directly for full control:
```go
type DogGen struct {
    nameGen  gen.Gen[string]
    breedGen gen.Gen[string]
}

func (dg *DogGen) Generate() Dog {
    return Dog{
        Name:  dg.nameGen.Generate(),
        Breed: dg.breedGen.Generate(),
    }
}
```

## License

See [LICENSE.md](LICENSE.md).

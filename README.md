# Gen #

Gen is a random generator library, which is safe (compile-time type checking), rational, and easier to use than `testing/quick.Generator`.

Let's take a look at Gen's base components, functions, variables, and "how to extend"!

# Gen interface #
Gen has a generic interface called `Gen[T]`, one of them is for a single value generation, the other one is to generate a slice of values:
```go
type Gen[T any] interface {
	Generate() T
	GenerateN(n uint) []T
}
```
The logic behind generating depends on the structs/interfaces implementing this interface. You can create various implementations of it as you with to. Now let's take a look at its most common implementations, which already exist in Gen:

## Only ##
`Only` is the most basic generator, as the name declares, it will only generate the value that it's given:
```go
onlyTwo := gen.Only(2)
values := onlyTwo.GenerateN(10000)
// values will be a slice, containing 10000 elements, 
// which all of them have the value of 2
```
Although it may seem unuseful in the first sight, but it can be super-useful specially in property-based-testings, or inferring/creating new generators!

## OneOf ##
`OneOf` is yet another generator, which as its name declares, can select a value among it's given values:
```go
nameGen := gen.OneOf("Bob", "Alice", "Peter", "John")
```
It's useful in many cases, especially when you want to avoid full-randomness, and try to rely on meaningful values.

## Between (Numeric) ##
`Between` is a generator, which only can be used with `Numeric` data types. `Numeric` is a simple type constraint:
```go
type Numeric interface {
	uint8 | uint16 | uint32 | uint64 | uint | int8 | int16 | int32 | int64 | int | float32 | float64
```
The order of the given arguments to `Between` actually does not matter, but it's always a good practice to code what you actually think of:
```go
ageGen := gem.Between(1, 100)
badPracticeAgeGen := gen.Between(100, 1) // still works though!
```
`time.Time` is not a `Numeric`, but there's a function which does the same thing for times!

## TimeBetween ##
The logic is pretty much the same as in `Between`:
```go
now := time.Now()
tenDaysAgo := now.Add(-10 * 24 * time.Hour)
timeGen := gen.TimeBetween(tenDaysAgo, now)
times := timeGen.GenerateN(10)
// will generate 10 time.Time instances in the given interval
```

# Composition #
Being able to generate simple values is not just enough, imagine given a struct below:
```go
type Person struct {
    Name string
    Age  int
}
```
It's not much convenient to create a person generator struct, and implement the functions and the logic from scratch. We should be able to compose already-existing generators to get a new one. In gen, there are two ways of doing this:

1. Putting a small effort and create them using functions.
2. Rely on reflection, and gen does the trick for you .

Both of them would work, the first approach might take a little bit of coding and functional programming involved, but surely it's worth the safety. Let's take a look at an example of each of them.

Note that the `Person` struct here is a simple struct, the case might be different  in your codebase!

## Safe way to compose generators ##
So given the `Person` struct as above, we can create the person generator as follow, using `UsingGen` and `Using` functions:
```go
nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)
personGen := gen.UsingGen(nameGen, func (name string) gen.Gen[Person] {
    return gen.Using(ageGen, func (age int) Person {
        return Person{name, age}
    })
})
```
That's all! Using this approach, you're manually designing the behavior of the generator, without creating a dedicated struct for it. `UsingGen` is basically the `bind` of `flatMap` (Monad) function in FP languages (if you're familiar with FP), while `Using` is basically the `Map` function (from functors). 

## Unsafe yet easy way to compose generators ##
Given the same scenario above, you can provide the base generators, and use the `Infer` function:
```go
nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)
personGen := gen.Infer[Person](
    gen.Wrap(nameGen), gen.Wrap(ageGen),
)
```
Notice that you have to use `Wrap`, because unfortunately, go does not yet support wildcards for generic types. It may seem more convenient than the first approach, so let's compare the two of them.

### Safe approach vs Unsafe approach ###
1- The first downside to the unsafe approach is that you cannot take the full control of the generation logic, firstly because it depends on reflection, and also, it depends on the types of generators.
Say our `Person` struct looked a bit different:
```go
type Perosn struct {
    Name    string
    Age     int
    Surname string
}

nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)
surnameGen := gen.Only("Potter")

personGen := gen.Infer[Person](
    gen.Wrap(nameGen), gen.Wrap(ageGen), gen.Wrap(surnameGen),
)
```
In this case, because the `Infer` function relies on the types of the generators, it uses the last generator of `string`, to generate both name ans the surname! while in the first approach, you're the one who rules!
```go
nameGen := gen.OneOf("Bob", "April")
ageGen := gen.Between(10, 90)

personGen := gen.UsingGen(nameGen, func (name string) gen.Gen[Person] {
    return gen.Using(ageGen, func (age int) Person {
        return Person{name, age, "Potter"}
    })
})
```
2- In the current version of the library, some types are not **yet** supported, like functions!

## Arbitrary Values ##
Generating arbitrary values is so common, that gen already has some arbitrary generators for most-common language types. There are arbitrary generators for these types:
```
int types, uint types, float types, rune and strings
```
They're caleld `Arbitrary` followed by their type name (e.g., `ArbitraryUint32`).

## Randomness ##
Gen uses `math/rand` to arbitrarily create random values under the hood, so it also makes sense if you could take control of that random value. You can use the `Seed` function to seed the random generator:
```go
gen.Seed(int64(6897235))
```
Gen uses current unix millis by default.
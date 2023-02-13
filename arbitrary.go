package gen

import "math"

// ------ int types ------
// ArbitraryInt is an arbitrary int generator within int min value / 2 and int max value / 2
var ArbitraryInt Gen[int] = Between((math.MinInt/2 + 1), (math.MaxInt/2 - 1))

// ArbitraryInt32 is an arbitrary int32 generator within int32 min value / 2 and int32 max value / 2
var ArbitraryInt32 Gen[int32] = Between(int32(math.MinInt32)/2+1, int32(math.MaxInt32)/2-1)

// ArbitraryInt64 is an arbitrary int64 generator within int64 min value / 2 and int64 max value / 2
var ArbitraryInt64 Gen[int64] = Between(int64(math.MinInt64)/2+1, int64(math.MaxInt64)/2-1)

// ------ uint types ------
// ArbitraryUint is an arbitrary uint generator within 0 and uint max value
var ArbitraryUint Gen[uint] = Between(uint(0), uint(math.MaxUint))

// ArbitraryUint8 is an arbitrary uint8 generator within 0 and uint8 max value
var ArbitraryUint8 Gen[uint8] = Between(uint8(0), uint8(math.MaxUint8))

// ArbitraryUint16 is an arbitrary uint16 generator within 0 and uint16 max value
var ArbitraryUint16 Gen[uint16] = Between(uint16(0), uint16(math.MaxUint16))

// ArbitraryUint32 is an arbitrary uint32 generator within 0 and uint32 max value
var ArbitraryUint32 Gen[uint32] = Between(uint32(0), uint32(math.MaxUint32))

// ArbitraryUint64 is an arbitrary uint64 generator within 0 and uint64 max value
var ArbitraryUint64 Gen[uint64] = Between(uint64(0), uint64(math.MaxUint64))

// ------ float types ------
// ArbitraryFloat32 is an arbitrary float32 generator within float32 min value / 2 and float32 max value / 2
var ArbitraryFloat32 Gen[float32] = Between(float32(math.MinInt32)/2+1, float32(math.MaxInt32)/2-1)

// ArbitraryFloat64 is an arbitrary float64 generator within float64 min value / 2 and float32 max value / 2
var ArbitraryFloat64 Gen[float64] = Between(float64(math.MinInt64)/2+1, float64(math.MaxInt64)/2-1)

// ------ rune ------
// ArbitraryRune is an arbitrary rune generator.
var ArbitraryRune Gen[rune] = ArbitraryInt32

// ------ string ------

type stringGen struct {
	alphabet             []rune
	minLength, maxLength int
}

func (s *stringGen) Generate() string {
	strlen := Between(s.minLength, s.maxLength).Generate()
	rs := GenerateN(OneOf(s.alphabet...), uint(strlen))
	return string(rs)
}

func (s *stringGen) GenerateN(n uint) []string {
	res := make([]string, n)
	for i := uint(0); i < n; i++ {
		res[i] = s.Generate()
	}
	return res
}

// StringGen is a string generator that generates random strings using the given alphabet and minLength and maxLength.
func StringGen(alphabet string, minLength uint, maxLength uint) Gen[string] {
	actualMin := numericMin(minLength, maxLength)
	actualMax := numericMax(minLength, maxLength)

	return &stringGen{[]rune(alphabet), int(actualMin), int(actualMax)}
}

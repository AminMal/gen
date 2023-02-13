package gen

// Numeric represents numeric types constraint.
type Numeric interface {
	uint8 | uint16 | uint32 | uint64 | uint | int8 | int16 | int32 | int64 | int | float32 | float64
}

func numericMin[T Numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func numericMax[T Numeric](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func randUint8(n uint8) uint8 {
	return uint8(random.Int31n(int32(n)))
}

func randUint16(n uint16) uint16 {
	return uint16(random.Int31n(int32(n)))
}

func randUint32(n uint32) uint32 {
	return uint32(random.Int31n(int32(n)))
}

func randUint64(n uint64) uint64 {
	return uint64(random.Int63n(int64(n)))
}

func randUint(n uint) uint {
	if n <= 1<<32-1 {
		return uint(randUint32(uint32(n)))
	} else {
		return uint(randUint64(uint64(n)))
	}
}

func randInt8(n int8) int8 {
	return int8(random.Int31n(int32(n)))
}

func randInt16(n int16) int16 {
	return int16(random.Int31n(int32(n)))
}

func randInt32(n int32) int32 {
	return random.Int31n(n)
}

func randInt64(n int64) int64 {
	return random.Int63n(n)
}

func randInt(n int) int {
	return random.Intn(n)
}

func randFloat32(n float32) float32 {
	return random.Float32() * n
}

func randFloat64(n float64) float64 {
	return random.Float64() * n
}

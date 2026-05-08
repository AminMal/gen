package gen

import (
	"math"
	"math/rand/v2"
)

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

// All randXxx functions return a value in [0, n] inclusive.

func randUint8(n uint8) uint8 {
	if n == 0 {
		return 0
	}
	if n == math.MaxUint8 {
		return uint8(rand.UintN(256))
	}
	return uint8(rand.UintN(uint(n) + 1))
}

func randUint16(n uint16) uint16 {
	if n == 0 {
		return 0
	}
	if n == math.MaxUint16 {
		return uint16(rand.UintN(65536))
	}
	return uint16(rand.UintN(uint(n) + 1))
}

func randUint32(n uint32) uint32 {
	if n == 0 {
		return 0
	}
	if n == math.MaxUint32 {
		return rand.Uint32()
	}
	return rand.Uint32N(n + 1)
}

func randUint64(n uint64) uint64 {
	if n == 0 {
		return 0
	}
	if n == math.MaxUint64 {
		return rand.Uint64()
	}
	return rand.Uint64N(n + 1)
}

func randUint(n uint) uint {
	if n == 0 {
		return 0
	}
	if n == math.MaxUint {
		return uint(rand.Uint64())
	}
	return rand.UintN(n + 1)
}

func randInt8(n int8) int8 {
	if n == 0 {
		return 0
	}
	return int8(rand.IntN(int(n) + 1))
}

func randInt16(n int16) int16 {
	if n == 0 {
		return 0
	}
	return int16(rand.IntN(int(n) + 1))
}

func randInt32(n int32) int32 {
	if n == 0 {
		return 0
	}
	return int32(rand.IntN(int(n) + 1))
}

func randInt64(n int64) int64 {
	if n == 0 {
		return 0
	}
	return rand.Int64N(n + 1)
}

func randInt(n int) int {
	if n == 0 {
		return 0
	}
	return rand.IntN(n + 1)
}

func randFloat32(n float32) float32 {
	return rand.Float32() * n
}

func randFloat64(n float64) float64 {
	return rand.Float64() * n
}

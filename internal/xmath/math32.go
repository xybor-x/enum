package xmath

import "math"

const (
	mask32  = 0xFF
	shift32 = 32 - 8 - 1 // total 32 bits - exponent 8 bits - sign 1 bit
	bias32  = 127
)

// Trunc32 is similar to math.Trunc, but for float32 instead.
//
// The implementation is copied and modified from math.Trunc.
func Trunc32(f float32) float32 {
	if f == 0 || IsNaN32(f) || IsInf32(f, 0) {
		return f
	}
	d, _ := Modf32(f)
	return d
}

func IsInf32(f float32, sign int) bool {
	return sign >= 0 && f > math.MaxFloat32 || sign <= 0 && f < -math.MaxFloat32
}

func IsNaN32(f float32) bool {
	return f != f
}

func Modf32(f float32) (int float32, frac float32) {
	return modf32(f)
}

func modf32(f float32) (i float32, frac float32) {
	if f < 1 {
		switch {
		case f < 0:
			i, frac = Modf32(-f)
			return -i, -frac
		case f == 0:
			return f, f // Return -0, -0 when f == -0
		}
		return 0, f
	}

	x := math.Float32bits(f)
	e := uint(x>>shift32)&mask32 - bias32

	// Keep the top 9+e bits, the integer part; clear the rest.
	if e < 32-9 {
		x &^= 1<<(32-9-e) - 1
	}
	i = math.Float32frombits(x)
	frac = f - i
	return
}

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import (
	"fmt"

	"robpike.io/ivy/config"
)

func NewComplex(r, i Value) Value {
	return Complex{real: r, imag: i}.shrink()
}

func newComplexReal(r Value) Complex {
	return Complex{real: r, imag: Int(0)}
}

func newComplexImag(i Value) Complex {
	return Complex{real: Int(0), imag: i}
}

type Complex struct {
	real, imag Value
}

func (z Complex) String() string {
	return "(" + z.Sprint(debugConf) + ")"
}

func (_ Complex) Rank() int {
	return 0
}

func (z Complex) Sprint(conf *config.Config) string {
	return fmt.Sprintf("%sj%s", z.real.Sprint(conf), z.imag.Sprint(conf))
}

func (z Complex) ProgString() string {
	return fmt.Sprintf("%sj%s)", z.real.ProgString(), z.imag.ProgString())
}

func (z Complex) Eval(Context) Value {
	return z
}

func (z Complex) Inner() Value {
	return z
}

func (z Complex) toType(op string, conf *config.Config, which valueType) Value {
	switch which {
	case complexType:
		return z
	case vectorType:
		return NewVector([]Value{z})
	case matrixType:
		return NewMatrix([]int{1}, []Value{z})
	}
	if toBool(z.imag) {
		Errorf("%s: cannot convert complex with non-zero imaginary part to %s", op, which)
		return nil
	}
	return z.real.toType(op, conf, which)
}

// Remove the imaginary part if it is zero.
func (z Complex) shrink() Value {
	if toBool(z.imag) {
		return z
	}
	return z.real
}

// Use EvalUnary to retain ints in real and imaginary parts.
func (z Complex) Floor(c Context) Complex {
	return Complex{c.EvalUnary("floor", z.real), c.EvalUnary("floor", z.imag)}
}

// Use EvalUnary to retain ints in real and imaginary parts.
func (z Complex) Ceil(c Context) Complex {
	return Complex{c.EvalUnary("ceil", z.real), c.EvalUnary("ceil", z.imag)}
}

func (z Complex) Real() Value {
	return z.real
}

func (z Complex) Imag() Value {
	return z.imag
}

// phase a + bi =
//  a = 0, b = 0:  0
//  a = 0, b > 0:  pi/2
//  a = 0, b < 0:  -pi/2
//  a > 0:         atan(b/y)
//  a < 0, b >= 0: atan(b/y) + pi
//  a < 0, b < 0:  atan(b/y) - pi
func (z Complex) Phase(c Context) Value {
	real := floatSelf(c, z.real).(BigFloat).Float
	imag := floatSelf(c, z.imag).(BigFloat).Float
	if real.Sign() == 0 {
		if imag.Sign() == 0 {
			return zero
		} else if imag.Sign() > 0 {
			return BigFloat{newFloat(c).Set(floatHalfPi)}
		} else {
			return BigFloat{newFloat(c).Set(floatMinusHalfPi)}
		}
	}
	slope := newFloat(c)
	slope.Quo(imag, real)
	atan := floatAtan(c, slope)
	if real.Sign() > 0 {
		return BigFloat{atan}.shrink()
	}
	if imag.Sign() >= 0 {
		atan.Add(atan, floatPi)
		return BigFloat{atan}.shrink()
	}
	atan.Sub(atan, floatPi)
	return BigFloat{atan}.shrink()
}

// Use EvalUnary to retain int/rational in real and imaginary parts.
func (z Complex) Neg(c Context) Complex {
	return Complex{c.EvalUnary("-", z.real), c.EvalUnary("-", z.imag)}
}

// sgn z = z / |z|
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Sign(c Context) Value {
	return c.EvalBinary(z, "/", z.Abs(c))
}

// |a+bi| = sqrt (a² + b²)
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Abs(c Context) Value {
	aSq := c.EvalBinary(z.real, "*", z.real)
	bSq := c.EvalBinary(z.imag, "*", z.imag)
	sumSq := c.EvalBinary(aSq, "+", bSq)
	return c.EvalUnary("sqrt", sumSq)
}

// principal square root:
// sqrt(z) = sqrt(|z|) * (z + |z|) / |(z + |z|)|
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Sqrt(c Context) Complex {
	// Avoid division by zero when imaginary part is zero.
	if !toBool(z.imag) {
		return sqrt(c, z.real).toType("sqrt", c.Config(), complexType).(Complex)
	}
	zMod := z.Abs(c)
	sqrtZMod := c.EvalUnary("sqrt", zMod)
	zPlusZMod := c.EvalBinary(z, "+", zMod)
	denom := c.EvalUnary("abs", zPlusZMod)
	num := c.EvalBinary(sqrtZMod, "*", zPlusZMod)
	return c.EvalBinary(num, "/", denom).toType("sqrt", c.Config(), complexType).(Complex)
}

func (z Complex) Cmp(c Context, right Complex) bool {
	return toBool(c.EvalBinary(z.real, "==", right.real)) && toBool(c.EvalBinary(z.imag, "==", right.imag))
}

// (a+bi) + (c+di) = (a+c) + (b+d)i
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Add(c Context, right Complex) Complex {
	return Complex{
		real: c.EvalBinary(z.real, "+", right.real),
		imag: c.EvalBinary(z.imag, "+", right.imag),
	}
}

// (a+bi) - (c+di) = (a-c) + (b-d)i
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Sub(c Context, right Complex) Complex {
	return Complex{
		real: c.EvalBinary(z.real, "-", right.real),
		imag: c.EvalBinary(z.imag, "-", right.imag),
	}
}

// (a+bi) * (c+di) = (ab - bd) + (ad - bc)i
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Mul(c Context, right Complex) Complex {
	ac := c.EvalBinary(z.real, "*", right.real)
	bd := c.EvalBinary(z.imag, "*", right.imag)
	ad := c.EvalBinary(z.real, "*", right.imag)
	bc := c.EvalBinary(z.imag, "*", right.real)
	return Complex{
		real: c.EvalBinary(ac, "-", bd),
		imag: c.EvalBinary(ad, "+", bc),
	}
}

// (a+bi) / (c+di) = (ac + bd)/(c² + d²) + ((bc - ad)/(c² + d²))i
// Use EvalBinary to retain int/rational in real and imaginary parts.
func (z Complex) Quo(c Context, right Complex) Complex {
	ac := c.EvalBinary(z.real, "*", right.real)
	bd := c.EvalBinary(z.imag, "*", right.imag)
	ad := c.EvalBinary(z.real, "*", right.imag)
	bc := c.EvalBinary(z.imag, "*", right.real)
	realNum := c.EvalBinary(ac, "+", bd)
	imagNum := c.EvalBinary(bc, "-", ad)
	cSq := c.EvalBinary(right.real, "*", right.real)
	dSq := c.EvalBinary(right.imag, "*", right.imag)
	denom := c.EvalBinary(cSq, "+", dSq)
	return Complex{
		real: c.EvalBinary(realNum, "/", denom),
		imag: c.EvalBinary(imagNum, "/", denom),
	}
}

// log z = log |z| + (arg z) i
func (z Complex) Log(c Context) Complex {
	return Complex{logn(c, z.Abs(c)), z.Phase(c)}
}

// z log y = log y / log z
func (z Complex) LogBaseU(c Context, right Complex) Complex {
	logZ := z.Log(c)
	logY := right.Log(c)
	return logY.Quo(c, logZ)
}

// exp(a+bi) = (exp(a) * cos b) + (exp(a) *sin b) i
func (z Complex) Exp(c Context) Complex {
	cosB := floatCos(c, floatSelf(c, z.imag).(BigFloat).Float)
	sinB := floatSin(c, floatSelf(c, z.imag).(BigFloat).Float)
	expA := floatPower(c, BigFloat{floatE}, floatSelf(c, z.real).(BigFloat))
	cosB.Mul(cosB, expA)
	sinB.Mul(sinB, expA)
	return Complex{
		real: BigFloat{cosB}.shrink(),
		imag: BigFloat{sinB}.shrink(),
	}
}

// principal solution:
// z**y = exp(y * log z)
func (z Complex) Pow(c Context, right Complex) Complex {
	return z.Log(c).Mul(c, right).Exp(c)
}

// sin(a + bi) = sin(a)*cosh(b) + i*cos(a)*sinh(b)
func (z Complex) Sin(c Context) Complex {
	sinA := floatSin(c, floatSelf(c, z.real).(BigFloat).Float)
	coshB := floatCosh(c, floatSelf(c, z.imag).(BigFloat).Float)
	cosA := floatCos(c, floatSelf(c, z.real).(BigFloat).Float)
	sinhB := floatSinh(c, floatSelf(c, z.imag).(BigFloat).Float)
	return Complex{
		real: BigFloat{newFloat(c).Mul(sinA, coshB)}.shrink(),
		imag: BigFloat{newFloat(c).Mul(cosA, sinhB)}.shrink(),
	}
}

// cos(a + bi) = cos(a)*cosh(b) - i*sin(a)*sinh(b)
func (z Complex) Cos(c Context) Complex {
	cosA := floatCos(c, floatSelf(c, z.real).(BigFloat).Float)
	coshB := floatCosh(c, floatSelf(c, z.imag).(BigFloat).Float)
	sinA := floatSin(c, floatSelf(c, z.real).(BigFloat).Float)
	sinhB := floatSinh(c, floatSelf(c, z.imag).(BigFloat).Float)
	return Complex{
		real: BigFloat{newFloat(c).Mul(cosA, coshB)}.shrink(),
		imag: BigFloat{newFloat(c).Neg(newFloat(c).Mul(sinA, sinhB))}.shrink(),
	}
}

// tan(a + bi) = (sin(2a) + i*sinh(2b))/(cos(2a) + cosh(2b))
func (z Complex) Tan(c Context) Complex {
	twoA := newFloat(c).Mul(floatSelf(c, z.real).(BigFloat).Float, floatTwo)
	twoB := newFloat(c).Mul(floatSelf(c, z.imag).(BigFloat).Float, floatTwo)
	denom := newFloat(c).Add(floatCos(c, twoA), floatCosh(c, twoB))
	return Complex{
		real: BigFloat{newFloat(c).Quo(floatSin(c, twoA), denom)}.shrink(),
		imag: BigFloat{newFloat(c).Quo(floatSinh(c, twoB), denom)}.shrink(),
	}
}

// asin(z) = i log (sqrt(1 - z²) - iz)
func (z Complex) Asin(c Context) Complex {
	sqrt := complexOne.Sub(c, z.Mul(c, z)).Sqrt(c)
	log := sqrt.Sub(c, complexI.Mul(c, z)).Log(c)
	return log.Mul(c, complexI)
}

// acos(z) = log(z + i * sqrt(1 - z²))/i
func (z Complex) Acos(c Context) Complex {
	return complexOne.Sub(c, z.Pow(c, complexTwo)).Sqrt(c).Mul(c, complexI).Add(c, z).Log(c).Quo(c, complexI)
}

// atan(z) = log((i - z)/(i + z))/2i
func (z Complex) Atan(c Context) Complex {
	a := floatSelf(c, z.real).(BigFloat).Float
	b := floatSelf(c, z.imag).(BigFloat).Float
	if a.Cmp(floatZero) == 0 && b.Cmp(floatOne) == 0 {
		Errorf("inverse tangent of 0j1")
	}
	if a.Cmp(floatZero) == 0 && b.Cmp(floatMinusOne) == 0 {
		Errorf("inverse tangent of 0j-1")
	}
	return complexI.Sub(c, z).Quo(c, complexI.Add(c, z)).Log(c).Quo(c, complexTwoI)
}

// sinh(a + bi) = sinh(a)cos(b) + i * cosh(a)sin(b)
func (z Complex) Sinh(c Context) Complex {
	a := floatSelf(c, z.real).(BigFloat).Float
	b := floatSelf(c, z.imag).(BigFloat).Float
	sinha := floatSinh(c, a)
	cosb := floatCos(c, b)
	cosha := floatCosh(c, a)
	sinb := floatSin(c, b)
	return Complex{
		real: BigFloat{sinha.Mul(sinha, cosb)}.shrink(),
		imag: BigFloat{cosha.Mul(cosha, sinb)}.shrink(),
	}
}

// cosh(a + bi) = cosh(a)cos(b) + i * sinh(a)sin(b)
func (z Complex) Cosh(c Context) Complex {
	a := floatSelf(c, z.real).(BigFloat).Float
	b := floatSelf(c, z.imag).(BigFloat).Float
	cosha := floatCosh(c, a)
	cosb := floatCos(c, b)
	sinha := floatSinh(c, a)
	sinb := floatSin(c, b)
	return Complex{
		real: BigFloat{cosha.Mul(cosha, cosb)}.shrink(),
		imag: BigFloat{sinha.Mul(sinha, sinb)}.shrink(),
	}
}

// tanh(z) = sinh(z)/cosh(z)
func (z Complex) Tanh(c Context) Complex {
	return z.Sinh(c).Quo(c, z.Cosh(c))
}

// asinh(z) = log(z + sqrt(z²+1))
func (z Complex) Asinh(c Context) Complex {
	return complexOne.Add(c, z.Pow(c, complexTwo)).Sqrt(c).Add(c, z).Log(c)
}

// acosh(z) = log(z + sqrt(z+1) * sqrt(z-1))
func (z Complex) Acosh(c Context) Complex {
	sqrtZAdd1 := z.Add(c, complexOne).Sqrt(c)
	sqrtZSub1 := z.Sub(c, complexOne).Sqrt(c)
	return z.Add(c, sqrtZAdd1.Mul(c, sqrtZSub1)).Log(c)
}

// alternate: same result?
// acosh(z) = acos(z) * sqrt(z-1) / sqrt(1-z)
/*func (z Complex) Acosh(c Context) Complex {
	sqrtZSub1 := z.Sub(c, complexOne).Sqrt(c)
	sqrt1SubZ := complexOne.Sub(c, z).Sqrt(c)
	return z.Acos(c).Mul(c, sqrtZSub1.Quo(c, sqrt1SubZ))
}*/

// atanh(z) = log((1+z)/(1-z))/2
func (z Complex) Atanh(c Context) Complex {
	onePlusZ := complexOne.Add(c, z)
	oneMinusZ := complexOne.Sub(c, z)
	return onePlusZ.Quo(c, oneMinusZ).Log(c).Quo(c, complexTwo)
}

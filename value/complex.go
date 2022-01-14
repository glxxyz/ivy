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

func (z Complex) shrink() Value {
	if toBool(z.imag) {
		return z
	}
	return z.real
}

func (z Complex) Floor(ctx Context) Complex {
	return Complex{ctx.EvalUnary("floor", z.real), ctx.EvalUnary("floor", z.imag)}
}

func (z Complex) Ceil(ctx Context) Complex {
	return Complex{ctx.EvalUnary("ceil", z.real), ctx.EvalUnary("ceil", z.imag)}
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
func (z Complex) Phase(ctx Context) Value {
	real := floatSelf(ctx, z.real).(BigFloat).Float
	imag := floatSelf(ctx, z.imag).(BigFloat).Float
	if real.Sign() == 0 {
		if imag.Sign() == 0 {
			return zero
		} else if imag.Sign() > 0 {
			return BigFloat{newFloat(ctx).Set(floatHalfPi)}
		} else {
			return BigFloat{newFloat(ctx).Set(floatMinusHalfPi)}
		}
	}
	slope := newFloat(ctx)
	slope.Quo(imag, real)
	atan := floatAtan(ctx, slope)
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

func (z Complex) Neg(ctx Context) Complex {
	return Complex{ctx.EvalUnary("-", z.real), ctx.EvalUnary("-", z.imag)}
}

// sgn z = z / |z|
func (z Complex) Sign(ctx Context) Value {
	return ctx.EvalBinary(z, "/", z.Abs(ctx))
}

// |a+bi| = sqrt (a² + b²)
func (z Complex) Abs(ctx Context) Value {
	aSq := ctx.EvalBinary(z.real, "*", z.real)
	bSq := ctx.EvalBinary(z.imag, "*", z.imag)
	sumSq := ctx.EvalBinary(aSq, "+", bSq)
	return ctx.EvalUnary("sqrt", sumSq)
}

// principal square root:
// sqrt(z) = sqrt(|z|) * (z + |z|) / |(z + |z|)|
func (z Complex) Sqrt(ctx Context) Value {
	zMod := z.Abs(ctx)
	sqrtZMod := ctx.EvalUnary("sqrt", zMod)
	zPlusZMod := ctx.EvalBinary(z, "+", zMod)
	denom := ctx.EvalUnary("abs", zPlusZMod)
	num := ctx.EvalBinary(sqrtZMod, "*", zPlusZMod)
	return ctx.EvalBinary(num, "/", denom)
}

func (z Complex) Cmp(ctx Context, right Complex) bool {
	return toBool(ctx.EvalBinary(z.real, "==", right.real)) && toBool(ctx.EvalBinary(z.imag, "==", right.imag))
}

// (a+bi) + (c+di) = (a+c) + (b+d)i
func (z Complex) Add(ctx Context, right Complex) Complex {
	return Complex{
		real: ctx.EvalBinary(z.real, "+", right.real),
		imag: ctx.EvalBinary(z.imag, "+", right.imag),
	}
}

// (a+bi) - (c+di) = (a-c) + (b-d)i
func (z Complex) Sub(ctx Context, right Complex) Complex {
	return Complex{
		real: ctx.EvalBinary(z.real, "-", right.real),
		imag: ctx.EvalBinary(z.imag, "-", right.imag),
	}
}

// (a+bi) * (c+di) = (ab - bd) + (ad - bc)i
func (z Complex) Mul(ctx Context, right Complex) Complex {
	ac := ctx.EvalBinary(z.real, "*", right.real)
	bd := ctx.EvalBinary(z.imag, "*", right.imag)
	ad := ctx.EvalBinary(z.real, "*", right.imag)
	bc := ctx.EvalBinary(z.imag, "*", right.real)
	return Complex{
		real: ctx.EvalBinary(ac, "-", bd),
		imag: ctx.EvalBinary(ad, "+", bc),
	}
}

// (a+bi) / (c+di) = (ac + bd)/(c² + d²) + ((bc - ad)/(c² + d²))i
func (z Complex) Quo(ctx Context, right Complex) Complex {
	ac := ctx.EvalBinary(z.real, "*", right.real)
	bd := ctx.EvalBinary(z.imag, "*", right.imag)
	ad := ctx.EvalBinary(z.real, "*", right.imag)
	bc := ctx.EvalBinary(z.imag, "*", right.real)
	realNum := ctx.EvalBinary(ac, "+", bd)
	imagNum := ctx.EvalBinary(bc, "-", ad)
	cSq := ctx.EvalBinary(right.real, "*", right.real)
	dSq := ctx.EvalBinary(right.imag, "*", right.imag)
	denom := ctx.EvalBinary(cSq, "+", dSq)
	return Complex{
		real: ctx.EvalBinary(realNum, "/", denom),
		imag: ctx.EvalBinary(imagNum, "/", denom),
	}
}

// principal solution:
// log a+bi = (log a² + b²)/2 + (atan b/a)i
func (z Complex) Log(ctx Context) Complex {
	aSq := ctx.EvalBinary(z.real, "*", z.real)
	bSq := ctx.EvalBinary(z.imag, "*", z.imag)
	sum := ctx.EvalBinary(aSq, "+", bSq)
	log := ctx.EvalUnary("log", sum)
	bdiva := ctx.EvalBinary(z.imag, "/", z.real)
	return Complex{
		real: ctx.EvalBinary(log, "/", Int(2)),
		imag: ctx.EvalUnary("atan", bdiva),
	}
}

// z log y = log y / log z
func (z Complex) LogBaseU(ctx Context, right Complex) Complex {
	logz := z.Log(ctx)
	logy := right.Log(ctx)
	return logy.Quo(ctx, logz)
}

// exp(a+bi) = (exp(a) * cos b) + (exp(a) *sin b) i
func (z Complex) Exp(ctx Context) Complex {
	cosb := floatCos(ctx, floatSelf(ctx, z.imag).(BigFloat).Float)
	sinb := floatSin(ctx, floatSelf(ctx, z.imag).(BigFloat).Float)
	expA := floatPower(ctx, BigFloat{floatE}, floatSelf(ctx, z.real).(BigFloat))
	cosb.Mul(cosb, expA)
	sinb.Mul(sinb, expA)
	return Complex{BigFloat{cosb}.shrink(), BigFloat{sinb}.shrink()}
}

// principal solution:
// z**y = exp(y * log z)
func (z Complex) Pow(ctx Context, right Complex) Complex {
	return z.Log(ctx).Mul(ctx, right).Exp(ctx)
}

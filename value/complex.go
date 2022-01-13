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

func (c Complex) String() string {
	return "(" + c.Sprint(debugConf) + ")"
}

func (c Complex) Rank() int {
	return 0
}

func (c Complex) Sprint(conf *config.Config) string {
	return fmt.Sprintf("%sj%s", c.real.Sprint(conf), c.imag.Sprint(conf))
}

func (c Complex) ProgString() string {
	return fmt.Sprintf("%sj%s)", c.real.ProgString(), c.imag.ProgString())
}

func (c Complex) Eval(Context) Value {
	return c
}

func (c Complex) Inner() Value {
	return c
}

func (c Complex) toType(op string, conf *config.Config, which valueType) Value {
	switch which {
	case complexType:
		return c
	case vectorType:
		return NewVector([]Value{c})
	case matrixType:
		return NewMatrix([]int{1}, []Value{c})
	}
	if toBool(c.imag) {
		Errorf("%s: cannot convert complex with non-zero imaginary part to %s", op, which)
		return nil
	}
	return c.real.toType(op, conf, which)
}

func (c Complex) shrink() Value {
	if toBool(c.imag) {
		return c
	}
	return c.real
}

func (c Complex) Floor(ctx Context) Complex {
	return Complex{ctx.EvalUnary("floor", c.real), ctx.EvalUnary("floor", c.imag)}
}

func (c Complex) Ceil(ctx Context) Complex {
	return Complex{ctx.EvalUnary("ceil", c.real), ctx.EvalUnary("ceil", c.imag)}
}

func (c Complex) Real(_ Context) Value {
	return c.real
}

func (c Complex) Imag(_ Context) Value {
	return c.imag
}

// phase a + bi =
//  a = 0, b = 0:  0
//  a = 0, b > 0:  pi/2
//  a = 0, b < 0:  -pi/2
//  a > 0:         atan(b/y)
//  a < 0, b >= 0: atan(b/y) + pi
//  a < 0, b < 0:  atan(b/y) - pi
func (c Complex) Phase(ctx Context) Value {
	if toBool(ctx.EvalBinary(c.real, "==", zero)) {
		if toBool(ctx.EvalBinary(c.imag, "==", zero)) {
			return zero
		} else if toBool(ctx.EvalBinary(c.imag, ">", zero)) {
			return BigFloat{newF(ctx.Config()).Set(floatHalfPi)}
		} else {
			return BigFloat{newF(ctx.Config()).Set(floatMinusHalfPi)}
		}
	}
	slope := ctx.EvalBinary(c.imag, "/", c.real)
	atan := ctx.EvalUnary("atan", slope)
	if toBool(ctx.EvalBinary(c.real, ">", zero)) {
		return atan
	}
	if toBool(ctx.EvalBinary(c.imag, ">=", zero)) {
		return ctx.EvalBinary(atan, "+", BigFloat{newF(ctx.Config()).Set(floatPi)})
	}
	return ctx.EvalBinary(atan, "-", BigFloat{newF(ctx.Config()).Set(floatPi)})
}

func (c Complex) Neg(ctx Context) Complex {
	return Complex{ctx.EvalUnary("-", c.real), ctx.EvalUnary("-", c.imag)}
}

// sgn z = z / |z|
func (c Complex) Sign(ctx Context) Value {
	return ctx.EvalBinary(c, "/", c.Abs(ctx))
}

// |a+bi| = sqrt (a^2 + b^2)
func (c Complex) Abs(ctx Context) Value {
	aSq := ctx.EvalBinary(c.real, "*", c.real)
	bSq := ctx.EvalBinary(c.imag, "*", c.imag)
	sumSq := ctx.EvalBinary(aSq, "+", bSq)
	return ctx.EvalUnary("sqrt", sumSq)
}

// principal square root:
// sqrt(z) = sqrt(|z|) * (z + |z|) / |(z + |z|)|
// sqrt(z) = sqrt(a + bi) = sqrt((|z|+a)/2) + sgn(b) * sqrt((|z|-a)/2)i
func (c Complex) Sqrt(ctx Context) Value {
	zMod := c.Abs(ctx)
	sqrtZMod := ctx.EvalUnary("sqrt", zMod)
	zPlusZMod := ctx.EvalBinary(c, "+", zMod)
	denom := ctx.EvalUnary("abs", zPlusZMod)
	num := ctx.EvalBinary(sqrtZMod, "*", zPlusZMod)
	return ctx.EvalBinary(num, "/", denom)
}

func (c Complex) Cmp(ctx Context, right Complex) bool {
	return toBool(ctx.EvalBinary(c.real, "==", right.real)) && toBool(ctx.EvalBinary(c.imag, "==", right.imag))
}

// (a+bi) + (c+di) = (a+c) + (b+d)i
func (c Complex) Add(ctx Context, right Complex) Complex {
	return Complex{
		real: ctx.EvalBinary(c.real, "+", right.real),
		imag: ctx.EvalBinary(c.imag, "+", right.imag),
	}
}

// (a+bi) - (c+di) = (a-c) + (b-d)i
func (c Complex) Sub(ctx Context, right Complex) Complex {
	return Complex{
		real: ctx.EvalBinary(c.real, "-", right.real),
		imag: ctx.EvalBinary(c.imag, "-", right.imag),
	}
}

// (a+bi) * (c+di) = (ab - bd) + (ad - bc)i
func (c Complex) Mul(ctx Context, right Complex) Complex {
	ac := ctx.EvalBinary(c.real, "*", right.real)
	bd := ctx.EvalBinary(c.imag, "*", right.imag)
	ad := ctx.EvalBinary(c.real, "*", right.imag)
	bc := ctx.EvalBinary(c.imag, "*", right.real)
	return Complex{
		real: ctx.EvalBinary(ac, "-", bd),
		imag: ctx.EvalBinary(ad, "+", bc),
	}
}

// (a+bi) / (c+di) = (ac + bd)/(c^2 + d^2) + ((bc - ad)/(c^2 + d^2))i
func (c Complex) Quo(ctx Context, right Complex) Complex {
	ac := ctx.EvalBinary(c.real, "*", right.real)
	bd := ctx.EvalBinary(c.imag, "*", right.imag)
	ad := ctx.EvalBinary(c.real, "*", right.imag)
	bc := ctx.EvalBinary(c.imag, "*", right.real)
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

// log a+bi = (log a^2 + b^2)/2 + (atan b/a)i
func (c Complex) Log(ctx Context) Complex {
	aSq := ctx.EvalBinary(c.real, "*", c.real)
	bSq := ctx.EvalBinary(c.imag, "*", c.imag)
	sum := ctx.EvalBinary(aSq, "+", bSq)
	log := ctx.EvalUnary("log", sum)
	bdiva := ctx.EvalBinary(c.imag, "/", c.real)
	return Complex{
		real: ctx.EvalBinary(log, "/", Int(2)),
		imag: ctx.EvalUnary("atan", bdiva),
	}
}

// u log v = log v / log u
func (c Complex) LogBaseU(ctx Context, right Complex) Complex {
	logu := c.Log(ctx)
	logv := right.Log(ctx)
	return logv.Quo(ctx, logu)
}

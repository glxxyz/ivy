// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import "math/big"

// domain: (−∞,∞)
func sinh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatSinh)
}

// domain: (−∞,∞)
func cosh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatCosh)
}

// domain: (−∞,∞)
func tanh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatTanh)
}

// domain: (−∞,∞)
func asinh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatAsinh)
}

// domain: [1,∞)
func acosh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatAcosh)
}

// domain: (−1,1)
func atanh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatAtanh)
}

// sinh x = (e**x - e**-x)/2
func floatSinh(c Context, x *big.Float) *big.Float {
	expX := exponential(c.Config(), x)
	expNegX := exponential(c.Config(), newFloat(c).Neg(x))
	num := newFloat(c).Sub(expX, expNegX)
	return newFloat(c).Quo(num, floatTwo)
}

// cosh x = (e**x + e**-x)/2
func floatCosh(c Context, x *big.Float) *big.Float {
	expX := exponential(c.Config(), x)
	expNegX := exponential(c.Config(), newFloat(c).Neg(x))
	num := newFloat(c).Add(expX, expNegX)
	return newFloat(c).Quo(num, floatTwo)
}

// tanh x = (e**x - e**-x)/(e**x + e**-x)
func floatTanh(c Context, x *big.Float) *big.Float {
	expX := exponential(c.Config(), x)
	expNegX := exponential(c.Config(), newFloat(c).Neg(x))
	num := newFloat(c).Sub(expX, expNegX)
	denom := newFloat(c).Add(expX, expNegX)
	return newFloat(c).Quo(num, denom)
}

// asinh x = log(x + sqrt(x**2 + 1))
func floatAsinh(c Context, x *big.Float) *big.Float {
	xSq := newFloat(c).Mul(x, x)
	xSq.Add(xSq, floatOne)
	return floatLog(c, floatSqrt(c, xSq))
}

// acosh x = log(x + sqrt(x**2 - 1))
// Domain: 1 <= x < +Inf
func floatAcosh(c Context, x *big.Float) *big.Float {
	if x.Cmp(floatOne) < 0 {
		Errorf("acosh of value less than 1")
	}
	xSq := newFloat(c).Mul(x, x)
	xSq.Sub(xSq, floatOne)
	return floatLog(c, floatSqrt(c, xSq))
}

// atanh x = log((1 + x)/(1 - x)/2
// Domain: -1 < x < 1
func floatAtanh(c Context, x *big.Float) *big.Float {
	if x.Cmp(floatMinusOne) <= 0 {
		Errorf("atanh of value less than or equal to -1")
	}
	if x.Cmp(floatOne) >= 0 {
		Errorf("atanh of value greater than or equal to 1")
	}
	oneAddX := newFloat(c).Add(floatOne, x)
	oneSubX := newFloat(c).Sub(floatOne, x)
	log := floatLog(c, newFloat(c).Quo(oneAddX, oneSubX))
	return log.Quo(log, floatTwo)
}

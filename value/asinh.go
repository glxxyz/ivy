// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import "math/big"

func asinh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatAsinh)
}

// domain: [1, ∞) - complex solution outside of domain
func acosh(c Context, v Value) Value {
	x := floatSelf(c, v).(BigFloat).Float
	if x.Cmp(floatOne) < 0 {
		return newComplexReal(v).Acosh(c).shrink()
	}
	return evalFloatFunc(c, v, floatAcosh)
}

// domain: (−1,1) - complex solution outside of domain
func atanh(c Context, v Value) Value {
	x := floatSelf(c, v).(BigFloat).Float
	if x.Cmp(floatMinusOne) < 0 || x.Cmp(floatOne) > 0 {
		return newComplexReal(v).Atanh(c).shrink()
	}
	return evalFloatFunc(c, v, floatAtanh)
}

// asinh x = log(x + sqrt(x² + 1))
func floatAsinh(c Context, x *big.Float) *big.Float {
	xSq := newFloat(c).Mul(x, x)
	sqrt := floatSqrt(c, newFloat(c).Add(xSq, floatOne))
	return floatLog(c, newFloat(c).Add(sqrt, x))
}

// acosh x = log(x + sqrt(x² - 1))
// domain: [1, ∞)
func floatAcosh(c Context, x *big.Float) *big.Float {
	if x.Cmp(floatOne) < 0 {
		Errorf("acosh of value less than 1")
	}
	xSq := newFloat(c).Mul(x, x)
	sqrt := floatSqrt(c, newFloat(c).Sub(xSq, floatOne))
	return floatLog(c, newFloat(c).Add(sqrt, x))
}

// atanh x = log((1 + x)/(1 - x)/2
// domain: (−1,1)
func floatAtanh(c Context, x *big.Float) *big.Float {
	switch x.Cmp(floatMinusOne) {
	case -1:
		Errorf("atanh of value less than -1")
	case 0:
		return floatMinusInf
	}
	switch x.Cmp(floatOne) {
	case 1:
		Errorf("atanh of value greater than 1")
	case 0:
		return floatInf
	}
	oneAddX := newFloat(c).Add(floatOne, x)
	oneSubX := newFloat(c).Sub(floatOne, x)
	log := floatLog(c, newFloat(c).Quo(oneAddX, oneSubX))
	return log.Quo(log, floatTwo)
}

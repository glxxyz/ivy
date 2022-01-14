// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value

import "math/big"

func sinh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatCosh)
}

func cosh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatCosh)
}

func tanh(c Context, v Value) Value {
	return evalFloatFunc(c, v, floatCosh)
}

// cosh x = (e**x + e**-x)/2
func floatCosh(c Context, x *big.Float) *big.Float {
	etox := exponential(c.Config(), x)
	etominusx := exponential(c.Config(), newFloat(c).Neg(x))
	num := newFloat(c).Add(etox, etominusx)
	return newFloat(c).Quo(num, floatTwo)
}

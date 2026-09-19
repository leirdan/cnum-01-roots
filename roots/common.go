package roots

import "cnum/types"

func Deriv(x float32, f types.Function) float32 {
	const h float32 = 1e-2
	return (f(x+h) - f(x)) / h
}

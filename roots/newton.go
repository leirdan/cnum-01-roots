package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

func Newton(a, b, epsilon float32, f types.Function, kmax uint16) (float32, uint16) {
	var root float32
	var counter uint16 = 1
	for root = a + float32(rand.Float64())*(b-a); math.Abs(float64(f(root))) > float64(epsilon) && counter < kmax; root = n_candidate(root, f) {
		counter++
	}

	return root, counter
}

func n_candidate(x float32, f types.Function) float32 {
	return x - (f(x) / Deriv(x, f))
}

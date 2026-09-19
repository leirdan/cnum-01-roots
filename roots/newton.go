package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

func Newton(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	var a = inst.Start
	var b = inst.End
	var root float64
	var counter uint16 = 1
	for root = a + rand.Float64()*(b-a); math.Abs(inst.Function(root)) > epsilon && counter < kmax; root = n_candidate(root, inst) {
		counter++
	}

	return root, counter
}

func n_candidate(x float64, inst types.Problem) float64 {
	return x - (inst.Function(x) / inst.Deriv(x))
}

package roots

import (
	"cnum/types"
	"math"
)

func Bisection(inst types.Problem, epsilon float64) (float64, uint16) {
	a_ := inst.Start
	b_ := inst.End
	var root float64
	var counter uint16 = 1

	for root = b_candidate(a_, b_); math.Abs(inst.Function(root)) > epsilon && (b_-a_) > epsilon; root = b_candidate(a_, b_) {
		if inst.Function(a_)*inst.Function(root) < 0 {
			b_ = root
		} else if inst.Function(b_)*inst.Function(root) < 0 {
			a_ = root
		}
		counter++
	}

	return root, counter
}

func b_candidate(a, b float64) float64 {
	return (a + b) / 2
}

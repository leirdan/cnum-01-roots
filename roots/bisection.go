package roots

import (
	"cnum/types"
	"math"
)

func Bisection(a, b, h, epsilon float32, f types.Function) (float32, int16) {
	a_ := a
	b_ := b
	var root float32
	var counter int16 = 1

	for root = b_candidate(a, b); math.Abs(float64(f(root))) > float64(epsilon) && (b_-a_) > epsilon; root = b_candidate(a_, b_) {
		if f(a_)*f(root) < 0 {
			b_ = root
		} else if f(b_)*f(root) < 0 {
			a_ = root
		}
		counter++
	}

	return root, counter
}

func b_candidate(a, b float32) float32 {
	return (a + b) / 2
}

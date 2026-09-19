package roots

import (
	"cnum/types"
	"math"
)

func FalsePosition(inst types.Problem, epsilon float64) (float64, uint16) {
	a_ := inst.Start
	b_ := inst.End
	fa_ := inst.Function(a_)
	fb_ := inst.Function(b_)
	root := fp_candidate(a_, b_, fa_, fb_)
	froot := inst.Function(root)
	var counter uint16 = 1

	for math.Abs(fb_-fa_) > epsilon && math.Abs(froot) > epsilon {
		if fa_*froot < 0 {
			b_ = root
			fb_ = froot
		} else {
			a_ = root
			fa_ = froot
		}

		root = fp_candidate(a_, b_, fa_, fb_)
		froot = inst.Function(root)

		counter++
	}

	return root, counter
}

func fp_candidate(a, b, fa, fb float64) float64 {
	return (a*fb - b*fa) / (fb - fa)
}

package roots

import (
	"cnum/types"
	"math"
)

func FalsePosition(a, b, h, epsilon float32, f types.Function) (float32, int16) {
	// root := candidate(a, b)
	a_ := a
	b_ := b
	fa_ := f(a)
	fb_ := f(b)
	root := step(a_, b_, fa_, fb_)
	froot := f(root)
	var counter int16 = 1

	for math.Abs(float64(fb_-fa_)) > float64(epsilon) && math.Abs(float64(froot)) > float64(epsilon) {

		if fa_*froot < 0 {
			b_ = root
			fb_ = froot
		} else {
			a_ = root
			fa_ = froot
		}

		root = step(a_, b_, fa_, fb_)
		froot = f(root)

		counter++
	}

	return root, counter
}

func step(a, b, fa, fb float32) float32 {
	return (a*fb - b*fa) / (fb - fa)
}

package roots

import (
	"cnum/types"
	"math"
)

// FixedPoint finds a root of inst.Function using the Fixed Point method.
// Instead of requiring a manually supplied iteration function g(x), it
// automatically builds g(x) = x - f(x)/λ, with λ estimated by pf_lambda from
// f itself, guaranteeing convergence (|g'(x)| < 1) within the isolated
// interval.
//
// Input:
//   - inst: already isolated problem, using inst.Function, inst.Start and
//     inst.End (inst.Start/inst.End define the interval [a,b] where the
//     root was isolated)
//   - epsilon: desired precision, stopping criterion |f(root)| <= epsilon
//   - kmax: maximum number of iterations allowed
//
// Output: the approximate root and the number of iterations performed
// until the stopping criterion is met (or kmax is reached).
func FixedPoint(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	lambda := pf_lambda(inst)
	root := b_candidate(inst.Start, inst.End)
	var counter uint16 = 1

	for math.Abs(inst.Function(root)) > epsilon && counter < kmax {
		root = pf_candidate(root, inst, lambda)
		counter++
	}

	return root, counter
}

// pf_candidate computes the next approximation x_{k+1} = x_k - f(x_k)/λ of
// the fixed point iteration.
//
// Input: x (current approximation), inst (to access inst.Function) and
// lambda (constant estimated by pf_lambda).
//
// Output: the next approximation x_{k+1}.
func pf_candidate(x float64, inst types.Problem, lambda float64) float64 {
	return x - inst.Function(x)/lambda
}

// pf_lambda estimates the constant λ used to build g(x) = x - f(x)/λ.
// It samples the derivative of f at several points of the interval [Start,
// End] and takes the largest absolute value found, with the same sign as
// the derivative at the midpoint and a 5% safety margin. This guarantees
// |λ| > max|f'(x)| in the interval, a sufficient condition for
// |g'(x)| = |1 - f'(x)/λ| < 1 and for the fixed point iteration to converge.
//
// Input: inst (uses inst.Start, inst.End and inst.Deriv).
//
// Output: the estimated value of λ.
func pf_lambda(inst types.Problem) float64 {
	const samples = 20

	mid := b_candidate(inst.Start, inst.End)
	sign := math.Copysign(1, inst.Deriv(mid))
	maxAbs := math.Abs(inst.Deriv(inst.Start))
	step := (inst.End - inst.Start) / samples

	for i := 1; i <= samples; i++ {
		x := inst.Start + float64(i)*step
		d := math.Abs(inst.Deriv(x))
		if d > maxAbs {
			maxAbs = d
		}
	}

	const margem = 1.05
	return sign * maxAbs * margem
}

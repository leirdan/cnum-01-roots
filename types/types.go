package types

type Function func(float64) float64
type Interval []float64
type RootIntervalList []Interval

type Problem struct {
	Function Function
	Start    float64
	End      float64
	Id       uint8
	H        float64
}

type RootMethod func(inst Problem, epsilon float64) (float64, uint16)

type NamedMethod struct {
	Name   string
	Method RootMethod
}

type RaceResult struct {
	Root    float64
	Counter uint16
}

type ResultStr struct {
	Root float64
	Iter uint16
}

// IsolateRoots runs the isolation step over the problem (phase 1 of the
// numerical method: locating subintervals of [Start, End] that contain a
// root, before refining with Bisection/False Position/Newton/Secant/Fixed
// Point).
//
// Input: none besides Problem itself (uses self.Start, self.End,
// self.Function and self.H as the tabulation step size).
//
// Output: the list of subintervals [xmin, xmax] in which the existence of
// at least one isolated root was verified (see isolateRec).
func (self Problem) IsolateRoots() RootIntervalList {
	return isolateRec(self.Start, self.End, self.Function, self.H)
}

// Deriv estimates f'(x) at point x, used to apply the unique-root corollary
// in isolateRec and to build λ in the Fixed Point method.
//
// Input: x (point where the derivative is evaluated).
// Output: the estimate of f'(x).
func (self Problem) Deriv(x float64) float64 {
	return deriv(x, self.Function)
}

// isolateRec sweeps [a,b] in steps of size "step" (tabulation) and applies
// the Bolzano Theorem: if f is continuous on [xmin,xmax] and
// f(xmin)*f(xmax) < 0, then there exists at least one root ξ in (xmin,xmax)
// such that f(ξ) = 0. Each sign change found between two consecutive
// tabulation points (done by nextInterval) marks one of these candidate
// subintervals.
//
// It then applies the root-uniqueness corollary: if f'(x) doesn't change
// sign on (xmin,xmax) (i.e. f is monotonic there, verified by d1*d2 >= 0 at
// the endpoints of the interval), the root in that subinterval is unique
// and the interval is accepted directly. Otherwise (f' changes sign, f is
// not monotonic), there may be more than one root there, so the
// subinterval itself is recursively re-isolated with the same step until
// unique roots are isolated.
//
// Input: a, b (bounds of the search interval), f (function to analyze),
// step (tabulation step size h).
//
// Output: list of subintervals [xmin, xmax], each containing exactly one
// isolated root of f.
func isolateRec(a, b float64, f Function, step float64) RootIntervalList {
	var collection = RootIntervalList{}
	var curr float64 = a

	for curr < b {
		sequence := nextInterval(curr, b, f, step)
		curr = sequence[len(sequence)-1] // update to the element it stopped at
		n := len(sequence)
		if n == 1 {
			continue
		}
		xmin := sequence[n-2]
		xmax := sequence[n-1]
		// Bolzano Theorem -> f(xmin)*f(xmax) < 0 guarantees at least 1 root in (xmin, xmax)
		// If there's no sign change, there's no guaranteed root and the segment is discarded
		if f(xmin)*f(xmax) < 0 {
			d1 := deriv(xmin, f)
			d2 := deriv(xmax, f)
			if d1*d2 >= 0 { // corollary: constant-sign derivative -> unique root, good
				collection = append(collection, Interval{xmin, xmax})
			} else {
				subintervals := isolateRec(xmin, xmax, f, step)
				collection = append(collection, subintervals...)
			}
		}
	}

	return collection
}

// deriv estimates the derivative of f at point x using the forward finite
// differences method: f'(x) ≈ (f(x+h) - f(x)) / h, with a small fixed h.
//
// Input: x (evaluation point), f (function to differentiate).
// Output: the estimate of f'(x).
func deriv(x float64, f Function) float64 {
	const h float64 = 1e-2
	return (f(x+h) - f(x)) / h
}

// nextInterval tabulates f starting from "a", walking in steps of size
// "step" until it finds a sign change relative to f(a) (Bolzano Theorem
// condition) or until it reaches "b". The last pair of points in the
// returned slice is the candidate subinterval for containing a root,
// consumed by isolateRec.
//
// Input: a, b (bounds of the remaining interval to sweep), f (tabulated
// function), step (step size h).
//
// Output: the sequence of points [a, a+step, a+2*step, ...] traversed until
// the sign change (or until b), used both for the tabulation and to find
// the subinterval [xmin, xmax] where f changes sign.
func nextInterval(a, b float64, f Function, step float64) Interval {
	var slice Interval
	curr := a
	for curr < b && f(a)*f(curr) >= 0 {
		slice = append(slice, curr)
		curr += step
	}

	if curr <= b {
		slice = append(slice, curr)
	} else {
		slice = append(slice, b)
	}

	return slice
}

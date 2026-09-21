package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
	"sync"
)

// Chaotic implements the strategy proposed in item 1.3: a "chaotic race"
// where Newton, Secant and Fixed Point compete for the same approximate
// root, held in a single shared variable (sharedRoot). Each method reads
// the most recent value written by any of the others and computes its own
// next approximation from it — it doesn't run in isolation with its own
// sequence, but "surfs" on the progress made by the others.
//
// sharedRoot is protected by a sync.Mutex (via read/write) because the four
// goroutines below (Newton, Secant, Fixed Point and the watcher) read and
// write it concurrently; without this mutual exclusion the access is a real
// data race (confirmed with `go run -race`), which could corrupt the
// read/written value.
//
// The watcher goroutine also guarantees that the value accepted as the
// result stays within the isolated interval [inst.Start, inst.End]: since
// the methods write to the same variable without coordination, a Newton or
// Fixed Point iteration can eventually "push" sharedRoot outside the
// isolated interval and converge to another root of f, outside the
// subinterval it was asked to refine. If that happens, the value is
// discarded and the search restarts within the interval, instead of
// accepting a root from outside it.
//
// Input:
//   - inst: already isolated problem (uses inst.Start/inst.End as the
//     bounds of the valid interval, inst.Function and, for the Fixed Point
//     goroutine, the same g(x) = x - f(x)/λ automatically built by
//     pf_lambda, as in FixedPoint)
//   - epsilon: desired precision, stopping criterion |f(root)| < epsilon
//   - kmax: maximum number of iterations of the watcher goroutine
//
// Output: the approximate root (always within [inst.Start, inst.End]) and
// the number of iterations the watcher goroutine took to accept it.
func Chaotic(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	var mu sync.Mutex
	sharedRoot := (inst.Start + inst.End) / 2.0

	read := func() float64 {
		mu.Lock()
		defer mu.Unlock()
		return sharedRoot
	}
	write := func(v float64) {
		mu.Lock()
		sharedRoot = v
		mu.Unlock()
	}

	done := make(chan struct{})
	chResult := make(chan types.ResultStr)

	// Newton goroutine
	go func() {
		h := 0.001
		for {
			select {
			case <-done:
				return
			default:
				x := read()

				fx := inst.Function(x)
				fxh := inst.Function(x + h)
				deriv := (fxh - fx) / h

				if deriv != 0 {
					write(x - (fx / deriv))
				}
			}
		}
	}()

	// Secant goroutine
	go func() {
		xPrev := inst.Start
		for {
			select {
			case <-done:
				return
			default:
				xCurr := read()
				fxCurr := inst.Function(xCurr)
				fxPrev := inst.Function(xPrev)

				if fxCurr != fxPrev {
					next := (xPrev*fxCurr - xCurr*fxPrev) / (fxCurr - fxPrev)
					write(next)
					xPrev = xCurr
				}
			}
		}
	}()

	// Fixed Point goroutine
	go func() {
		lambda := pf_lambda(inst)
		for {
			select {
			case <-done:
				return
			default:
				write(pf_candidate(read(), inst, lambda))
			}
		}
	}()

	// Watcher goroutine
	go func() {
		for i := uint16(1); i < kmax; i++ {
			snapshot := read()

			foraDoIntervalo := snapshot < inst.Start || snapshot > inst.End
			if math.IsNaN(snapshot) || math.IsInf(snapshot, 0) || foraDoIntervalo {
				// some goroutine diverged or "escaped" the isolated interval;
				// discard it and restart the search within [Start, End]
				write(inst.Start + rand.Float64()*(inst.End-inst.Start))
				continue
			}
			if math.Abs(inst.Function(snapshot)) < epsilon {
				chResult <- types.ResultStr{Root: snapshot, Iter: i}
				return
			}
		}
	}()

	final := <-chResult
	close(done)

	return final.Root, final.Iter
}

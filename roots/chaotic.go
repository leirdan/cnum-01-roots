package roots

import (
	"cnum/types"
	"math"
	"math/rand/v2"
)

// Chaotic implements the "chaotic race" of item 1.3: Newton, Secant and
// Fixed Point goroutines compete over a single shared approximation
// (sharedRoot, initially the midpoint of the interval), each computing its
// next value from whatever the others wrote last.
//
// sharedRoot is deliberately accessed with no synchronization, so this is a
// real data race (`go run -race` reports it) and results vary between runs.
//
// A watcher goroutine snapshots sharedRoot and, if it is NaN, ±Inf or outside
// [inst.Start, inst.End], resets it to a random point in the interval;
// otherwise it accepts the snapshot once |f(snapshot)| < epsilon. It runs at
// most kmax-1 iterations; if none is accepted, Chaotic blocks forever.
//
// Returns the root (within [inst.Start, inst.End]) and the watcher iteration
// at which it was accepted.
func Chaotic(inst types.Problem, epsilon float64, kmax uint16) (float64, uint16) {
	sharedRoot := (inst.Start + inst.End) / 2.0

	read := func() float64 {
		return sharedRoot
	}
	write := func(v float64) {
		sharedRoot = v
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

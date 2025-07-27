package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

// Qubit represents the state of a single qubit as amplitudes of |0> and |1>.
type Qubit struct {
	amp [2]complex128
}

// hadamard applies the Hadamard gate to the qubit.
func (q *Qubit) hadamard() {
	a0, a1 := q.amp[0], q.amp[1]
	inv := complex(1/math.Sqrt2, 0)
	q.amp[0] = (a0 + a1) * inv
	q.amp[1] = (a0 - a1) * inv
}

// measureX measures the qubit in the X basis and returns 0 for |+> and 1 for |->.
func (q *Qubit) measureX() int {
	// After applying H, |0> -> |+> and |1> -> |->.
	// The larger amplitude determines the deterministic measurement result in this
	// simplified simulation since the state will be either |0> or |1>.
	if cmplx.Abs(q.amp[0]) > cmplx.Abs(q.amp[1]) {
		return 0
	}
	return 1
}

// distinguishIZ takes a single-qubit operation op which is either the identity
// gate or the Z gate, and returns 0 if op is I and 1 if op is Z.
// The operation is allowed to be applied exactly once.
func distinguishIZ(op func(*Qubit)) int {
	var q Qubit
	q.amp[0] = 1
	q.hadamard() // prepare |+>
	op(&q)       // apply the unknown gate once
	q.hadamard() // rotate to computational basis
	return q.measureX()
}

// identity gate implementation
func iGate(q *Qubit) {}

// zGate applies the Pauli Z gate.
func zGate(q *Qubit) {
	q.amp[1] = -q.amp[1]
}

func main() {
	fmt.Println(distinguishIZ(iGate)) // expected 0
	fmt.Println(distinguishIZ(zGate)) // expected 1
}

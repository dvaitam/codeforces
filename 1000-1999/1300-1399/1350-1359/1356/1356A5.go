package main

// DistinguishGate determines whether the provided single-qubit operation
// is the Z gate or the -Z gate. The algorithm uses the operation exactly
// once. It prepares the |+\u27e9 state, applies the unknown gate, then
// measures in the X basis. Measurement result 0 corresponds to Z, and
// measurement result 1 corresponds to -Z.
func DistinguishGate(q Qubit, op func(Qubit)) int {
	// Prepare |+> = H|0>
	H(q)
	// Apply the unknown operation
	op(q)
	// Rotate back to computational basis
	H(q)
	// Measure: 0 -> Z, 1 -> -Z
	return M(q)
}

// Qubit abstracts a single qubit in the testing harness.
type Qubit interface{}

// The following gate and measurement stubs are assumed to be provided by the
// quantum simulator environment.
func H(Qubit)     {}
func M(Qubit) int { return 0 }

func main() {
	// Execution is handled by the harness which provides a qubit and the
	// operation to test; this stub leaves main empty.
}
